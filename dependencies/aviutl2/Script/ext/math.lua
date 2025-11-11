local math = {}
for k,v in pairs(require 'math') do math[k] = v end

math.nan = 0/0

math.e = math.exp(1)

-- luajit and lua 5.1 compat ...
if not math.atan2 then math.atan2 = math.atan end
-- also note, code that uses math.atan(y,x) in luajit will instead just call math.atan(y) ...

-- some luas don't have hyperbolic trigonometric functions

if not math.sinh then
	function math.sinh(x)
		local ex = math.exp(x)
		return .5 * (ex - 1/ex)
	end
end

if not math.cosh then
	function math.cosh(x)
		local ex = math.exp(x)
		return .5 * (ex + 1/ex)
	end
end

if not math.tanh then
	function math.tanh(x)
--[[ this isn't so stable.
		local ex = math.exp(x)
		return (ex - 1/ex) / (ex + 1/ex)
--]]
-- [[ instead...
-- if e^-2x < smallest float epsilon
-- then consider (e^x  - e^-x) ~ e^x .. well, it turns out to be 1
-- and if e^2x < smallest float epsilon then -1
		if x < 0 then
			local e2x = math.exp(2*x)
			return (e2x - 1) / (e2x + 1)
		else
			local em2x = math.exp(-2*x)
			return (1 - em2x) / (1 + em2x)
		end
--]]
	end
end

function math.asinh(x)
	return math.log(x + math.sqrt(x*x + 1))
end

function math.acosh(x)
	return math.log(x + math.sqrt(x*x - 1))
end

function math.atanh(x)
	return .5 * math.log((1 + x) / (1 - x))
end

function math.cbrt(x)
	return math.sign(x) * math.abs(x)^(1/3)
end

function math.clamp(v,min,max)
	return math.max(math.min(v, max), min)
end

function math.sign(x)
	if x < 0 then return -1 end
	if x > 0 then return 1 end
	return 0
end

function math.trunc(x)
	if x < 0 then return math.ceil(x) else return math.floor(x) end
end

function math.round(x)
	return math.floor(x+.5)
end

function math.isnan(x) return x ~= x end
function math.isinf(x) return x == math.huge or x == -math.huge end
function math.isfinite(x) return tonumber(x) and not math.isnan(x) and not math.isinf(x) end

function math.isprime(n)
	if n < 2 then return false end	-- 1 isnt prime
	for i=2,math.floor(math.sqrt(n)) do
		if n%i == 0 then
			return false
		end
	end
	return true
end

-- assumes n is a non-negative integer.  this isn't the Gamma function
function math.factorial(n)
	local prod = 1
	for i=1,n do
		prod = prod * i
	end
	return prod
end

function math.factors(n)
	local table = require 'ext.table'
	local f = table()
	for i=1,n do
		if n%i == 0 then
			f:insert(i)
		end
	end
	return f
end

-- returns a table containing the prime factorization of the number
function math.primeFactorization(n)
	local table = require 'ext.table'
	n = math.floor(n)
	local f = table()
	while n > 1 do
		local found = false
		for i=2,math.floor(math.sqrt(n)) do
			if n%i == 0 then
				n = math.floor(n/i)
				f:insert(i)
				found = true
				break
			end
		end
		if not found then
			f:insert(n)
			break
		end
	end
	return f
end

function math.gcd(a,b)
	return b == 0 and a or math.gcd(b, a % b)
end

-- if this math lib gets too big ...
function math.mix(a,b,s)
	return a * (1 - s) + b * s
end

return math
