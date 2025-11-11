--[[
the original assert() asserts that the first arg is true, and returns all args, therefore we can assert that the first returned value will also always coerce to true
1) should asserts always return a true value?
2) or should asserts always return the forwarded value?
I'm voting for #2 so assert can be used for wrapping args and not changing behvaior.  when would you need to assert the first arg is true afer the assert has already bene carried out anyways? 
... except for certain specified operations that cannot return their first argument, like assertindex()
--]]

-- cheap 'tolua'
local function tostr(x)
	--[[ just tostring
	return tostring(x)
	--]]
	--[[ also quotes to help distinguish strings-of-numbers from numbers
	if type(x) == 'string' then return ('%q'):format(x) end
	return tostring(x)
	--]]
	--[[ full-on lua serialization ... might have trouble with cdata, especially cdata-primitives
	return require 'ext.tolua'(x)
	--]]
	-- [[ lua type and value
	return type(x)..'('..tostring(x)..')'
	--]]
end

local function prependmsg(msg, str)
	if type(msg) == 'number' then
		msg = tostring(msg)
	end
	if type(msg) == 'nil' then
		return str
	end
	if type(msg) == 'string' then
		return msg..': '..str
	end
	-- not implicitly converted to string -- return as is without message
	return msg
end

local function asserttype(x, t, msg, ...)
	local xt = type(x)
	if xt ~= t then
		error(prependmsg(msg, "expected "..tostring(t).." found "..tostring(xt)))
	end
	return x, t, msg, ...
end

local function assertis(obj, cl, msg, ...)
	if not cl.isa then
		error(prependmsg(msg, "assertis expected 2nd arg to be a class"))
	end
	if not cl:isa(obj) then
		error(prependmsg(msg, "object "..tostring(obj).." is not of class "..tostring(cl)))
	end
	return obj, cl, msg, ...
end

-- how to specify varargs...
-- for now: (msg, N, type1, ..., typeN, arg1, ..., argN)
local function asserttypes(msg, n, ...)
	asserttype(n, 'number', prependmsg(msg, "asserttypes number of args"))
	for i=1,n do
		asserttype(select(n+i, ...), select(i, ...), prependmsg(msg, "asserttypes arg "..i))
	end
	return select(n+1, ...)
end

local function asserteq(a, b, msg, ...)
	if not (a == b) then
		error(prependmsg(msg, "expected "..tostr(a).." == "..tostr(b)))
	end
	return a, b, msg, ...
end

local function asserteqeps(a, b, eps, msg, ...)
	eps = eps or 1e-7
	local normval = math.abs(a - b)
	if normval > eps then
		error((msg and msg..': ' or '').."expected |"..tostr(a).." - "..tostr(b).."| <= "..eps..' but found norm to be '..tostr(normval))
	end
	return a, b, eps, msg, ...
end

local function absdiff(a,b) return math.abs(a - b) end
local function asserteqepsnorm(a, b, eps, norm, msg, ...)
	eps = eps or 1e-7
	norm = norm or absdiff
	local normval = norm(a, b)
	if normval > eps then
		error((msg and msg..': ' or '').."expected |"..tostr(a)..", "..tostr(b).."| <= "..eps..' but found norm to be '..tostr(normval))
	end
	return a, b, eps, norm, msg, ...
end

local function assertne(a, b, msg, ...)
	if not (a ~= b) then
		error(prependmsg(msg, "expected "..tostr(a).." ~= "..tostr(b)))
	end
	return a, b, msg, ...
end

local function assertlt(a, b, msg, ...)
	if not (a < b) then
		error(prependmsg(msg, "expected "..tostr(a).." < "..tostr(b)))
	end
	return a, b, msg, ...
end

local function assertle(a, b, msg, ...)
	if not (a <= b) then
		error(prependmsg(msg, "expected "..tostr(a).." <= "..tostr(b)))
	end
	return a, b, msg, ...
end

local function assertgt(a, b, msg, ...)
	if not (a > b) then
		error(prependmsg(msg, "expected "..tostr(a).." > "..tostr(b)))
	end
	return a, b, msg, ...
end

local function assertge(a, b, msg, ...)
	if not (a >= b) then
		error(prependmsg(msg, "expected "..tostr(a).." >= "..tostr(b)))
	end
	return a, b, msg, ...
end

-- this is a t[k] operation + assert
local function assertindex(t, k, msg, ...)
	if not t then
		error(prependmsg(msg, "object is nil"))
	end
	local v = t[k]
	assert(v, prependmsg(msg, "expected "..tostr(t).."["..tostr(k).." ]"))
	return v, msg, ...
end

-- assert integer indexes 1 to len, and len of tables matches
-- maybe I'll use ipairs... maybe
local function asserttableieq(t1, t2, msg, ...)
	asserteq(#t1, #t2, msg)
	for i=1,#t1 do
		asserteq(t1[i], t2[i], msg)
	end
	return t1, t2, msg, ...
end

-- for when you want to assert a table's length but still want to return the table
-- TODO should this be like assertindex() where it performs the operation and returns the operator value, i.e. returns the length instead of the table?
-- or would that be less usable than asserting the length and returning the table?
local function assertlen(t, n, msg, ...)
	asserteq(#t, n, msg)
	return t, n, msg, ...
end

local function asserterror(f, msg, ...)
	local result, errmsg = pcall(f, ...)
	asserteq(result, false, prependmsg(msg, errmsg))
	-- I'd like to forward all arguments like every assert above
	--return f, msg, ...
	-- but by its nature, "asserterror" means "we expect a discontinuity in execution from this code"
	-- and the calling code wants to see the resulting error information
	-- and since I already error'd if no error was found,
	-- then we already know the pcall's result at this point is false
	-- so I'll do this:
	return errmsg
end

local origassert = _G.assert
return setmetatable({
	type = asserttype,
	types = asserttypes,
	is = assertis,
	eq = asserteq,
	ne = assertne,
	lt = assertlt,
	le = assertle,
	gt = assertgt,
	ge = assertge,
	index = assertindex,
	eqeps = asserteqeps,
	eqepsnorm = asserteqepsnorm,
	tableieq = asserttableieq,
	len = assertlen,
	error = asserterror,
}, {
	-- default `assert = require 'ext.assert'` works, as well as `assertle = assert.le`
	__call = function(t, ...)
		return origassert(...)
	end,
})
