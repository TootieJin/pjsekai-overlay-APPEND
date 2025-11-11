local io = {}
for k,v in pairs(require 'io') do io[k] = v end

-- io or os?  io since it is shorthand for io.open():read()
function io.readfile(fn)
	local f, err = io.open(fn, 'rb')
	if not f then return false, err end

	-- file.read compat (tested on Windows)
	-- 						*a	a	*l	l
	-- lua-5.3.5:			yes	yes	yes	yes		jit == nil and _VERSION == 'Lua 5.3'
	-- lua-5.2.4:			yes	no	yes	no		jit == nil and _VERSION == 'Lua 5.2'
	-- lua-5.1.5:			yes	no	yes	no		jit == nil and _VERSION == 'Lua 5.1'
	-- luajit-2.1.0-beta3:	yes	yes	yes	yes		(jit.version == 'LuaJIT 2.1.0-beta3' / jit.version_num == 20100)
	-- luajit-2.0.5			yes	no	yes	no		(jit.version == 'LuaJIT 2.0.5' / jit.version_num == 20005)
	local d = f:read('*a')
	f:close()
	return d
end

function io.writefile(fn, d)
	local f, err = io.open(fn, 'wb')
	if not f then return false, err end
	if d then f:write(d) end
	f:close()
	return true
end

function io.appendfile(fn, d)
	local f, err = io.open(fn, 'ab')
	if not f then return false, err end
	if d then f:write(d) end
	f:close()
	return true
end

function io.readproc(cmd)
	local f, err = io.popen(cmd)
	if not f then return false, err end
	local d = f:read('*a')
	f:close()
	return d
end

function io.getfiledir(fn)
	local dir, name = fn:match'^(.*)/([^/]-)$'
	if dir == '' then
		-- "/" => "/", "/"
		if name == '' then return '/', '/' end
		-- "/x" => "/", "x"
		return '/', name
	elseif not dir then
		return '.', fn
	end
	return dir, name
end

-- this should really return the extension first.
-- that is the function name, after all.
function io.getfileext(fn)
	local front, ext = fn:match('^(.*)%.([^%./]-)$')
	if front then
		return front, ext
	end
	-- no ext? then leave that field nil - just return the base filename
	return fn, nil
end

-- in Lua 5.3.5 at least:
-- (for file = getmetatable(io.open(something)))
-- io.read ~= file.read
-- file.__index == file
-- within meta.lua, simply modifying the file metatable
-- but if someone requires ext/io.lua and not lua then io.open and all subsequently created files will need to be modified
--[[ TODO FIXME
if jit or (not jit and _VERSION < 'Lua 5.2') then

	local function fixfilereadargs(...)
		print(...)
		if select('#', ...) == 0 then return ... end
		local fmt = select(1, ...)
		if fmt == 'a' then fmt = '*a'
		elseif fmt == 'l' then fmt = '*l'
		elseif fmt == 'n' then fmt = '*n'
		end
		return fmt, fixfilereadargs(select(2, ...))
	end

	-- even though io.read is basically the same as file.read, they are still different functions
	-- so file.read will still have to be separately overridden
	local oldfileread
	local function newfileread(...)
		return oldfileread(fixfilereadargs(...))
	end
	io.read = function(...)
		return newfileread(io.stdout, ...)
	end

	local oldfilemeta = debug.getmetatable(io.stdout)
	local newfilemeta = {}
	for k,v in pairs(oldfilemeta) do
		newfilemeta[k] = v
	end

	-- override file:read
	oldfileread = oldfilemeta.read
	newfilemeta.read = newfileread

	-- should these be overridden in this case, or only when running ext/meta.lua?
	debug.setmetatable(io.stdin, newfilemeta)
	debug.setmetatable(io.stdout, newfilemeta)
	debug.setmetatable(io.stderr, newfilemeta)

	local function fixfilemeta(...)
		if select('#', ...) > 0 then
			local f = select(1, ...)
			if f then
				debug.setmetatable(f, newfilemeta)
			end
		end
		return ...
	end

	local oldioopen = io.open
	function io.open(...)
		return fixfilemeta(oldioopen(...))
	end
end
--]]

-- [[ add lfs lock/unlock to files
do
	local detect_lfs = require 'ext.detect_lfs'
	local lfs = detect_lfs()
	if lfs then
		-- can I do this? yes on Lua 5.3.  Yes on LuaJIT 2.1.0-beta3
		local filemeta = debug.getmetatable(io.stdout)
		filemeta.lock = lfs.lock
		filemeta.unlock = lfs.unlock
	end
end
--]]

return io
