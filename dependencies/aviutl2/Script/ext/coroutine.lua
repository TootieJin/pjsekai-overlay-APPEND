local coroutine = {}
for k,v in pairs(require 'coroutine') do coroutine[k] = v end

local function safehandle(thread, res, ...)
	if not res then
		local err = tostring(...)..'\n'..debug.traceback(thread)
		--[[
		reminder to myself:
		this is here because where else should it go?
		I don't want assertresume to error, but I do want it to print a stacktrace
		it is especially used with threadmanager, which calls assertresume, then gets the status to determine if the thread is running or not
		I could put an assert() around that, but assert(assertresume()) seems strange, and I'd still have to only print the exceptions when the status ~= dead
		So things seem to work best with threadmanager if I put this here.
		Maybe I shouldn't call it 'assertresume' but instead something else like 'resumeAndPrintErrorIfThereIsOne' ?
		--]]
		io.stderr:write(err..'\n')
		io.stderr:flush()
		return false, err
	end
	return true, ...
end

-- resumes thread
-- if the thread is dead, return false
-- if the thread is alive and resume failed due to error, prints the stack trace of the thread upon error
-- 	as opposed to assert(coroutine.resume(thread)), which only prints the stack trace of the resume statement
-- if the thread is alive and resume succeeded, returns true
function coroutine.assertresume(thread, ...)
	if coroutine.status(thread) == 'dead' then return false, 'dead' end
	return safehandle(thread, coroutine.resume(thread, ...))
end

return coroutine
