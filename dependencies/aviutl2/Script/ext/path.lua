--[[
path(pathtofile):open(mode) - to get a file handle
path(pathtofile):read() - to read a file in entirety
path(pathtofile):write() - to write to a file
path(pathtofile):dir() - to iterate through a directory listing
path(pathtofile):attr() - to get file attributes


TODO

- `path:cwd()` returns the *absolute* cwd, but `path` returns the directory `.` ...
... maybe path:cwd() should return `path'.'` and `path:abs()` should return the absolute path (using lfs.currentdir() for evaluation of '.')

- maybe `path` shoudl be the class, so I can use `path:isa` instead of `path.class:isa` ?

- right now path(a)(b)(c) is the same as path(a)/b/c
... maybe just use /'s and use call for something else? or don't use call at all?
--]]


-- [[ TODO - this block is also in ext/os.lua and ext/file.lua
local detect_os = require 'ext.detect_os'
local detect_lfs = require 'ext.detect_lfs'
local assert = require 'ext.assert'

local io = require 'ext.io'
local os = require 'ext.os'
local string = require 'ext.string'
local class = require 'ext.class'


-- TODO this goes in file or os or somewhere in ext
local function simplifypath(s)
	local p = string.split(s, '/')
	for i=#p-1,1,-1 do
		-- convert a//b's to a/b
		if i > 1 then	-- don't remove empty '' as the first entry - this signifies a root path
			while p[i] == '' do p:remove(i) end
		end
		-- convert a/./b's to a/b
		while p[i] == '.' do p:remove(i) end
		-- convert a/b/../c's to a/c
		if p[i+1] == '..'
		and p[i] ~= '..'
		then
			if i == 1 and p[1] == '' then
				error("/.. absolute + previous doesn't make sense: "..tostring(s))	-- don't allow /../ to remove the base / ... btw this is invalid anyways ...
			end
			p:remove(i)
			p:remove(i)
		end
	end
	-- remove trailing '.''s except the first
	while #p > 1 and p[#p] == '.' do
		p:remove()
	end
	return p:concat'/'
end


-- PREPEND the path if fn is relative, otherwise use fn
-- I should reverse these arguments
-- but this function is really specific to the Path path state variable
local function appendPath(...)
	local fn, p = assert.types('appendPath', 2, 'string', 'string', ...)
	--[[ dont change to path sep ... always use / internally
	if detect_os() then
		fn = os.path(fn)
		if not (
			fn:sub(1,1) == '\\'
			or fn:match'^[A-Z,a-z]:\\'
		) then
			fn = p .. '\\' .. fn
		end
	else
	--]]
		if fn:sub(1,1) ~= '/' then
			fn = p .. '/' .. fn
		end
	--end
	fn = fn:gsub('/%./', '/')
	fn = fn:gsub('/+', '/')
	if #fn > 2 and fn:sub(1,2) == './' then
		fn = fn:sub(3)
	end
	return fn
end


local Path = class()

--Path.sep = os.sep	-- TOO redundant?

function Path:init(args)
	self.path = assert.type(
		assert.type(
			args,
			'table',
			'Path:init args'
		).path,
		'string',
		'Path:init args.path'

	-- always use / internally
	):gsub(
		os.sep,
		'/'
	)
	assert.ne(self.path, nil)
end

-- wrappers
local mappings = {
	[io] = {
		lines = 'lines',
		open = 'open',
		read = 'readfile',
		write = 'writefile',
		append = 'appendfile',
		--getdir = 'getfiledir',	-- defined later, wrapped in Path
		--getext = 'getfileext',	-- defined later, wrapped in Path
	},
	[os] = {
		-- vanilla
		remove = 'remove',
		-- ext
		mkdir = 'mkdir',	-- using os.mkdir instead of lfs.mkdir becasuse of fallbacks ... and 'makeParents' flag
		rmdir = 'rmdir',
		--move = 'move',	-- defined later for picking out path from arg
		exists = 'fileexists',
		isdir = 'isdir',
		--dir = 'listdir',		-- wrapping in path
		--rdir = 'rlistdir',

		-- TODO rename to 'fixpath'? 'fixsep'?
		fixpathsep = 'path',
	},
}
local lfs = detect_lfs()
if lfs then
	mappings[lfs] = {
		attr = 'attributes',
		symattr = 'symlinkattributes',
		cd = 'chdir',
		link = 'link',
		setmode = 'setmode',
		touch = 'touch',
		--cwd = 'currentdir',		-- TODO how about some kind of cwd or something ... default 'file' obj path is '.', so how about relating this to the default. path storage?
		--mkdir = 'mkdir',			-- in 'ext.os'
		--rmdir = 'rmdir',			-- in 'ext.os'
		--lock = 'lock',			-- in 'file' objects via ext.io.open
		--unlock = 'unlock',		-- in 'file' objects via ext.io.open
		lockdir = 'lock_dir',		-- can this be combined with lock() nah since lock() needs an open file handle.
	}
end

for obj,mapping in pairs(mappings) do
	for k,v in pairs(mapping) do
		Path[k] = function(self, ...)
			return obj[v](self.path, ...)
		end
	end
end

-- Path wrapping function, but return wraps in Path
function Path:getdir(...)
	local dir, name = io.getfiledir(self.path, ...)
	return Path{path=dir}, Path{path=name}
end

-- Path wrapping
function Path:getext(...)
	local base, ext = io.getfileext(self.path)
	return Path{path=base}, ext
end

-- [[ same as above but with non-lfs options.
-- TODO put them in io or os like I am doing to abstract non-lfs stuff elsewhere?

function Path:cwd()
	if lfs then
		return Path{path=lfs.currentdir()}
	else
		--[=[ TODO should I even bother with the non-lfs fallback?
		-- if so then use this:
		require 'ffi.req' 'c.stdlib'
		local dirp = unistd.getcwd(nil, 0)
		local dir = ffi.string(dirp)
		ffi.C.free(dirp)
		return dir
		--]=]
		if detect_os() then
			return Path{path=string.trim(io.readproc'cd')}
		else
			return Path{path=string.trim(io.readproc'pwd')}
		end
	end
end
--]]

-- convert relative to absolute paths
function Path:abs()
	if self.path:sub(1,1) == '/' then
		return self
	end
	return Path:cwd()/self
end

-- os.listdir wrapper

function Path:move(to)
	if Path:isa(to) then to = to.path end
	return os.move(self.path, to)
end

function Path:dir()
	if not os.isdir(self.path) then
		error("can't dir() a non-directory: "..tostring(self.path))
	end
	return coroutine.wrap(function()
		for fn in os.listdir(self.path) do
			coroutine.yield(Path{path=fn})
		end
	end)
end

function Path:rdir(callback)
	if not os.isdir(self.path) then
		error("can't rdir() a non-directory: "..tostring(self.path))
	end
	return coroutine.wrap(function()
		for fn in os.rlistdir(self.path, callback) do
			coroutine.yield(Path{path=fn})
		end
	end)
end


-- shorthand for splitting off ext and replacing it
-- Path:getext() splits off the last '.' and returns the letters after it
-- but for files with no ext it returns nil afterwards
-- so a file with only a single '.' at the end will produce a '' for an ext
-- and a file with no '.' will produce ext nil
-- so for compat, handle it likewise
-- for newext == nil, remove the last .ext from the filename
-- for newext == "", replace the last .ext with just a .
function Path:setext(newext)
	local base = self:getext().path
	if newext then
		base = base .. '.' .. newext
	end
	return Path{path=base}
end

-- iirc setting __index and __newindex outside :init() is tough, since so much writing is still going on
--[[
TODO how to do the new interface?
specifically reading vs writing?

instead of __newindex for writing new files, how about path(path):write()
--]]
function Path:__call(k)
	assert.ne(self.path, nil)
	if k == nil then return self end
	if Path:isa(k) then k = k.path end
	local fn = assert.type(
		appendPath(k, self.path),
		'string',
		"Path:__call appendPath(k, self.path)")
	-- is this safe?
	fn = simplifypath(fn)

	-- one option is checking existence and returning nil if no file
	-- but then how about file writing?

	return Path{
		path = assert.type(fn, 'string', "Path:__call simplifypath"),
	}
end

-- clever stl idea: path(a)/path(b) = path(a..'/'..b)
Path.__div = Path.__call

-- return the path but for whatever OS we're using
function Path:__tostring()
	return self:fixpathsep()
end

-- This is intended for shell / cmdline use
-- TODO it doesn't need to quote if there are no special characters present
-- also TODO make sure its escaping matches up with whatever OS is being used
function Path:escape()
	return('%q'):format(self:fixpathsep())
end

Path.__concat = string.concat

-- don't coerce strings just yet
-- don't coerce abs path either
function Path.__eq(a,b) return a.path == b.path end
function Path.__lt(a,b) return a.path < b.path end
function Path.__le(a,b) return a.path <= b.path end

local pathSys = Path{path='.'}

return pathSys
