local assert = require 'ext.assert'
local table = require 'ext.table'
local string = require 'ext.string'
local coroutine = require 'ext.coroutine'
local number = require 'ext.number'
local op = require 'ext.op'

-- fix up lua type metatables

-- nil
debug.setmetatable(nil, {__concat = string.concat})

-- booleans
debug.setmetatable(true, {
	__concat = string.concat,
	__index = {
		and_ = op.land,
		or_ = op.lor,
		not_ = op.lnot,
		xor = function(a,b) return a ~= b end,
		implies = function(a,b) return not a or b end,
	}
})

debug.setmetatable(0, number)

-- strings
getmetatable('').__concat = string.concat
getmetatable('').__index = string

-- It'd be fun if I could apply the operator to all return values, and not just the first ...
-- like (function() return 1,2 end + function() return 3,4 end)() returns 4,6
local function combineFunctionsWithBinaryOperator(f, g, opfunc)
	if type(f) == 'function' and type(g) == 'function' then
		return function(...)
			return opfunc(f(...), g(...))
		end
	elseif type(f) == 'function' then
		return function(...)
			return opfunc(f(...), g)
		end
	elseif type(g) == 'function' then
		return function(...)
			return opfunc(f, g(...))
		end
	else
		-- shouldn't get called unless __add is called explicitly
		return function()
			return opfunc(f, g)
		end
	end
end

-- primitive functions.  should these be public?  or put in a single table?

-- function operators generate functions
-- f(x) = y, g(x) = z, (f+g)(x) = y+z
local functionMeta = {
	-- I could make this a function composition like the rest of the meta operations,
	-- but instead I'm going to have it follow the default __concat convention I have with other primitive types
	__concat = string.concat,
	dump = function(f) return string.dump(f) end,
	__add = function(f, g) return combineFunctionsWithBinaryOperator(f, g, op.add) end,
	__sub = function(f, g) return combineFunctionsWithBinaryOperator(f, g, op.sub) end,
	__mul = function(f, g) return combineFunctionsWithBinaryOperator(f, g, op.mul) end,
	__div = function(f, g) return combineFunctionsWithBinaryOperator(f, g, op.div) end,
	__mod = function(f, g) return combineFunctionsWithBinaryOperator(f, g, op.mod) end,
	__pow = function(f, g) return combineFunctionsWithBinaryOperator(f, g, op.pow) end,
	__unm = function(f) return function(...) return -f(...) end end,
	__len = function(f) return function(...) return #f(...) end end,
	-- boolean operations aren't overloaded just yet.  should they be?
	--__call doesn't work anyways
	-- TODO comparison operations too?  probably not equivalence for compatability with sort() and anything else
	-- TODO boolean operations?  anything else?
	--[[ here's one option for allowing any function object dereference to be mapped to a new function
	__index = function(f, k) return function(...) return f(...)[k] end end,
	__newindex = function(f, k, v) return function(...) f(...)[k] = v end end,
	--]]
	-- [[ ... but that prevents us from overloading our own methods.
	-- so here's "index" to be used in its place
	-- while we can provide more of our own methods as we desire
	__index = {
		-- takes a function that returns an object
		--  returns a function that returns that object's __index to the key argument
		-- so if f() = {a=1} then f:index'a'() == 1
		index = function(f, k)
			return function(...)
				return f(...)[k]
			end
		end,
		-- takes a function that returns an object
		--  returns a function that applies that object's __newindex to the key and value arguments
		-- so if t={} and f()==t then f:assign('a',1)() assigns t.a==1
		assign = function(f, k, v)
			return function(...)
				f(...)[k] = v
			end
		end,

		-- f:compose(g1, ...) returns a function that evaluates to f(g1(...(gn(args))))
		compose = function(...)
			local funcs = table.pack(...)
			for i=1,funcs.n do
				assert.type(funcs[i], 'function')
			end
			return function(...)
				local args = table.pack(...)
				for i=funcs.n,1,-1 do
					args = table.pack(funcs[i](table.unpack(args,1,args.n)))
				end
				return table.unpack(args,1,args.n)
			end
		end,

		-- f:compose_n(n, g) returns a function that evaluates to f(arg[1], ... arg[j-1], g(arg[j]), arg[j+1], ..., arg[n])
		compose_n = function(f, n, ...)
			local funcs = table.pack(...)
			return function(...)
				local args = table.pack(...)

				local ntharg = {args[n]}
				ntharg.n = n <= args.n and 1 or 0
				for i=funcs.n,1,-1 do
					ntharg = table.pack(funcs[i](table.unpack(ntharg,1,ntharg.n)))
				end

				args[n] = ntharg[1]
				args.n = math.max(args.n, n)

				return f(table.unpack(args, 1, args.n))
			end
		end,

		-- bind / partial apply -- currying first args, and allowing vararg rest of args
		bind = function(f, ...)
			local args = table.pack(...)
			return function(...)
				local n = args.n
				local callargs = {table.unpack(args, 1, n)}
				for i=1,select('#', ...) do
					n=n+1
					callargs[n] = select(i, ...)
				end
				return f(table.unpack(callargs, 1, n))
			end
		end,

		-- bind argument n, n+1, n+2, ... to the values provided
		bind_n = function(f, n, ...)
			local nargs = table.pack(...)
			return function(...)
				local args = table.pack(...)
				for i=1,nargs.n do
					args[n+i-1] = nargs[i]
				end
				args.n = math.max(args.n, n+nargs.n-1)
				return f(table.unpack(args, 1, args.n))
			end
		end,

		-- Takes a function and a number of arguments,
		-- returns a function that applies them individually,
		-- first to the function, then to each function returned
		-- (a1 -> (a2 -> ... (an -> b))) -> (a1, a2, ..., an -> b)
		uncurry = function(f, n)
			return function(...)
				local s = f
				for i=1,n do
					s = s(select(i, ...))
				end
				return s
			end
		end,
		-- grows/shrinks the number of args passed.  pads with nil.
		nargs = function(f, n)
			return function(...)
				local t = {}
				for i=1,n do
					t[i] = select(i, ...)
				end
				return f(table.unpack(t, 1, n))
			end
		end,
		-- swaps the next two arguments
		swap = function(f)
			return function(a, b, ...)
				return f(b, a, ...)
			end
		end,
		dump = string.dump,
		-- coroutine access
		wrap = coroutine.wrap,
		co = coroutine.create,
	}
	--]]
}
-- shorthand
functionMeta.__index._ = functionMeta.__index.index
functionMeta.__index.o = functionMeta.__index.compose
functionMeta.__index.o_n = functionMeta.__index.compose_n
debug.setmetatable(function() end, functionMeta)

-- coroutines
debug.setmetatable(coroutine.create(function() end), {__index = coroutine})

-- TODO lightuserdata, if you can create it within lua somehow ...
