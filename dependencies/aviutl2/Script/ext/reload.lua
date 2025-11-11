-- this goes wherever packages and the likes would go
-- it would be awesome to track loaded dependencies and auto-reload them too
-- should be easy to do by overriding 'require'
local function reload(n)
	package.loaded[n] = nil
	return require(n)
end

return reload
