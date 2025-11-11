--[[
I use this too often, so now it goes here

override T.getTime to change the timing function
override T.out to change the output file handle
--]]

local hasffi, ffi = pcall(require, 'ffi')

local T = {}

T.out = io.stderr

if not hasffi or ffi.os == 'Windows' then
	-- in linux this is the live time, or something other than the actual time
	T.getTime = os.clock
else
	require 'ffi.req' 'c.time'		-- timegm, gmtime
	require 'ffi.req' 'c.sys.time'	-- gettimeofday
	require 'ffi.req' 'c.string'	-- strerror
	local errno = require 'ffi.req' 'c.errno'
	local gettimeofday_tv = ffi.new'struct timeval[1]'
	function T.getTime()
		local result = ffi.C.gettimeofday(gettimeofday_tv, nil)
		if result ~= 0 then
			error(ffi.string(ffi.C.strerror(errno.errno())))
		end
		return tonumber(gettimeofday_tv[0].tv_sec) + tonumber(gettimeofday_tv[0].tv_usec) / 1000000
	end

	local tm_1 = ffi.typeof'struct tm[1]'

	-- takes in a timestamp (your timezone? do timestamps consider timezone, or are they all UTC?)
	-- spits out UTC date info
	-- TODO add a first 'format' option that formats this... ?
	-- TODO TODO either rename this to 'ext.time' or move it into its own ... repo? file? idk...
	function T.timegm(t)
		local ts = tm_1()
		ts[0].tm_year = (t.year or 1900) - 1900
		ts[0].tm_mon = (t.month or 1) - 1
		ts[0].tm_mday = t.day or 0
		ts[0].tm_hour = t.hour or 12
		ts[0].tm_min = t.min or 0
		ts[0].tm_sec = t.sec or 0
		ts[0].tm_isdst = t.isdst or false
		return ffi.C.timegm(ts)
	end

	function T.time()
		return ffi.C.time(nil)
	end

	local time_t_1 = ffi.typeof'time_t[1]'

	-- takes in UTC date info, spits out a timestamp
	-- pass it unix timestamp, or nil for the current time
	-- returns a date stucture with .year .month .day .hour .min .sec .isdst hopeully with the same range as Lua's os.date
	function T.gmtime(t)
		local tp = time_t_1()
		tp[0] = t or T.time()
		local ts = ffi.C.gmtime(tp)
		return {
			year = ts[0].tm_year + 1900,
			month = ts[0].tm_mon + 1,
			day = ts[0].tm_mday,
			hour = ts[0].tm_hour,
			min = ts[0].tm_min,
			sec = ts[0].tm_sec,
			isdst = ts[0].tm_isdst ~= 0,
		}
	end
end

T.depth = 0
T.tab = ' '

local function timerReturn(name, startTime, indent, ...)
	local endTime = T.getTime()
	local dt = endTime - startTime

	-- this is all printing ...
	T.depth = T.depth - 1
	T.out:write(indent..'...done ')
	if name then
		T.out:write(name..' ')
	end
	T.out:write('('..dt..'s)\n')
	T.out:flush()

	return dt, ...
end

function T.timer(name, cb, ...)
	local indent = T.tab:rep(T.depth)
	if name then
		T.out:write(indent..name..'...\n')
	end
	T.out:flush()
	local startTime = T.getTime()
	T.depth = T.depth + 1
	return timerReturn(name, startTime, indent, cb(...))
end


-- same as above but without printing (and no name too cuz we're not printing)
local function timerReturnQuiet(startTime, ...)
	local endTime = T.getTime()
	local dt = endTime - startTime
	return dt, ...
end

function T.timerQuiet(cb, ...)
	local startTime = T.getTime()
	return timerReturnQuiet(startTime, cb(...))
end


setmetatable(T, {
	-- call forwards to timer:
	__call = function(self, ...)
		return self.timer(...)
	end,
})

return T
