local lfs
local function detect_lfs()
	if lfs == nil then
		-- ok in a naive world luajit => lfs_ffi and lua => lfs
		-- esp in a windows world where giving your lua and luajit scripts the same package.path/cpath's will cause crashes
		-- but on openresty we're using luajit but using lfs so ...
		-- don't mix up your package.path/cpath's
		-- and i'm gonna try both
		for _,try in ipairs{'lfs', 'lfs_ffi'} do
			local result
			result, lfs = pcall(require, try)
			lfs = result and lfs
			if lfs then break end
		end
	end
	return lfs
end
return detect_lfs
