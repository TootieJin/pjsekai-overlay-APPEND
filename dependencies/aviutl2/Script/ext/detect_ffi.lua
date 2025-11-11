-- this runs just once and then gets cached in package.path
-- but return the function, not the value, because the value could be false, which (I don't believe) package.path will cache
local ffi
local function detect_ffi()
	if ffi == nil then
		local result
		result, ffi = pcall(require, 'ffi')
		ffi = result and ffi
	end
	return ffi
end
return detect_ffi
