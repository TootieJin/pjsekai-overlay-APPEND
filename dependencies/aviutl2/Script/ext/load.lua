--[[
shim load() and loadstring() and loadfile() layer
comes with callback handler for various AST or string modifications

How to make this more modular so that it could be applied to specific environments and not just _G ?
require 'ext.load'(_G) 	-- modify _G so load, loadstring, loadfile, require all use new load() function with xform shim layer
but if it's returning a function instead of caching result in package.laded after a one-time global modification, then how to store it and associate it per-env, so we don't re-add the layer every time we require it?
How about caching the results per-input-env table?
Of course if anyone else replaces _G, or passes in multiple _ENV's, then we're gonna re-apply this a few times possibly ...

local state = require 'ext.load'(env)
state has the following:
	xforms = table of transforms for the load() function
	oldload = old load function
	load = new load function that is now also assigned to env.load
	oldloadfile = old loadfile function
	loadfile = new loadfile function that is now also assigned to env.loadfile
	olddofile = old dofile function
	dofile = new dofile function that is now also assigned to env.dofile
	oldloadstring = (if env.loadstring exists) old loadstring function
	loadstring = (if env.loadstring exists) new loadstring function that is now also assigned to env.loadstring
	oldsearchfile = old package.searchers[2] or package.loaders[2]
	searchfile = new package.searchers[2] or package.loaders[2]
--]]

local stateForEnv = {}

return function(env)
	env = env or _G

	local state = stateForEnv[env]
	if state then return state end
	state = {}

	require 'ext.xpcall'(env)
	require 'ext.require'(env)

	local package = env.package or _G.package

	-- package.searchpath is missing in 5.1 (but not luajit) ... I need a package.searchpath
	-- Put in its own file?  does anyone else need this?  or not since it pairs with my require-override, which is in this file too.
	local searchpath = package.searchpath
	if not searchpath then
		function searchpath(name, path, sep, rep)
			sep = sep or ';'
			rep = rep or '/'	-- or \ for windows ... TODO meh?
			local namerep = name:gsub('%.', rep)
			local attempted = {}
			for w in path:gmatch('[^'..sep..']*') do	-- TODO escape sep? meh?
				local fn = w
				if fn == '' then	-- search default locations ... which are ... builtin ... or not?
				else
					fn = fn:gsub('%?', namerep)
					-- path(fn):exists() implementation:
					local exists = io.open(fn,'rb')
					if exists then
						exists:close()
						return fn
					end
					table.insert(attempted, "\n\tno file '"..fn.."'")
				end
			end
			return nil, table.concat(attempted)
		end
	end

	state.xforms = setmetatable({}, {__index=table})

	-- TODO proper test?  like if load'string' fails?
	local loadUsesFunctions = (_VERSION == 'Lua 5.1' and not env.jit)
	state.oldload = loadUsesFunctions and (env.loadstring or _G.loadstring) or (env.load or _G.load)

	-- ok here's my modified load behavior
	-- it's going to parse the lua 5.4 code and spit out the luajit code
	state.load = function(data, ...)
		if type(data) == 'function' then
			local s = {}
			repeat
				-- "A return of an empty string, nil, or no value signals the end of the chunk."
				local chunk = data()
				if chunk == '' or chunk == nil then break end
				table.insert(s, chunk)
			until false
			data = table.concat(s)
		end
		-- 5.1 behavior: load(func, name) versus loadstring(data, name)
		-- 5.2..5.4 behavior: load(chunk, name, mode, env)
		-- TODO mind you the formatting on re-converting it will be off ...
		-- errors won't match up ...
		-- so I'll re-insert the generated code
		-- TODO would be nice to save whitespace and re-insert that ... hmm maybe long into the future ...
		-- TODO xpcall behavior testing for when we are allowed to forward the args ... maybe that compat behavior belongs in ext ?
		local source = ... or ('['..data:sub(1,10)..'...]')
		for i,xform in ipairs(state.xforms) do
			local reason
			data, reason = xform(data, source)
			if not data then
				return false, "ext.load.xform["..i.."]: "..(reason and tostring(reason) or '')
			end
		end
		return state.oldload(data, ...)
	end

	-- override global load() function, and maybe loadfile() if it's present too
	-- (maybe loadstring() too ?)
	if env.loadstring ~= nil or _G.loadstring ~= nil then
		state.oldloadstring = env.loadstring or _G.loadstring
		state.loadstring = state.load
		env.loadstring = state.loadstring
	end
	-- TODO if we're in luajit (_VERSION=Lua 5.1) then load() will handle strings, but if we're in lua 5.1 then it will only handle functions (according to docs?) right?
	env.load = state.load

	state.oldloadfile = env.loadfile or _G.loadfile
	-- NOTICE when specifying args (filename, mode, env) explicitly, and forwarding them explicitly, the CLI had some trouble with _ENV var in lua 5.4 ...
	-- so for vanilla lua cli the number of args matters for some reason
	state.loadfile = function(...)
		local filename = ...
		local data, err
		if filename then
			local f
			f, err = io.open(filename, 'rb')
			if not f then return nil, err end
			data, err = f:read'*a'
			f:close()
		else
			data, err = io.read'*a'
		end
		if err then return nil, err end

		-- luajit loadfile/dofile will skip the first # just fine
		-- but lua 5.4 will not ... lua 5.4 seems to only skip the leading # if you run it directly
		-- so here's some luajit compat ...
		if data then
			data = data:match'^#[^\n]*\n(.*)$' or data
		end

		return state.load(data, ...)
	end
	env.loadfile = state.loadfile

	state.olddofile = env.dofile or _G.dofile
	state.dofile = function(filename)
		return assert(state.loadfile(filename))()
	end
	env.dofile = state.dofile

	-- next TODO here , same as ext.debug (consider making modular)
	-- ... wedge in new package.seachers[2]/package.loaders[2] behavior to use my modified load()
	-- replace the package.loaders[2] / package.searchers[2] table entry
	-- make it to replace file contents before loading
	local searchers = assert(package.searchers or package.loaders, "couldn't find searchers")
	state.oldsearchfile = searchers[2]
	state.searchfile = function(req, ...)
		local filename, err = searchpath(req, package.path)
		if not filename then return err end
		local f, err = state.loadfile(filename)
		return f or err
	end
	searchers[2] = state.searchfile

	stateForEnv[env] = state
	return state
end
