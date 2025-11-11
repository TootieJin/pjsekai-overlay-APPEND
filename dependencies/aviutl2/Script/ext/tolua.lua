--[[
how to handle recursion ...
a={}
b={}
a.b=b
b.a=a

tolua(a) would give ...
{b={a=error('recursive reference')}}

but how about, if something is found that is marked in touched tables ...
1) wrap everything in a function block
2) give root a local
3) add assignments of self-references after-the-fact

(function()
	local _tmp={b={}}
	_tmp.b.a= _tmp
	return _tmp
end)()
--]]

local table = require 'ext.table'

local function builtinPairs(t)
	return next,t,nil
end

local _0byte = ('0'):byte()
local _9byte = ('9'):byte()
local function escapeString(s)
	--[[ multiline strings
	-- it seems certain chars can't be encoded in Lua multiline strings
	-- TODO find out exactly which ones
	-- TODO if 's' begins with [ or ends with ] then you're gonna have a bad time ...
	-- in fact ... does Lua support multiline strings that begin with [ or end with ] ?  especially the latter ...
	local foundNewline
	local foundBadChar
	for i=1,#s do
		local b = s:byte(i)
		if b == 10
		-- TODO still lua loading will convert \r \n's into \n's ... so this is still not guaranteed to reproduce the original ...
		-- I could disable multi line strs when encoding \r's ...
		-- but this gets to be a bit os-specific ...
		-- a better solution would be changing fromlua() to properly handle newline formatting within multiline strings ...
		--or b == 13
		then
			foundNewline = true
		elseif b < 32 or b > 126 then
			foundBadChar = true
			break	-- don't need to keep looking
		end
	end
	if foundNewline and not foundBadChar then
		for neq=0,math.huge do
			local eq = ('='):rep(neq)
			local open = '['..eq..'['
			local close = ']'..eq..']'
			if not s:find(open, 1, true)
			and not s:find(close, 1, true)
			and s:sub(-neq-1) ~= close:sub(1,neq+1) -- if the very end of the string is close without the last ] then it could still do a false match ...
			then
				-- open and close aren't in the string, we can use this to escape the string
				-- ... ig all I have to search for is close, but meh
				local ret = open .. '\n' 	-- \n cuz lua ignores trailing spaces/newline after the opening
					.. s .. close
--DEBUG: require 'ext.assert'.eq(load('return '..ret)(), s)
				return ret
			end
		end
	end
	--]]

	-- [[
	-- this will only escape escape codes
	-- will respect unicode
	-- but it skips \r \t and encodes them as \13 \9
	local o = ('%q'):format(s)
	o = o:gsub('\\\n','\\n')
	return o
	--]]
	--[==[ this gets those that builtin misses
	-- but does it in lua so it'll be slow
	-- and requires implementations of iscntrl and isdigit
	--
	-- it's slow and has bugs.
	--
	-- TODO
	-- for size-minimal strings:
	-- if min(# single-quotes, # double-quotes) within the string > 2 then use [[ ]] (so long as that isn't used either)
	-- otherwise use as quotes whatever the min is
	-- or ... use " to wrap if less than 1 " is embedded
	-- then use ' to wrap if less than 1 ' is embedded
	-- then use [[ ]] to wrap if no [[ ]] is embedded
	-- ... etc for [=...=[ all string escape options
	local o = '"'
	for i=1,#s do
		local c = s:sub(i,i)
		if c == '"' then
			o = o .. '\\"'
		elseif c == '\\' then
			o = o .. '\\\\'
		elseif c == '\n' then
			o = o .. '\\n'
		elseif c == '\r' then
			o = o .. '\\r'
		elseif c == '\t' then
			o = o .. '\\t'
		elseif c == '\a' then
			o = o .. '\\a'
		elseif c == '\b' then
			o = o .. '\\b'
		elseif c == '\f' then
			o = o .. '\\f'
		elseif c == '\v' then
			o = o .. '\\v'
		else
			local b = c:byte()
			assert(b < 256)
			if b < 0x20 or b == 0x7f then	-- if iscntrl(c)
-- make sure the next character isn't a digit because that will mess up the encoded escape code
				local b2 = c:byte(i+1)
				if not (b2 and b2 >= _0byte and b2 <= _9byte) then	-- if not isdigit(c2) then
					o = o .. ('\\%d'):format(b)
				else
					o = o .. ('\\%03d'):format(b)
				end
			else
				-- TODO for extended ascii, why am I seeing different things here vs encoding one character at a time?
				o = o .. c
			end
		end
	end
	o = o .. '"'
	o:gsub('\\(%d%d%d)', function(d)
		if tonumber(d) > 255 then
			print('#s', #s)
			print'o'
			print(o)
			print's'
			print(s)
			error("got an oob escape code: "..d)
		end
	end)
	local f = require 'ext.fromlua'(o)
	if f ~= s then
		print('#s', #s)
		print('#f', #f)
		print'o'
		print(o)
		print's'
		print(s)
		print'f'
		print(f)
		print("failed to reencode as the same string")
		for i=1,math.min(#s,#f) do
			if f:sub(i,i) ~= s:sub(i,i) then
				print('char '..i..' differs')
				break
			end
		end
		error("here")

	end
	return o
	--]==]
end

-- as of 5.4.  I could modify this based on the Lua version (like removing 'goto') but misfiring just means wrapping in quotes, so meh.
local reserved = {
	["and"] = true,
	["break"] = true,
	["do"] = true,
	["else"] = true,
	["elseif"] = true,
	["end"] = true,
	["false"] = true,
	["for"] = true,
	["function"] = true,
	["goto"] = true,
	["if"] = true,
	["in"] = true,
	["local"] = true,
	["nil"] = true,
	["not"] = true,
	["or"] = true,
	["repeat"] = true,
	["return"] = true,
	["then"] = true,
	["true"] = true,
	["until"] = true,
	["while"] = true,
}

-- returns 'true' if k is a valid variable name, but not a reserved keyword
local function isVarName(k)
	return type(k) == 'string' and k:match('^[_a-zA-Z][_a-zA-Z0-9]*$') and not reserved[k]
end

local toLuaRecurse

local function toLuaKey(state, k, path)
	if isVarName(k) then
		return k, true
	else
		local result = toLuaRecurse(state, k, nil, path, true)
		if result then
			return '['..result..']', false
		else
			return false, false
		end
	end
end


-- another copy of maxn, with custom pairs
local function maxn(t, state)
	local max = 0
	local count = 0
	for k,v in state.pairs(t) do
		count = count + 1
		if type(k) == 'number' then
			max = math.max(max, k)
		end
	end
	return max, count
end


local defaultSerializeForType = {
	number = function(state,x)
		if x == math.huge then return 'math.huge' end
		if x == -math.huge then return '-math.huge' end
		if x ~= x then return '0/0' end
		return tostring(x)
	end,
	boolean = function(state,x) return tostring(x) end,
	['nil'] = function(state,x) return tostring(x) end,
	string = function(state,x) return escapeString(x) end,
	['function'] = function(state, x)
		local result, s = pcall(string.dump, x)

		if result then
			s = 'load('..escapeString(s)..')'
		else
			-- if string.dump failed then check the builtins
			-- check the global object and one table deep
			-- todo maybe, check against a predefined set of functions?
			if s == "unable to dump given function" then
				local found
				for k,v in state.pairs(_G) do
					if v == x then
						found = true
						s = k
						break
					elseif type(v) == 'table' then
						-- only one level deep ...
						for k2,v2 in state.pairs(v) do
							if v2 == x then
								s = k..'.'..k2
								found = true
								break
							end
						end
						if found then break end
					end
				end
				if not found then
					s = "error('"..s.."')"
				end
			else
				return "error('got a function I could neither dump nor lookup in the global namespace nor one level deep')"
			end
		end

		return s
	end,
	table = function(state, x, tab, path, keyRef)
		local result

		local newtab = tab .. state.indentChar
		-- TODO override for specific metatables?  as I'm doing for types?

		if state.touchedTables[x] then
			if state.skipRecursiveReferences then
				result = 'error("recursive reference")'
			else
				result = false	-- false is used internally and means recursive reference
				state.wrapWithFunction = true

				-- we're serializing *something*
				-- is it a value?  if so, use the 'path' to dereference the key
				-- is it a key?  if so the what's the value ..
				-- do we have to add an entry for both?
				-- maybe the caller should be responsible for populating this table ...
				if keyRef then
					state.recursiveReferences:insert('root'..path..'['..state.touchedTables[x]..'] = error("can\'t handle recursive references in keys")')
				else
					state.recursiveReferences:insert('root'..path..' = '..state.touchedTables[x])
				end
			end
		else
			state.touchedTables[x] = 'root'..path

			-- prelim see if we can write it as an indexed table
			local numx, count = maxn(x, state)
			local intNilKeys, intNonNilKeys
			-- only count if our max isn't too high
			if numx < 2 * count then
				intNilKeys, intNonNilKeys = 0, 0
				for i=1,numx do
					if x[i] == nil then
						intNilKeys = intNilKeys + 1
					else
						intNonNilKeys = intNonNilKeys + 1
					end
				end
			end

			local hasSubTable

			local s = table()

			-- add integer keys without keys explicitly. nil-padded so long as there are 2x values than nils
			local addedIntKeys = {}
			if intNonNilKeys
			and intNilKeys
			and intNonNilKeys >= intNilKeys * 2
			then	-- some metric for when to create in-order tables
				for k=1,numx do
					if type(x[k]) == 'table' then hasSubTable = true end
					local nextResult = toLuaRecurse(state, x[k], newtab, path and path..'['..k..']')
					if nextResult then
						s:insert(nextResult)
					-- else x[k] is a recursive reference
					end
					addedIntKeys[k] = true
				end
			end

			-- sort key/value pairs added here by key
			local mixed = table()
			for k,v in state.pairs(x) do
				if not addedIntKeys[k] then
					if type(v) == 'table' then hasSubTable = true end
					local keyStr, usesDot = toLuaKey(state, k, path)
					if keyStr then
						local newpath
						if path then
							newpath = path
							if usesDot then newpath = newpath .. '.' end
							newpath = newpath .. keyStr
						end
						local nextResult = toLuaRecurse(state, v, newtab, newpath)
						if nextResult then
							mixed:insert{keyStr, nextResult}
						-- else x[k] is a recursive reference
						end
					end
				end
			end
			mixed:sort(function(a,b) return a[1] < b[1] end)	-- sort by keys
			mixed = mixed:map(function(kv) return table.concat(kv, '=') end)
			s:append(mixed)

			local thisNewLineChar, thisNewLineSepChar, thisTab, thisNewTab
			if not hasSubTable and not state.alwaysIndent then
				thisNewLineChar = ''
				thisNewLineSepChar = ' '
				thisTab = ''
				thisNewTab = ''
			else
				thisNewLineChar = state.newlineChar
				thisNewLineSepChar = state.newlineChar
				thisTab = tab
				thisNewTab = newtab
			end

			local rs = '{'..thisNewLineChar
			if #s > 0 then
				rs = rs .. thisNewTab .. s:concat(','..thisNewLineSepChar..thisNewTab) .. thisNewLineChar
			end
			rs = rs .. thisTab .. '}'

			result = rs
		end
		return result
	end,
}

local function defaultSerializeMetatableFunc(state, m, x, tab, path, keyRef)
	-- only serialize the metatables of tables
	-- otherwise, assume the current metatable is the default one (which is usually nil)
	if type(x) ~= 'table' then return 'nil' end
	return toLuaRecurse(state, m, tab..state.indentChar, path, keyRef)
end

toLuaRecurse = function(state, x, tab, path, keyRef)
	if not tab then tab = '' end

	local xtype = type(x)
	local serializeFunction
	if state.serializeForType then
		serializeFunction = state.serializeForType[xtype]
	end
	if not serializeFunction then
		serializeFunction = defaultSerializeForType[xtype]
	end

	local result
	if serializeFunction then
		result = serializeFunction(state, x, tab, path, keyRef)
	else
		result = '['..type(x)..':'..tostring(x)..']'
	end
	assert(result ~= nil)

	if state.serializeMetatables then
		local m = getmetatable(x)
		if m ~= nil then
			local serializeMetatableFunc = state.serializeMetatableFunc or defaultSerializeMetatableFunc
			local mstr = serializeMetatableFunc(state, m, x, tab, path, keyRef)
			-- make sure you get something
			assert(mstr ~= nil)
			-- but if that something is nil, i.e. setmetatable(something newly created with a nil metatable, nil), then don't bother modifing the code
			if mstr ~= 'nil' and mstr ~= false then
				-- if this is false then the result was deferred and we need to add this line to wherever else...
				assert(result ~= false)
				result = 'setmetatable('..result..', '..mstr..')'
			end
		end
	end

	return result
end

--[[
args:
	indent = default to 'true', set to 'false' to make results concise, true will skip inner-most tables. set to 'always' for always indenting.
	pairs = default to a form of pairs() which iterates over all fields using next().  Set this to your own custom pairs function, or 'pairs' if you would like serialization to respect the __pairs metatable (which it does not by default).
	serializeForType = a table with keys of lua types and values of callbacks for serializing those types
	serializeMetatables = set to 'true' to include serialization of metatables
	serializeMetatableFunc = function to override the default serialization of metatables
	skipRecursiveReferences = default to 'false', set this to 'true' to not include serialization of recursive references
--]]
local function tolua(x, args)
	local state = {
		indentChar = '',
		newlineChar = '',
		wrapWithFunction = false,
		recursiveReferences = table(),
		touchedTables = {},
	}
	local indent = true
	if args then
		-- indent == ... false => none, true => some, "always" => always
		if args.indent == false then indent = false end
		if args.indent == 'always' then state.alwaysIndent = true end
		state.serializeForType = args.serializeForType
		state.serializeMetatables = args.serializeMetatables
		state.serializeMetatableFunc = args.serializeMetatableFunc
		state.skipRecursiveReferences = args.skipRecursiveReferences
	end
	if indent then
		state.indentChar = '\t'
		state.newlineChar = '\n'
	end
	state.pairs = builtinPairs

	local str = toLuaRecurse(state, x, nil, '')

	if state.wrapWithFunction then
		str = '(function()' .. state.newlineChar
			.. state.indentChar .. 'local root = '..str .. ' ' .. state.newlineChar
			-- TODO defer self-references to here
			.. state.recursiveReferences:concat(' '..state.newlineChar..state.indentChar) .. ' ' .. state.newlineChar
			.. state.indentChar .. 'return root ' .. state.newlineChar
			.. 'end)()'
	end

	return str
end

return setmetatable({}, {
	__call = function(self, x, args)
		return tolua(x, args)
	end,
	__index = {
		-- escaping a Lua string for load() to use
		escapeString = escapeString,
		-- returns 'true' if the key passed is a valid Lua variable string, 'false' otherwise
		isVarName = isVarName,
		-- table of default serialization functions indexed by each time
		defaultSerializeForType = defaultSerializeForType,
		-- default metatable serialization function
		defaultSerializeMetatableFunc = defaultSerializeMetatableFunc,
	}
})
