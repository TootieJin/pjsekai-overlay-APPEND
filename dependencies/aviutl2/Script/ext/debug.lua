--[[
Here's a quick hack for debugging
It's not meta-lua, not dependent on my parser lib, nothing like that
It does use my new load() shim layer ext.load
Just this: when you require() a file, if the debug tag is set, then grep all --DEBUG: lines to remove the comment
Usage: `lua -lext.debug ...`
or if you want to query specific dags: `lua -e "require 'ext.debug' 'tag1,tag2,tag3'" ...`


But now I'm tired of always tagging files, so how about a proper query script.
Something that looks like this:
	lua -e "require 'ext.debug' 'source:match'sdl*' and level>3 and func:match'SDLAPP' and tag:match'testing''"

That's great but those are a lot of variables that I haven't got yet.

What kind of variables would we like:
	source = string of the current loaded source
	level = log-level, needs to be specified in the DEBUG tag
	func = current function ... how to pick out, aside from making this use the parser (that means moving it out of ext and into its own library...)
	tag = same as before

How to specify things like tag and level ...

DEBUG(tag@level): ?
DEBUG(@1):

What should default log level be?  1.
--]]

function string_split(s, exp)
	exp = exp or ''
	s = tostring(s)
	local t = {}
	-- handle the exp='' case
	if exp == '' then
		for i=1,#s do
			table.insert(t, s:sub(i,i))
		end
	else
		local searchpos = 1
		local start, fin = s:find(exp, searchpos)
		while start do
			table.insert(t, s:sub(searchpos, start-1))
			searchpos = fin+1
			start, fin = s:find(exp, searchpos)
		end
		table.insert(t, s:sub(searchpos))
	end
	return t
end

local oldload = load or loadstring
local logcond

--[[
Strip out DEBUG: and DEBUG(...): tags based on what tags were requested via `require 'ext.debug'(tag1,tag2,...)`

d = data
source = for require'd files will be an absolute file path, from ext.load, from package.searchers[2]

Tags will look like:
`--DEBUG(tag):` to specify a tag.
`--DEBUG(@level):` to specify a level.
`--DEBUG(tag@level):` to specify both.
--]]
table.insert(require 'ext.load'().xforms, function(d, source)
	local result = {}
	local ls = string_split(d, '\n')
	for lineno=1,#ls do
		local l = ls[lineno]
		local found
		repeat
			found = false

			local start, fin = l:find'%-%-DEBUG:'
			if start then
				if logcond(source, lineno, 1) then
					l = l:sub(1, start-1)..l:sub(fin+1)
					ls[lineno] = l
					found = true
				end
			end

			local start, fin = l:find'%-%-DEBUG%b():'
			if start then
				local inside = l:sub(start+8, fin-2)
				local tag, level = table.unpack(string_split(inside, '@'))
				level = tonumber(level)
				if logcond(source, lineno, level, tag) then
					l = l:sub(1, start-1)..l:sub(fin+1)
					ls[lineno] = l
					found = true
				end
			end
		until not found
	end
	d = table.concat(ls, '\n')
	return d
end)

--[[
Use this with: `lua -e "require 'ext.debug' 'source:match'sdl*' or level>3 and tag:match'testing''"`

Conditions can use the variables: source, line, level, tag
--]]
local function setCond(condstr)
	local code = 'local source, line, level, tag = ... return '..condstr
	logcond = assert(oldload(code))
end

setCond'level == 1'

return setCond
