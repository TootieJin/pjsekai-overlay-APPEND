--[[
notice that,
while this does override the 'string' and add some extra stuff,
it does not explicitly replace the default string metatable __index
to do that, require 'ext.meta' (or do it yourself)
--]]
local string = {}
for k,v in pairs(require 'string') do string[k] = v end

local table = require 'ext.table'

-- table.concat(string.split(a,b),b) == a
function string.split(s, exp)
	exp = exp or ''
	s = tostring(s)
	local t = table()
	-- handle the exp='' case
	if exp == '' then
		for i=1,#s do
			t:insert(s:sub(i,i))
		end
	else
		local searchpos = 1
		local start, fin = s:find(exp, searchpos)
		while start do
			t:insert(s:sub(searchpos, start-1))
			searchpos = fin+1
			start, fin = s:find(exp, searchpos)
		end
		t:insert(s:sub(searchpos))
	end
	return t
end

function string.trim(s)
	return s:match('^%s*(.-)%s*$')
end

-- should this wrap in a table?
function string.bytes(s)
	return table{s:byte(1,#s)}
end

string.load = load or loadstring

--[[
-- drifting further from standards...
-- this string-converts everything concat'd (no more errors, no more print(a,b,c)'s)
getmetatable('').__concat = function(a,b)
	return tostring(a)..tostring(b)
end
--]]

-- a C++-ized accessor to subsets
-- indexes are zero-based inclusive
-- sizes are zero-based-exclusive (or one-based-inclusive depending on how you think about it)
-- parameters are (index, size) rather than (start index, end index)
function string.csub(d, start, size)
	if not size then return string.sub(d, start + 1) end	-- til-the-end
	return string.sub(d, start + 1, start + size)
end

--d = string data
--l = length of a column.  default 32
--w = hex word size.  default 1
--c = extra column space.  default 8
function string.hexdump(d, l, w, c)
	d = tostring(d)
	l = tonumber(l)
	w = tonumber(w)
	c = tonumber(c)
	if not l or l < 1 then l = 32 end
	if not w or w < 1 then w = 1 end
	if not c or c < 1 then c = 8 end
	local s = table()
	local rhs = table()
	local col = 0
	for i=1,#d,w do
		if i % l == 1 then
			s:insert(string.format('%.8x ', (i-1)))
			rhs = table()
			col = 1
		end
		s:insert' '
		for j=w,1,-1 do
			local e = i+j-1
			local sub = d:sub(e,e)
			if #sub > 0 then
				local b = string.byte(sub)
				s:insert(string.format('%.2x', b))
				rhs:insert(b >= 32 and sub or '.')
			end
		end
		if col % c == 0 then
			s:insert' '
		end
		if (i + w - 1) % l == 0 or i+w>#d then
			s:insert' '
			s:insert(rhs:concat())
		end
		if (i + w - 1) % l == 0 then
			s:insert'\n'
		end
		col = col + 1
	end
	return s:concat()
end

-- escape for pattern matching
local escapeFind = '[' .. ([[^$()%.[]*+-?]]):gsub('.', '%%%1') .. ']'
function string.patescape(s)
	return (s:gsub(escapeFind, '%%%1'))
end

-- this is a common function, especially as a __concat metamethod
-- it is nearly table.concat, except table.concat errors upon non-string/number instead of calling tostring() automatically
-- (should I change table.concat's default behavior and use that instead?  nah, because why require a table creation.)
-- tempted to make this ext.op.concat ... but that's specifically a binary op ... and that shouldn't call tostring() while this should ...
-- maybe I should move this to ext.op as 'tostringconcat' or something?
function string.concat(...)
	local n = select('#', ...)
	if n == 0 then return end	-- base-case nil or "" ?
	local s = tostring((...))
	if n == 1 then return s end
	return s .. string.concat(select(2, ...))
end

-- another common __tostring metamethod
-- since luajit doesn't support __name metafield
function string:nametostring()
	-- NOTICE this will break for anything that overrides its __metatable metafield
	local mt = getmetatable(self)

	-- invoke a 'rawtostring' call / get the builtin 'tostring' result
	setmetatable(self, nil)
	local s = tostring(self)
	setmetatable(self, mt)

	local name = mt.__name
	return name and tostring(name)..s:sub(6) or s
end

-- I use this too often ....
function string.hex(s, uppercase)
	local fmt = uppercase and '%02X' or '%02x'
	return (tostring(s):gsub('.', function(c)
		return fmt:format(c:byte())
	end))
end

function string.unhex(h)
	if bit.band(#h, 1) == 1
	or h:find'[^0-9a-fA-F]'
	then
		return nil, "string is not hex"
	end
	return h:gsub('..', function(d)
		return string.char(assert(tonumber(d, 16)))
	end)
end

-- TODO other encoders?

return string
