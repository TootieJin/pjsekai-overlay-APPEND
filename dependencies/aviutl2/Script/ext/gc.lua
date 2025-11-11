--[[
lua 5.2+ supports __gc for tables and userdata
lua 5.1 supports __gc for userdata only
luajit (5.1+...?) supports __gc for userdata and cdata ... but not tables

but luckily 5.1 and luajit have this 'newproxy' function

To use this, just call "require 'ext.gc'" once to override the 'setmetatable' function to a new one that uses 'newproxy' to provide table metatables
If this doesn't find 'newproxy' (not in Lua 5.2+) then it doesn't override anything.
--]]

if not newproxy then return end

-- from https://stackoverflow.com/a/77702023/2714073
-- TODO accept an 'env' param like so many other of my override functions ... maybe ... tho I don't do this with all, do I?
local gcProxies = setmetatable({}, {__mode='k'})
local oldsetmetatable = setmetatable
function setmetatable(t, mt)
	local oldp = gcProxies[t]
	if oldp then
		getmetatable(oldp).__gc = nil
		--oldsetmetatable(oldp, nil)
	end

	if mt and mt.__gc then
		local p = newproxy(true)
		gcProxies[t] = p
		getmetatable(p).__gc = function()
			if type(mt.__gc) == 'function' then
				mt.__gc(t)
			end
		end
	end

	return oldsetmetatable(t, mt)
end
