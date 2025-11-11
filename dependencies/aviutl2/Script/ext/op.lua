--[[
make lua functions for each operator.
it looks like i'm mapping 1-1 between metamethods and fields in this table.
useful for using Lua as a functional language.

TODO rename to 'ops'?
--]]

--local load = require 'string'.load	-- string.load = loadstring or load
local load = loadstring or load

-- test if we hae lua 5.3 bitwise operators
-- orrr I could just try each op and bail out on error
-- and honestly I should be defaulting to the 'bit' library anyways, esp in the case of luajit where it is translated to an asm opcode
local lua53 = _VERSION >= 'Lua 5.3'

local symbolscode = [[

	-- which fields are unary operators
	local unary = {
		unm = true,
		bnot = true,
		len = true,
		lnot = true,
	}

	local symbols = {
		add = '+',
		sub = '-',
		mul = '*',
		div = '/',
		mod = '%',
		pow = '^',
		unm = '-',			-- unary
		concat = '..',
		eq = '==',
		ne = '~=',
		lt = '<',
		le = '<=',
		gt = '>',
		ge = '>=',
		land = 'and',		-- non-overloadable
		lor = 'or',			-- non-overloadable
		len = '#',			-- unary
		lnot = 'not',		-- non-overloadable, unary
]]
if lua53 then
	symbolscode = symbolscode .. [[
		idiv = '//',		-- 5.3
		band = '&',			-- 5.3
		bor = '|',			-- 5.3
		bxor = '~',			-- 5.3
		shl = '<<',			-- 5.3
		shr = '>>',			-- 5.3
		bnot = '~',			-- 5.3, unary
]]
--[[ alternatively, luajit 'bit' library:
I should probably include all of these instead
would there be a perf hit from directly assigning these functions to my own table,
 as there is a perf hit for assigning from ffi.C func ptrs to other variables?  probably.
 how about as a tail call / vararg forwarding?
I wonder if luajit adds extra metamethods

luajit 2.0		lua 5.2		lua 5.3
band			band		&
bnot			bnot		~
bor				bor			|
bxor			bxor		~
lshift			lshift		<<
rshift			rshift		>>
arshift			arshift
rol				lrotate
ror				rrotate
bswap (reverses 32-bit integer endian-ness of bytes)
tobit (converts from lua number to its signed 32-bit value)
tohex (string conversion)
				btest (does some bitflag stuff)
				extract (same)
				replace (same)
--]]
end
symbolscode = symbolscode .. [[
	}
]]

local symbols, unary = assert(load(symbolscode..' return symbols, unary'))()

local code = symbolscode .. [[
	-- functions for operators
	local ops
	ops = {
]]
for name,symbol in pairs(symbols) do
	if unary[name] then
		code = code .. [[
		]]..name..[[ = function(a) return ]]..symbol..[[ a end,
]]
	else
		code = code .. [[
		]]..name..[[ = function(a,b) return a ]]..symbol..[[ b end,
]]
	end
end
code = code .. [[
		index = function(t, k) return t[k] end,
		newindex = function(t, k, v)
			t[k] = v
			return t, k, v	-- ? should it return anything ?
		end,
		call = function(f, ...) return f(...) end,

		symbols = symbols,

		-- special pcall wrapping index, thanks luajit.  thanks.
		-- while i'm here, multiple indexing, so it bails out nil early, so it's a chained .? operator
		safeindex = function(t, ...)
			if select('#', ...) == 0 then return t end
			local res, v = pcall(ops.index, t, ...)
			if not res then return nil, v end
			return ops.safeindex(v, select(2, ...))
		end,
	}
	return ops
]]
return assert(load(code))()
