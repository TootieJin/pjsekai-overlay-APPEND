local ffi = require 'ffi'
require 'ffi.req' 'c.stdlib'

--[[
uses C malloc paired with ffi-based garbage collection
typecasts correctly
and retains the ability to manually free
(so you don't have to manually free it)
NOTICE casting *AFTER* wrapping will crash, probably due to the gc thinking the old pointer is gone
also ffi.gc retains type, so no worries about casting before
...
that was true, but now it's always losing the ptr and crashing, so I'm going to fall back on ffi.new
--]]
local function gcnew(T, n)
	-- [[
	local ptr = ffi.C.malloc(n * ffi.sizeof(T))
	ptr = ffi.cast(T..'*', ptr)
	ptr = ffi.gc(ptr, ffi.C.free)
	--]]
	--[[
	local ptr = ffi.new(T..'['..n..']')
	--]]
	return ptr
end

--[[
manual free of a pointer
frees the ptr and removes it from the gc
(just in case you want to manually free a pointer)
--]]
local function gcfree(ptr)
	ffi.C.free(ffi.gc(ptr, nil))
end

return {
	new = gcnew,
	free = gcfree,
}
