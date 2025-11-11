-- number metatable
local math = require 'ext.math'

-- flatten ext.math into this just like we flatten lua's string into ext.string

local hasutf8, utf8 = pcall(require, 'utf8')

local number = {}
number.__index = number

for k,v in pairs(math) do number[k] = math[k] end

-- [[ tostring machine precision of arbitrary base
number.alphabets = {
	{('0'):byte(), ('9'):byte()},	-- latin numbers
	{('a'):byte(), ('z'):byte()},	-- english
	{0x3b1, 0x3c9},	-- greek
	{0x430, 0x45f},	-- cyrillic
	{0x561, 0x586},	-- armenian
	{0x905, 0x939},	-- devanagari
	{0x2f00, 0x2fd5},	-- kangxi
	{0x3041, 0x3096},	-- hiragana
	{0x30a1, 0x30fa},	-- katakana
	{0x4e00, 0x9fd0},	-- chinese, japanese, korean characters
}

function number.charfor(digit)
	local table = require 'ext.table'
	for _,alphabet in ipairs(number.alphabets) do
		local start,fin = table.unpack(alphabet)
		if digit <= fin - start then
			digit = digit + start
			if hasutf8 then
				return utf8.char(digit)
			else
				-- TODO this will fail with utf8 chars beyond ascii
				return string.char(digit)
			end
		end
		digit = digit - (fin - start + 1)
	end
	error 'you need more alphabets to represent that many digits'
end

-- TODO rename above function to 'tochar' ?
function number.todigit(ch)
	local table = require 'ext.table'
	local indexInAlphabet
	if hasutf8 then
		indexInAlphabet = utf8.codepoint(ch)
	else
		-- TODO this will fail with utf8 chars beyond ascii
		indexInAlphabet = string.byte(ch)
	end
	local lastTotalIndex = 0
	for _,alphabet in ipairs(number.alphabets) do
		local start,fin = table.unpack(alphabet)
		if indexInAlphabet >= start and indexInAlphabet <= fin then
			return lastTotalIndex + (indexInAlphabet - start)
		end
		lastTotalIndex = lastTotalIndex + (fin - start + 1)
	end
	error"couldn't find the character in all the alphabets"
end

number.base = 10
number.maxdigits = 50
-- I'm not going to set this as __tostring by default, but I will leave it as part of the meta
-- feel free to use it with a line something like (function(m)m.__tostring=m.tostring end)(debug.getmetatable(0))
number.tostring = function(t, base, maxdigits)
	local s = {}
	if t < 0 then
		t = -t
		table.insert(s, '-')
	end
	if t == 0 then
		table.insert(s, '0')
		table.insert(s, '.')
	else
		if not base then base = number.base end
		if not maxdigits then maxdigits = number.maxdigits end
		local i = math.floor(math.log(t,base))+1
		if i == math.huge then error'infinite number of digits' end
		t = t / base^i
		local dot
		while true do
			if i < 1 then
				if not dot then
					dot = true
					table.insert(s, '.')
					if i < 0 then
						table.insert(s, ('0'):rep(-i))
					end
				end
				if t == 0 then break end
				if i <= -maxdigits then break end
			end
			t = t * base
			local digit = math.floor(t)
			t = t - digit
			-- at this point 'digit' holds an integer value in [0,base)
			-- there's two ways we can go about it:
			-- 1) traditionally, where each digit represents an integer, times some base^i for some i
			--	in this case, for fractional bases, the last (fraction) digit is only considered up to the fraction
			--  i.e. in base 2.5, from 0b..1b we have the span of 1, partitioned by digits 0.1b at 2/5ths the distance, 0.2b at 4/5ths the distance, and 1.0 and 5/5ths the distance
			-- 		from 1b..2b we have the same span,
			--		and from 2b..10b we have half that span, from 2 to 2.5
			--		therefore the value 2.2b would represent 2.8, which is also represetned by 10.0112210002002002...
			-- so as long as the span between digits can represent fractions of base^i rather than whole base^i's
			--	we can have multiple representations of numbers
			-- 2) stretched.  in this case a fractional span, such as from 2b to 10b in base 2.5, would be stretched
			-- 		this is harder to convey when descibing the number system
			table.insert(s, number.charfor(digit))
			i = i - 1
		end
	end
	if s[#s] == '.' then s[#s] = nil end
	return table.concat(s)
end
--]]

-- ('a'):byte():char() == 'a'
number.char = string.char

-- me getting tired of typing 'tostring(16)' etc ...
function number.bin(x) return number.tostring(x, 2) end
function number.oct(x) return number.tostring(x, 8) end
function number.hex(x) return number.tostring(x, 16) end

-- so the lookup goes: primitive number -> number -> math
return number
