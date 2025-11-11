local class = require 'ext.class'
local table = require 'ext.table'
local range = require 'ext.range'
local number = require 'ext.number'
local math = require 'ext.math'
local assert = require 'ext.assert'

local BigNumber = class()

-- default values
BigNumber.negative = false
BigNumber.nan = false
BigNumber.infinity = false
BigNumber.base = 10

function BigNumber:init(n, base)
-- hmm... this causes errors in longIntDiv
-- but without this I can't make the assert.ion in toBase() that we are getting an integer
--	self.minExp = 0
--	self.maxExp = 0
	self.base = base
	assert.gt(self.base, 1, "can't set to a base of 1 or lower")
	-- TODO add repeating fraction construction using []'s in the string
	if type(n) == 'string' then
		if n == 'nan' then self.nan = true end
		if n == '0' then return end
		if n:sub(1,1) == '-' then
			n = n:sub(2)
			self.negative = true
		end
		local d = n:find'%.'
		local found = d
		d = d or #n+1
		for i=1,d-1 do
			local v = number.todigit(n:sub(i,i))
			assert(v >= 0 and v < self.base)
			self[d-1-i] = math.floor(v)
		end
		self.maxExp = d-1
		if found then
			for i=d+1,#n do
				local v = number.todigit(n:sub(i,i))
				assert.ge(v, 0)
				assert.lt(v, self.base)
				self[d-i] = math.floor(v)
			end
			self.minExp = d-#n
		else
			self.minExp = 0
		end
	elseif type(n) == 'number' then
		if n ~= n then
			self.nan = true
		else
			if n < 0 then
				self.negative = true
				n = -n
			end
			if n == math.huge then
				self.infinity = true
			else
				-- TODO use longIntDiv, esp with non-10 bases, so we can take advantage of repeating digits
				-- TODO between this and toBase(), the code is very similar
				-- and toBase uses BigNumbers for accuracy...
				-- and this allows for fractional bases...
				-- TODO first convert the number in base-10 or base-2 or something
				-- (using ffi we can assign it to a double[1] and then use its bits to determine the number exactly)
				-- then use toBase to convert it to whatever base is desired
				-- and change that code to take decimals (and repeating decimals) into account
				if n > 0 then
					local f = math.floor(n)
					local i = 0
					self.minExp = 0
					while f > 0 do
						self[i] = f % self.base
						f = math.floor((f - self[i]) / self.base)
						self.maxExp = i
						i = i + 1
					end
					local d = n - math.floor(n)
					if d > 0 then
						d = d * self.base
						i = -1
						while d > 0 do
							self[i] = math.floor(d) % self.base
							d = (d - self[i]) * self.base
							self.minExp = i
							i = i - 1
						end
					end
				end
			end
		end
	elseif type(n) == 'table' then
		for k,v in pairs(n) do
			self[k] = v
		end
	elseif n == nil then
	else
		error("don't know how to handle the input of type "..type(n))
	end
	self:removeExtraZeroes()
end

function BigNumber:removeExtraZeroes()
	return self
		:removeLeadingZeroes()
		:removeTrailingZeroes()
end

function BigNumber:removeLeadingZeroes()
	while true do
		self:calcMaxExp()
		assert(not self:isFinite() or self:isZero() or self[self.maxExp] ~= nil)
		if self[self.maxExp] ~= 0 then break end
		if self.repeatFrom then
			assert.ge(self.maxExp, self.repeatFrom)
			if self.maxExp == self.repeatFrom then break end
		end
		-- TODO
		-- make sure maxExp is always non-negative, and make sure minExp is always non-positive
--		if self.maxExp < 0 then break end
		self[self.maxExp] = nil
	end
	self:calcMaxExp()
	return self
end

function BigNumber:removeTrailingZeroes()
	while true do
		self:calcMinExp()
		assert(not self:isFinite() or self:isZero() or self[self.minExp] ~= nil)
		if self[self.minExp] ~= 0 then break end
		if self.repeatTo then
			assert.le(self.minExp, self.repeatTo)
			if self.minExp == self.repeatTo then break end
		end
--		if self.minExp > 0 then break end
		self[self.minExp] = nil
	end
	self:calcMinExp()
	return self
end

-- in-place
function BigNumber:carry()
	local i = 0
	local carry = 0
	while true do
		self[i] = (self[i] or 0) + carry
		carry = 0
		while self[i] >= self.base do
			self[i] = self[i] - self.base
			carry = carry + 1
		end
		i = i + 1
		if i > self.maxExp then
			if carry == 0 and (self[i] == 0 or self[i] == nil) then break end
		end
	end
	self:removeLeadingZeroes()
end

local function gcd(a,b)
	assert.ge(a, 1)
	assert.ge(b, 1)
	repeat
		a, b = b, a % b
	until b == 0
	return a
end

local function lcm(a,b)
	if a == 0 then return 0 end
	if b == 0 then return 0 end
	return a * b / gcd(a,b)
end

--[[
Thanks Sean Moore for the help on this one
--]]
function BigNumber.__add(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if b.base ~= a.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return BigNumber.constant.nan end
	if a.infinity and b.infinity then
		if a.negative == b.negative then
			return a
		else
			return BigNumber.constant.nan
		end
	elseif a.infinity then
		return a
	elseif b.infinity then
		return b
	end
	-- if signs don't match, treat it like an addition
	if a.negative ~= b.negative then
		if a.negative and not b.negative then
			return b - (-a)
		elseif b.negative and not a.negative then
			return a - (-b)
		end
	end
	if a:isZero() then return b end
	if b:isZero() then return a end

	local repeatFrom, repeatTo
	if a.repeatFrom and b.repeatFrom then
		-- just use getRepeatAsFrac() ...
		-- (ai + an / ad) + (bi + bn / bd) = ai + bi + (an * bd + bn * ad) / (ad * bd)
		-- but it's probably much slower
		a = BigNumber(a)
		b = BigNumber(b)
		-- find lcm of #aRep and #bRep
		-- stretch both to this size
		-- adjust them
		local aRepLen = a.repeatFrom - a.repeatTo + 1
		local bRepLen = b.repeatFrom - b.repeatTo + 1
		local lcmab = lcm(aRepLen, bRepLen)

		local aRepeatFrom = a.repeatFrom
		while aRepeatFrom-a.repeatTo+1 < lcmab do
			a:repeatRepeat()
		end
		a.repeatFrom = aRepeatFrom
		local bRepeatFrom = b.repeatFrom
		while bRepeatFrom-b.repeatTo+1 < lcmab do
			b:repeatRepeat()
		end
		b.repeatFrom = bRepeatFrom
		-- so now the repeat of a and b is as long as one another
		while a.repeatTo < b.repeatTo do
			b:shiftRepeat()
		end
		while b.repeatTo < a.repeatTo do
			a:shiftRepeat()
		end
		assert.eq(a.repeatFrom, b.repeatFrom)
		assert.eq(a.repeatTo, b.repeatTo)
		repeatFrom = a.repeatFrom
		repeatTo = a.repeatTo
	elseif a.repeatFrom or b.repeatFrom then
		-- just use getRepeatAsFrac() ...
		-- (ai + an / ad) + b = ai + b + an / ad
		-- but it's probably much slower
		local function adjustAtoB(a_,b_)
			if b_.minExp < a_.repeatFrom then
				a_ = BigNumber(a_)
				-- prepend
				while b_.minExp < a_.repeatFrom do
					a_:repeatRepeat()
				end
			end
			return a_
		end
		if a.repeatFrom and not b.repeatFrom then
			a = adjustAtoB(a,b)
			repeatFrom = a.repeatFrom
			repeatTo = a.repeatTo
		elseif b.repeatFrom and not a.repeatFrom then
			b = adjustAtoB(b,a)
			repeatFrom = b.repeatFrom
			repeatTo = b.repeatTo
		end
	end

	local c = BigNumber()
	c.negative = a.negative	-- == b.negative
	c.minExp = math.min(a.minExp, b.minExp)
	c.base = a.base
	local carry = 0
	local i = c.minExp
	while i <= a.maxExp or i <= b.maxExp or carry > 0 do
		local digit = (a[i] or 0) + (b[i] or 0) + carry
		carry = math.floor(digit / a.base)
		c[i] = digit % a.base
		c.maxExp = i
		i = i + 1
	end

	c.repeatFrom = repeatFrom
	c.repeatTo = repeatTo

	if c.repeatFrom or c.repeatTo then
		assert(c.repeatFrom and c.repeatTo)
		assert.le(c.repeatTo, c.repeatFrom)
		local all = true
		for j=c.repeatTo,c.repeatFrom do
			if c[j] ~= c.base-1 then
				all = false
				break
			end
		end
		-- all repeating
		if all then
			local minExp = c.repeatFrom+1
			for j=c.repeatTo,c.repeatFrom do
				c[j] = nil
			end
			c.repeatFrom = nil
			c.repeatTo = nil
			c:removeExtraZeroes()
			c = c + BigNumber{[minExp]=1, base=c.base}
		end
	end
	return c:removeExtraZeroes()
end

function BigNumber.__sub(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.base ~= b.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return BigNumber.constant.nan end
	if a.infinity and b.infinity then
		if a.negative == b.negative then
			return BigNumber.constant.nan
		else
			return a
		end
	elseif a.infinity then
		return a
	elseif b.infinity then
		return -b
	end
	-- if signs don't match, treat it like an addition
	if a.negative ~= b.negative then
		if a.negative and not b.negative then
			return -((-a) + b)
		elseif b.negative and not a.negative then
			return a + (-b)
		end
	end
	-- signs match, make sure abs biggest is 'a'
	if a.negative then
		if b < a then return -(b - a) end
	else
		if a < b then return -(b - a) end
	end
	if a:isZero() then return -b end
	if b:isZero() then return a end

	local repeatFrom, repeatTo
	if a.repeatFrom or b.repeatFrom then
		a = BigNumber(a)
		b = BigNumber(b)

		-- adjust a, the larger #, by changing 1 => 0[9]
		if not a.repeatFrom then
			-- decrement digit a.minExp
			for i=a.minExp,a.maxExp do
				a[i] = a[i] - 1
				if a[i] >= 0 then break end
				a[i] = a.base - 1
				assert.ne(i, a.maxExp)
			end

			a.minExp = a.minExp - 1
			a[a.minExp] = a.base - 1
			a.repeatFrom = a.minExp
			a.repeatTo = a.minExp
		end
		-- adjust b, the smaller #, by changing 1 => 1[0]
		if not b.repeatFrom then
			b.minExp = b.minExp - 1
			b[b.minExp] = 0
			b.repeatFrom = b.minExp
			b.repeatTo = b.minExp
		end

		-- find lcm of #aRep and #bRep
		-- stretch both to this size
		-- adjust them
		local aRepLen = a.repeatFrom - a.repeatTo + 1
		local bRepLen = b.repeatFrom - b.repeatTo + 1
		local lcmab = lcm(aRepLen, bRepLen)

		local aRepeatFrom = a.repeatFrom
		while aRepeatFrom-a.repeatTo+1 < lcmab do
			a:repeatRepeat()
		end
		a.repeatFrom = aRepeatFrom
		local bRepeatFrom = b.repeatFrom
		while bRepeatFrom-b.repeatTo+1 < lcmab do
			b:repeatRepeat()
		end
		b.repeatFrom = bRepeatFrom
		-- so now the repeat of a and b is as long as one another
		while a.repeatTo < b.repeatTo do
			b:shiftRepeat()
		end
		while b.repeatTo < a.repeatTo do
			a:shiftRepeat()
		end
		assert.eq(a.repeatFrom, b.repeatFrom)
		assert.eq(a.repeatTo, b.repeatTo)
		repeatFrom = a.repeatFrom
		repeatTo = a.repeatTo
	end

	local c = BigNumber()
	c.negative = a.negative
	c.minExp = math.min(a.minExp, b.minExp)
	c.base = a.base
	local borrow = 0	-- or -1
	local i = c.minExp
	while i <= a.maxExp or i <= b.maxExp do	-- shouldn't allow borrow past the end
		local digit = (a[i] or 0) - (b[i] or 0) + borrow
		borrow = 0
		if digit < 0 then
			borrow = -1
			digit = digit + a.base
		end
		c[i] = digit
		c.maxExp = i
		i = i + 1
	end
	-- this fails for big(2.5, 2.5) / big(5)
	assert.eq(borrow, 0)

	c.repeatFrom = repeatFrom
	c.repeatTo = repeatTo

	if c.repeatFrom or c.repeatTo then
		assert(c.repeatFrom and c.repeatTo)
		assert.le(c.repeatTo, c.repeatFrom)
		local all = true
		for j=c.repeatTo,c.repeatFrom do
			if c[j] ~= c.base-1 then
				all = false
				break
			end
		end
		-- all repeating
		if all then
			local minExp = c.repeatFrom+1
			for j=c.repeatTo,c.repeatFrom do
				c[j] = nil
			end
			c.repeatFrom = nil
			c.repeatTo = nil
			c:removeExtraZeroes()
			c = c + BigNumber{[minExp]=1, base=c.base}
		end
	end
	return c:removeExtraZeroes()
end

-- shift the repeat pattern once to the right
function BigNumber:shiftRepeat()
	if self.minExp ~= self.repeatTo then
		local tolua = require 'ext.tolua'
		error("expected self.minExp == self.repeatTo "..tolua{
			['self.minExp'] = self.minExp,
			['self.repeatFrom'] = self.repeatFrom,
			['self.repeatTo'] = self.repeatTo,
		})
	end
	local replen = self.repeatFrom - self.repeatTo + 1
	local i = self.minExp-1+replen
	if not self[i] then
		local tolua = require 'ext.tolua'
		error("failed to read from element "..i..' '..tolua{
			['self.minExp'] = self.minExp,
			['self.repeatFrom'] = self.repeatFrom,
			['self.repeatTo'] = self.repeatTo,
			replen = replen,
		})
	end
	self[self.minExp-1] = self[i]
	self.minExp = self.minExp - 1
	self.repeatFrom = self.repeatFrom - 1
	self.repeatTo = self.repeatTo - 1
	return self
end

-- duplicate the repeating pattern and adjust the repeat range
function BigNumber:repeatRepeat()
	local replen = self.repeatFrom - self.repeatTo + 1
	for i=1,replen do
		self:shiftRepeat()
	end
end

-- returns the numerator, denominator, remainder
-- all as whole integer
-- denominator is 9's for as many repeating decimals as there are, followed by 0's for how many non-repeating decimals there are,
-- such that numerator / denominator = repeating fraction
-- and numerator / denominator + remainder = original number
function BigNumber:getRepeatAsFrac()
	if not self.repeatFrom or not self.repeatTo then
		assert(not self.repeatFrom and not self.repeatTo)
		return BigNumber(0, self.base), BigNumber(1, self.base), BigNumber(self)
	end

	assert.lt(self.repeatFrom, 0)
	assert.lt(self.repeatTo, 0)
	assert.le(self.repeatTo, self.repeatFrom)

	local num = range(self.repeatTo, self.repeatFrom):mapi(function(i)
		return self[i], i-self.repeatTo
	end)
	num.base = self.base
	num = BigNumber(num)

	local denom = {}
	denom.base = self.base
	for i=self.repeatFrom,self.repeatTo,-1 do
		denom[-i-1] = 9
	end
	denom = BigNumber(denom)

	local remainder = {}
	for i=self.repeatFrom+1,self.maxExp do
		remainder[i] = self[i]
	end
	remainder.base = self.base
	remainder = BigNumber(remainder)

--DEBUG: assert.eq(num / denom + remainder, self) -- inf loop if called from div ?

	return num, denom, remainder
end

function BigNumber.__unm(a)
	a = BigNumber(a)
	a.negative = not a.negative
	return a
end

function BigNumber:shiftLeft(n)
	local result = BigNumber()
	result.base = self.base
	result.negative = self.negative
	result.nan = self.nan
	result.infinity = self.infinity
	if self.minExp then result.minExp = self.minExp + n end
	if self.maxExp then result.maxExp = self.maxExp + n end
	if self.repeatFrom then result.repeatFrom = self.repeatFrom + n end
	if self.repeatTo then result.repeatTo = self.repeatTo + n end
	for k,v in pairs(self) do
		if type(k) == 'number' then
			result[k+n] = self[k]
		end
	end
	result:removeExtraZeroes()
	return result
end
function BigNumber:shiftRight(n)
	return self:shiftLeft(-n)
end

function BigNumber:truncMinExp(newMinExp)
	local result = BigNumber(self)
	for k,v in pairs(result) do
		if type(k) == 'number' and k < newMinExp then
			result[k] = nil
		end
	end
	result.minExp = newMinExp
	result:removeExtraZeroes()
	return result
end

function BigNumber:truncMaxExp(newMaxExp)
	local result = BigNumber(self)
	for k,v in pairs(result) do
		if type(k) == 'number' and k > newMaxExp then
			result[k] = nil
		end
	end
	result.maxExp = newMaxExp
	result:removeExtraZeroes()
	return result
end


function BigNumber.simpleMul(a,b)	-- TODO better multiplication algorithm!
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.base ~= b.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return BigNumber.nan end
	if a.maxExp == nil or b.maxExp == nil then return BigNumber() end
	local bWasNegative = b.negative
	if b.negative then
		b = -b
	end

	local aMinExp = a.minExp
	local bMinExp = b.minExp
	if aMinExp ~= 0 then a = a:shiftLeft(-aMinExp) end
	if bMinExp ~= 0 then b = b:shiftLeft(-bMinExp) end

	local c = BigNumber()
	c.base = a.base
	local counter = BigNumber()
	while counter ~= b do
		c = c + a
		counter = counter + BigNumber(1)
	end
	c.negative = a.negative ~= bWasNegative

	local cMinExp = aMinExp + bMinExp
	if cMinExp ~= 0 then
		c = c:shiftRight(-cMinExp)
	end
	c:removeExtraZeroes()
	return c
end

--[[
(a + b/99..m) * (c + d/99..n)
= a * c + b/99..m * c + a + d/99..n + b*d/(99..m * 99..n)
= a * c
	+ bbb/99..m
	+ ddd/99..n
	+ bbb*ddd* numer of 1/(99..m * 99..n)

for m-rep * n-rep, the results rep is a (m+n)-rep where
m+n	(m+n)-rep
0	0
1	1
2	9
3	81
4	729
5	6561
6	59049
...
n	9^(n-1)
for a 1-rep * 1-rep, resulting rep is 9

This means, with multiplication, repeated digits increases exponentially.
So in order to implement this with any practicality, time to implement a limit on precision.

--]]
function BigNumber.longMul(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.base ~= b.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return BigNumber.constant.nan end
	if a.maxExp == nil or b.maxExp == nil then return BigNumber() end

	if a.repeatFrom and not b.repeatFrom then
		-- TODO I'm sure this is a bad implementation. I didn't think hard about edge cases.  fixme plz.
--print('before shifting repeat: '..a..' * '..b)
		local shl
		if b.minExp < 0 then
			shl = -b.minExp
			b = b:shiftLeft(shl)
			assert.eq(b.minExp, 0)
		end
		local bnumdigits = b.maxExp - b.minExp + 1
--print('rhs number of digits', bnumdigits)
		-- can I just push/pop the repeat range and otherwise multiply like a non-repeating number?
		local from = a.repeatFrom
		local to = a.repeatTo
--print('lhs repeat range', from, to)
		-- ok how about shift repeat left function
		-- and shift left by b.maxExp ... only if b.maxExp > 0 though ... hmm
		a = BigNumber(a)
		for i=1,bnumdigits do
			a:shiftRepeat()
		end
--print('after shifting repeat: '..a..' * '..b)
		a.repeatFrom = nil
		a.repeatTo = nil
		local c = a * b
--print('result without repeat', c)
		c.repeatFrom = from
		c.repeatTo = to
		while c.minExp < c.repeatTo do
			c[c.minExp] = nil
			c.minExp = c.minExp + 1
		end
		if shl then
			c = c:shiftRight(shl)
		end
		return c
	elseif not a.repeatFrom and b.repeatFrom then
		return b * a
	elseif a.repeatFrom and b.repeatFrom then
		-- (ai + an / ad) * (bi + bn / bd)
		-- = ai * (bi + bn / bd) + an / ad * (bi + bn / bd)
		-- = ai * bi + ai * (bn / bd) + (an / ad) * bi + (an / ad) * (bn / bd)
		-- = ai * bi + (ai * bn) / bd + (an * bi) / ad + (an * bn) / (ad * bd)
		local anum, adenom, aint = a:getRepeatAsFrac()
		local bnum, bdenom, bint = b:getRepeatAsFrac()
		return aint * bint
			+ (aint * bnum) / bdenom
			+ (anum * bint) / adenom
			+ (anum * bnum) / (adenom * bdenom)
	end

	local c = BigNumber()
	c.base = a.base
	c.negative = a.negative ~= b.negative

	local aMinExp = a.minExp
	local bMinExp = b.minExp
	if aMinExp ~= 0 then a = a:shiftRight(aMinExp) end
	if bMinExp ~= 0 then b = b:shiftRight(bMinExp) end

	c.minExp = 0
	c.maxExp = a.maxExp + b.maxExp
	for i=0,c.maxExp do
		c[i] = 0
	end
	-- 1) add digits
	for i=0,a.maxExp do
		local ai = a[i] or 0
		for j=0,b.maxExp do
			local bj = b[j] or 0
			c[i+j] = c[i+j] + ai * bj
		end
	end
	-- 2) apply carries
	c:carry()

	local cMinExp = aMinExp + bMinExp
	if cMinExp ~= 0 then
		c = c:shiftLeft(cMinExp)
	end
	return c:removeExtraZeroes()
end
BigNumber.__mul = BigNumber.longMul

--[[
a^b = (an..a0)^b = (10^n an + ... + 10^0 a0)^b
 = (10^n an + ... + a0) * ... * (10^n an + ... + a0) (b times)

(a0 + x a1)^2 = a0^2 + x 2 a0 a1 + x^2 a1^2
(a0 + x a1 + x^2 a2)^2 =
	a0 a0
	+ x * ((a0 + a1) + (a1 + a0))
	+ x^2 * (a0 a2 + a1 a1 + a2 a0)
	+ x^3 * (a2 a1 + a1 a2)
	+ x^4 * a2 a2

this one turns out to run *really* slow,
and breaks for n > 80 or so
--]]
function BigNumber.intPow_binomial(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.base ~= b.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return BigNumber.constant.nan end
	if b:isZero() then return BigNumber{[0]=1, base=a.base} end
	if a:isZero() then return BigNumber{base=a.base} end
	if b.negative then error('no support for negative powers!') end
	--'b' iterators cycling through the coefficients of 'a'
	-- so iters will be sized by a's digits and increment until it reaches b
	local iters = table()
	local bNum = b:tonumber()	-- for our iterator table ... this limits this power method to only exponents of 4billion (or whatever the precision of doubles is)
	for i=1,bNum do
		iters[i] = a.minExp
	end
	local c = BigNumber()
	c.base = a.base
	c.minExp = a.minExp ^ bNum
	c.maxExp = a.maxExp ^ bNum
	if b:intmod(BigNumber(2)) == 1 then c.negative = a.negative end
	local done = false
	while not done do
		--apply iteration
		local power = 0
		local coeff = 1
		for i=1,bNum do
			power = power + iters[i]
			coeff = coeff * a[iters[i]]
		end
		c[power] = (c[power] or 0) + coeff

		for i=1,bNum do
			iters[i] = iters[i] + 1
			if iters[i] <= a.maxExp then break end
			iters[i] = a.minExp
			if i == bNum then
				done = true
			end
		end
	end
	c:carry()
	return c
end
function BigNumber.intPow_simple(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.base ~= b.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return BigNumber.constant.nan end
	if b.negative then error('no support for negative powers!') end
	local c = BigNumber(1)
	c.base = a.base
	local counter = BigNumber()
	while counter ~= b do
		c = c * a
		counter = counter + BigNumber(1)
	end
	return c
end

-- returns a table of digits, with the 0's digit in [0], and the n'th digit in [n]
function BigNumber.toBase(n, base)
	assert.gt(base, 1, "can't set to a base of 1 or lower")

	n = BigNumber(n)
	if n.base == base then return n end

--	assert.eq(n.minExp, 0, "can only handle integers, but got a minExp "..n.minExp)	-- can't handle fractions yet

	-- construct 'p' as value 'base' in base 'n.base'
	-- don't rely on 'toBase' to convert it -- or we'll get a recursive call
	local p = BigNumber(base, n.base)

	local orign = n
	local result = BigNumber()
	result.base = base
	result.minExp = 0
	do
		local i = 0
		local b
		local zero = BigNumber(0, n.base)
		while n > zero do
			n, b = n:intdiv(p)
			n = n:truncMinExp(0)
			result[i] = b:tonumber()
			result.maxExp = i
			i = i + 1
		end
	end

	local f = orign:truncMaxExp(-1)
	if f > 0 then
		local i = -1
		f = f * p
		while f > 0 do
			result[i] = (f:truncMinExp(0) % p):tonumber()
			f = (f - result[i]) * p
			result.minExp = i
			i = i - 1
		end
	end

	result:removeExtraZeroes()

	return result
end

--[[
-- convert a table of digits back to a number
-- use caching of powers-of-two to do this
local function bin2num(b)
	local res = BigNumber()
	local _2d = BigNumber(1)	--2^d
	for d,ni in ipairs(b) do	--d=1 == 0's digit
		if ni == 1 then
			res = res + _2d
		end
		_2d = _2d * 2
	end
	return res
end
--]]

-- a^b = a^sum_i b_i for b_i the powers-of-two of b
-- TODO use another power? other than 2?
-- this runs faster than intPow_simple for n>80 (and they're both quick enough for n<=80)
function BigNumber.intPow_binDigits(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.nan or b.nan then return BigNumber.constant.nan end
	if b.negative then error('no support for negative powers!') end
	if b:isZero() then return BigNumber(1) end
	if a:isZero() then return BigNumber(0) end

	local bb = b:toBase(2)	-- binary form of 'b'. TODO arbitrary base bignumbers
	local res = BigNumber(1):toBase(a.base)
	local _2ToTheI = BigNumber(1):toBase(a.base)	-- 2^i
	local aToThe2ToThei = a	--a^(2^i)
	for i=0,bb.maxExp do
	--for i=0,#bb do
		local bbi = bb[i]
		if bbi == 1 then
			res = res * aToThe2ToThei
		end
		aToThe2ToThei = aToThe2ToThei * aToThe2ToThei
		_2ToTheI = _2ToTheI + _2ToTheI
	end
	return res
end

-- simple runs *very* fast for b<100
-- binomial *might* run faster later ... but probably not
-- binDigits runs faster than simple for b>80
-- ... and should be easily extendible to decimal numbers
--BigNumber.intpow = BigNumber.intPow_binomial
--BigNumber.intpow = BigNumber.intPow_simple
BigNumber.intpow = BigNumber.intPow_binDigits
BigNumber.__pow = BigNumber.intpow

function BigNumber.simpleIntDiv(a,b)	-- TODO negative support!
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.base ~= b.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return BigNumber.constant.nan end
	if b.maxExp == nil then return BigNumber.constant.nan end
	local n = BigNumber(b)
	local c = BigNumber(0, a.base)
	c.base = a.base
	while n <= a do
		n = n + b
		c = c + BigNumber(1, a.base)
	end
	n = b - (n - a)
	return c, n
end

function BigNumber.longIntDiv(a,b, getRepeatingDecimals)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if b.base ~= a.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return BigNumber.constant.nan end
	if b:isZero() then
		if a:isZero() then return BigNumber.constant.nan end
		if a.negative then
			return -BigNumber.constant.infinity
		else
			return BigNumber.constant.infinity
		end
	end
	if a:isZero() then return a end

	-- repeated decimal division ...
	if a.repeatFrom and not b.repeatFrom then
		-- (ai + an / ad) / b
		-- = (ai ad + an) / (b ad)
		local anum, adenom, aint = a:getRepeatAsFrac()
		assert(not anum.repeatFrom)	-- or else infinite recursion ...
		assert(not adenom.repeatFrom)
		return (aint * adenom + anum) / (b * adenom)
	elseif not a.repeatFrom and b.repeatFrom then
		-- a / (bi + bn / bd)
		-- = a / ((b bd + bn) / bd)
		-- = (a bd) / (b bd + bn)
		local bnum, bdenom, bint = b:getRepeatAsFrac()
		assert(not bnum.repeatFrom)	-- or else infinite recursion ...
		assert(not bdenom.repeatFrom)
		return (a * bdenom) / (bint * bdenom + bnum)
	elseif a.repeatFrom and b.repeatFrom then
		-- (ai + an/ad) / (bi + bn/bd)
		-- = (ai ad + an)/ad / ((bi bd + bn)/bd)
		-- = ((ai ad + an) * bd) / ((bi bd + bn) * ad)
		local anum, adenom, aint = a:getRepeatAsFrac()
		local bnum, bdenom, bint = b:getRepeatAsFrac()
		assert(not anum.repeatFrom)	-- or else infinite recursion ...
		assert(not adenom.repeatFrom)
		assert(not bnum.repeatFrom)
		assert(not bdenom.repeatFrom)
		return ((aint * adenom + anum) * bdenom) / ((bint * bdenom + bnum) * adenom)
	end

	-- TODO decimal division, especially picking the # of digits of accuracy
	--local aMinExp = 0
	-- [[
	local aMinExp = a.minExp
	if aMinExp < 0 then a = a:shiftRight(aMinExp) else aMinExp = 0 end
	--]]
	-- [[
	local bMinExp = b.minExp
	if bMinExp < 0 then b = b:shiftRight(bMinExp) else bMinExp = 0 end
	--]]

	local dividendDigits = BigNumber(a)
	local place = dividendDigits.maxExp or 0
	local dividendCurrentDigits = BigNumber(dividendDigits[place], a.base)	-- most significant digit
	dividendDigits[place] = nil
	dividendDigits:calcMaxExp()
	local divisor = BigNumber(b)
	local results = BigNumber(0, a.base)
	results.maxExp = place
	results.minExp = place
	local remainder = BigNumber(0, a.base)
	local digitLastRemainderPairs = table()	-- keys are digit,lastRemainder as strings
	local repeatFrom, repeatTo
	for precision=1,math.huge do
		local digit = dividendCurrentDigits:simpleIntDiv(divisor)
		dividendCurrentDigits = dividendCurrentDigits - divisor * digit
		local lastRemainder = remainder
		remainder = dividendCurrentDigits
		--print('place',place,'digit', digit, 'last remainder', lastRemainder)
		local key = digit..','..lastRemainder
		if place < 0 then
			if digitLastRemainderPairs[key] then
				repeatFrom = digitLastRemainderPairs[key]
				repeatTo = place + 1
				break
			end
			digitLastRemainderPairs[key] = place
		end
		assert.ge(dividendCurrentDigits, BigNumber(0, a.base))
		results[place] = digit[0] or 0
		results.minExp = place
		if place <= 0 and not getRepeatingDecimals then
			-- return int part and remainder
			break
		end
		place = place - 1
		dividendCurrentDigits = dividendCurrentDigits * BigNumber(a.base, a.base) + BigNumber(dividendDigits[place], a.base)
	end
	-- final check that we're not repeating zeroes ...
	if repeatFrom then
		if repeatFrom == repeatTo and results[repeatFrom] == 0 then
			results[repeatFrom] = nil
			results.minExp = repeatFrom+1
			repeatFrom = nil
			repeatTo = nil
		end
	end
	-- and we're good to test
	--[[ TODO this looks unfinished ...
	-- TODO define maxRepeatingLength somewhere
	-- and do something with maxRepeatingP
	if repeatFrom then
		local repeatingLength = repeatFrom - repeatTo + 1
		if not maxRepeatingLength or repeatingLength > maxRepeatingLength then
			maxRepeatingLength = repeatingLength
			maxRepeatingP = b
		end
	end
	--]]
	results:removeLeadingZeroes()
	if getRepeatingDecimals then
		results.repeatFrom = repeatFrom
		results.repeatTo = repeatTo
	end

	-- [[
	local cMinExp = aMinExp - bMinExp
	if cMinExp ~= 0 then
		results = results:shiftLeft(cMinExp)
	end
	results:removeExtraZeroes()
	--]]

	results.negative = a.negative ~= b.negative

	return results, dividendCurrentDigits
end

function BigNumber.decdiv(a,b)
	return (BigNumber.longIntDiv(a,b,true))
end

BigNumber.intdiv = BigNumber.longIntDiv

function BigNumber.__div(a,b)
	--return (BigNumber.intdiv(a,b))	-- return only one value
	return (BigNumber.decdiv(a,b))	-- return only one value
end

function BigNumber.intmod(a,b)
	return select(2, BigNumber.intdiv(a,b))
end

BigNumber.__mod = BigNumber.intmod

function BigNumber:isPosInf()
	return self.infinity and not self.negative
end

function BigNumber:isNegInf()
	return self.infinity and self.negative
end

function BigNumber.__eq(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.base ~= b.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return false end
	if a.maxExp == nil and b.maxExp == nil then return true end	--negative zero equals positive zero
	if a.maxExp ~= b.maxExp then return false end
	if a.minExp ~= b.minExp then return false end
	if a.negative ~= b.negative then return false end
	if a.repeatFrom ~= b.repeatFrom then return false end
	if a.repeatTo ~= b.repeatTo then return false end
	for i=a.minExp,a.maxExp do
		if a[i] ~= b[i] then return false end
	end
	return true
end
function BigNumber.__ne(a,b)
	return not a == b
end
function BigNumber.__lt(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.nan or b.nan then return false end
	if b:isPosInf() and not a:isPosInf() then return true end
	if a:isZero() then
		if b:isZero() then
			return false
		end
		return not b.negative
	end
	if b:isZero() then
		return a.negative
	end
	if a.negative and not b.negative then return true end
	if b.negative and not a.negative then return false end

	-- do zero/negative tests before converting number bases
	if a.base ~= b.base then b = b:toBase(a.base) end

	if a.maxExp < b.maxExp then return not a.negative end
	if a.maxExp > b.maxExp then return a.negative end
	for i=a.maxExp,math.min(a.minExp, b.minExp),-1 do
		if (a[i] or 0) < (b[i] or 0) then return not a.negative end
		if (a[i] or 0) > (b[i] or 0) then return a.negative end
	end
	return false -- equal
end
function BigNumber.__le(a,b)
	if not BigNumber:isa(a) then a = BigNumber(a) end
	if not BigNumber:isa(b) then b = BigNumber(b) end
	if a.base ~= b.base then b = b:toBase(a.base) end
	if a.nan or b.nan then return false end
	if b:isPosInf() then return true end
	if a:isZero() then
		if b:isZero() then
			return true
		end
		return not b.negative
	end
	if b:isZero() then
		return a.negative
	end
	if a.negative and not b.negative then return true end
	if b.negative and not a.negative then return false end
	if a.maxExp < b.maxExp then return not a.negative end
	if a.maxExp > b.maxExp then return a.negative end
	for i=a.maxExp,math.min(a.minExp, b.minExp),-1 do
		if (a[i] or 0) < (b[i] or 0) then return not a.negative end
		if (a[i] or 0) > (b[i] or 0) then return a.negative end
	end
	return true -- equal
end
function BigNumber.tonumber(n)
	if n.nan then return 0/0 end
	if n:isZero() then return 0 end
	local sum = 0
	for i=n.maxExp,0,-1 do
		sum = sum * n.base
		sum = sum + (n[i] or 0)
	end
	if n.negative then sum = -sum end
	return sum
end
--[[
	number of digits that the integer form has
	zero has 0 digits
	nan has nan digits
	negatives are ignored
	decimal places are ignored
--]]
function BigNumber.calcMaxExp(n)
	if n.nan then return n end
	n.maxExp = nil
	for k,v in pairs(n) do
		if type(k) == 'number' then
			if n.maxExp == nil then
				n.maxExp = k
			else
				n.maxExp = math.max(n.maxExp, k)
			end
		end
	end
	if n.maxExp == nil then n.minExp = nil end
	return n.maxExp
end
function BigNumber.calcMinExp(n)
	if n.nan then return n end
	n.minExp = nil
	for k,v in pairs(n) do
		if type(k) == 'number' then
			if n.minExp == nil then
				n.minExp = k
			else
				n.minExp = math.min(n.minExp, k)
			end
		end
	end
	if n.minExp == nil then n.maxExp = nil end
	return n.minExp
end
function BigNumber.__tostring(n)
	if n.nan then return 'nan' end
	if n.infinity then
		if n.negative then
			return '-inf'
		else
			return 'inf'
		end
	end
	if n:isZero() then return '0' end
	local s = ''
	for i=n.maxExp,0,-1 do
		if i == n.repeatFrom then s = s .. '[' end
		if n[i] then
			s = s .. number.charfor(math.floor(n[i]))
		else
			-- repeated decimals starting at positive digits means repeating through
			if n.repeatTo and i < n.repeatTo then
				local j = (i - n.repeatTo) % (n.repeatFrom - n.repeatTo + 1) + n.repeatTo
				s = s .. number.charfor(math.floor(n[j] or 0))
			else
				s = s .. number.charfor(0)
			end
		end
		if i == n.repeatTo then s = s .. ']' end
	end
	if n.minExp < 0 then
		s = s .. '.'
		for i=-1,n.minExp,-1 do
			if i == n.repeatFrom then s = s .. '[' end
			s = s .. number.charfor(math.floor(n[i] or 0))
			if i == n.repeatTo then s = s .. ']' end
		end
	end
	if n.negative then s = '-'..s end
	if n.base ~= 10 then
		s = s .. '_'..n.base
	end
	if n.trailsoff then
		s = s .. '...'
	end
	return s
end
function BigNumber.__concat(a,b)
	return tostring(a) .. tostring(b)
end

function BigNumber.factorial(n)
	local f = BigNumber(1)
	local i = BigNumber(2)
	while i <= n do
		f = f * i
		i = i + BigNumber(1)
	end
	return f
end

-- returns a lua table, indexes 1-n, of all the digits
-- t[1] is the 1's place, t[2] is the 10's place, etc
function BigNumber:digits()
	local t = table()
	if self:isZero() then return t end
	for i=math.min(0,self.minExp or 0),math.max(0,self.maxExp or 0) do
		t:insert(self[i] or 0)
	end
	return t
end
-- ridiculously often used
function BigNumber:sumOfDigits()
	return self:digits():sum()
end

function BigNumber:isZero()
	return not self.maxExp
end

function BigNumber:isFinite()
	if self.nan then return false end
	if self.infinity then return false end
	return true
end

-- hmm, need a different name
BigNumber.constant = {
	nan = BigNumber{nan=true},
	infinity = BigNumber{infinity=true},
}

-- https://en.wikipedia.org/wiki/Integer_square_root#Algorithm_using_binary_search
function BigNumber.sqrt(x)
--DEBUG(BigNumber.sqrt): local tolua = require 'ext.tolua'
--DEBUG(BigNumber.sqrt): print'BigNumber.sqrt'
	if not BigNumber:isa(x) then x = BigNumber(x) end
--DEBUG(BigNumber.sqrt): print('x =', x)
	x = x:floor()
--DEBUG(BigNumber.sqrt): print('floor: x =', x)
	local l = BigNumber()
--DEBUG(BigNumber.sqrt): print('l =', l)
	local r = x + 1
--DEBUG(BigNumber.sqrt): print('r =', r)
	while l ~= r - 1 do
--DEBUG(BigNumber.sqrt): print()
--DEBUG(BigNumber.sqrt): print('r = '..r)
--DEBUG(BigNumber.sqrt): print('tolua(r) = '..tolua(r))
--DEBUG(BigNumber.sqrt): print('r - 1 = '..(r - 1))
--DEBUG(BigNumber.sqrt): print('tolua(r - 1) = '..tolua(r - 1))
--DEBUG(BigNumber.sqrt): print('l = '..l)
--DEBUG(BigNumber.sqrt): print('tolua(l) = '..tolua(l))
--DEBUG(BigNumber.sqrt): print('l == r - 1 = '..tostring(l == r - 1))
--DEBUG(BigNumber.sqrt): print('l ~= r - 1 = '..tostring(l ~= r - 1))
		local m = (l + r):intdiv(2)
		if m * m <= x then
			l = m
		else
			r = m
		end
	end
	return l
end

--[[
this is also in ext.math
but it has a few floor() functions there
and it has separated % and / operations
would be nice if Lua would let us access int division-with-remainder
--]]
function BigNumber.primeFactorization(n)
--DEBUG(BigNumber.primeFactorization): print('BigNumber.primeFactorization')
--DEBUG(BigNumber.primeFactorization): print('n =', n)
	n = BigNumber(n)
--DEBUG(BigNumber.primeFactorization): print('n =', n)
	local f = table()
	while n > 1 do
		local found = false
		local i = BigNumber(2)
--DEBUG(BigNumber.primeFactorization): print('i =', i)
		local sqrtn = n:sqrt()
--DEBUG(BigNumber.primeFactorization): print('sqrtn =', sqrtn)
		while i <= sqrtn do
			local a, b = n:intdiv(i)
--DEBUG(BigNumber.primeFactorization): print(i, 'intdiv', a, b)
			if b:isZero() then
				n = a
				f:insert(i)
				found = true
				break
			end
			i = i + 1
		end
		if not found then
--DEBUG(BigNumber.primeFactorization): print('inserting', n)
			f:insert(n)
			break
		end
	end
	return f
end

function BigNumber.floor(n)
	n = BigNumber(n)
	while n.minExp < 0 do
		n[n.minExp] = 0
		n.minExp = n.minExp + 1
		if n.minExp > n.maxExp then return BigNumber() end
	end
	return n
end

return BigNumber
