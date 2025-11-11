-- auto serialize all cmdline params and store them in the global 'cmdline'
-- well, don't store them here, instead return the table from the require()
-- have ext.env store them in the global

local fromlua = require 'ext.fromlua'
local assert = require 'ext.assert'
local table = require 'ext.table'
local string = require 'ext.string'
local tolua = require 'ext.tolua'


local function getCmdline(...)
	-- let cmdline[1...] work as well
	local cmdline = {...}

	for _,w in ipairs{...} do
		local k,v = w:match'^(.-)=(.*)$'
		if k then
			pcall(function()
				cmdline[k] = fromlua(v)
			end)
			if cmdline[k] == nil then cmdline[k] = v end
		else
			cmdline[w] = true
		end
	end

	return cmdline
end



-- this might have officially gotten out of hand...
local function showHelp(cmdValue, cmdKey, cmdline, desc)
	print'specify commands via `command` or `command=value`'
	print()
	print'commands:'
	for _,k in ipairs(table.keys(desc):sort()) do
		local descValue = desc[k]
		if descValue.desc then
			print('\t'..descValue.name..' = '..string.trim(descValue.desc):gsub('\n', '\n\t\t'))
		else
			print('\t'..descValue.name)
		end
		print()
	end
end
-- .validate() signature:
local function showHelpAndQuit(...)
	showHelp(...)
	-- 'and quit' means willingly , vs showing help if we fail, so ...
	-- TODO brief help (if a cmd goes bad, show brief help)
	-- vs full help (i.e. --help etc)
	os.exit(0)
end

--[[
ok now for some stern support
errors if it does't get a validated argument
TODO how to support key'd, indexed, or both cmds ...
for now just validate key'd commands.
TODO equivalent processing for integer indexes, so 'a=1 b=2 c=3' is the same as 'a 1 b 2 c 3'

desc =
	[cmdline key] =
		true = use anything
		string = Lua type of the argument
			or 'number' will cast strings to numbers for you
			or 'file' = string + validate the file exists
		table =
			.type = type of the argument \_ ... use one of these
			.validate = input validation /
			.must = error if the argument isn't present
			.desc = description of the argument (to print for help?)
		function
			= function for validating the argument

--]]
local function validate(desc)
	-- here transform desc into the same format
	for _,name in ipairs(table.keys(desc)) do
		local descValue = desc[name]

		-- use desc[name] to handle the cmdline
		if descValue == true then
			descValue = {}
		elseif type(descValue) == 'string' then
			descValue = {
				type = descValue,
			}
		elseif type(descValue) == 'function' then
			descValue = {
				validate = descValue,
			}
		elseif type(descValue) == 'table' then
			-- fallthru and handle next
		else
			error('idk how to handle this cmdline description '..tolua(descValue))
		end

		if not descValue.type
		and not descValue.validate then
			-- no type/validate is provided? use always
			descValue.validate = function() end
		end

		-- convert desc's with .type into .validate
		if descValue.type then
			assert(not descValue.validate, "you should provide either a .type or a .validate, but not both")
			local descType = descValue.type
			descValue.validate = function(cmdValue, key, cmdline)
				-- special-case casting numbers etc?
				if descType == 'number' then
					cmdline[name] = assert(tonumber(cmdValue))
				elseif descType == 'file' then
					assert.type(cmdValue, 'string')
					assert(require 'ext.path'(cmdValue):exists(), "failed to find file "..tolua(cmdValue))
				else
					assert.type(cmdValue, descType)
				end
			end
			--desc.type = nil	-- still useful?
		end

		descValue.name = name

		desc[name] = descValue
	end

	-- now build our object that we're going to return to the caller

	local cmdlineValidation = {}
	-- cmdline.validate(descs...):fromTable(getCmdline(cmdline...)) is equivalent to cmdline.validate(descs)(cmdline...)
	-- or can be used to handle already-produced cmdline objects
	cmdlineValidation.fromTable = function(self, cmdline)
		-- make sure all cmdline keys are in the description
		-- TODO this is going to assume the int key'd cmdline are from ...
		--  and the string key'd cmdline are from 'k' or 'k=v'
		-- so if someon does '[1]=true', then yes, it will overwrite the int-key'd cmdline, and that should probably be prevented in 'getCmdline'
		for _,k in ipairs(table.keys(cmdline)) do
			local cmdValue = cmdline[k]
			if type(k) == 'number' then
				-- assume its part of the ... sequence
			elseif type(k) == 'string' then
				local descValue = desc[k]
				if not descValue then
					error("got an unknown command "..tolua(k))
				else
					-- assert validation, sometimes overwriting cmdline[k] as we go
					descValue.validate(cmdValue, k, cmdline, desc)
				end
			else
				error("got a cmdline with an unknown key type: "..tolua(k))
			end
		end

		-- make sure all must-be keys are in the command-line
		for k,v in pairs(desc) do
			if v.must then
				if not cmdline[k] then
					error("expected to find key "..k)
				end
			end
		end

		return cmdline
	end
	-- cmdline.validate(descs ...)(cmdline...) operates on the cmdline
	setmetatable(cmdlineValidation, {
		__call = function(self, ...)
			return self:fromTable(getCmdline(...))
		end,
	})
	return cmdlineValidation
end

return setmetatable({
		getCmdline = getCmdline,
		validate = validate,
		showHelp = showHelp,
		showHelpAndQuit = showHelpAndQuit,
	}, {
	__call = function(t,...)
		return getCmdline(...)
	end,
})
