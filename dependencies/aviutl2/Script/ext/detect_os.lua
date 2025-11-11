local detect_ffi = require 'ext.detect_ffi'
local result
local function detect_os()
	if result ~= nil then return result end
	local ffi = detect_ffi()
	if ffi then
		result = ffi.os == 'Windows'
	else
		-- TODO what if uname doesn't exist? then this will output to stderr.  does it exist in Windows?
		-- to get around that on Windows I can pipe to > NUL
		-- TODO what if it's not Windows?  then this will create a NUL file ...
		-- honestly I could just use the existence of piping to NUL vs /dev/null to determine Windows vs Unix ...
		result = ({
			msys = true,
			ming = true,
		--})[(io.popen'uname 2> NUL':read'*a'):sub(1,4):lower()]
		})[(io.popen'uname':read'*a'):sub(1,4):lower()] or false
	end

	-- right now just true/value for windows/not
	return result
end
return detect_os
