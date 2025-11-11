-- pure-lua require since the luajit one has C-call boundary yield problems (maybe?  10yo bug?  https://github.com/openresty/lua-nginx-module/issues/376#issuecomment-46104284 )
return function(env)
	env = env or _G
	--local oldrequire = env.require
	local package = env.package
	env.require = function(modname)
		local mod = package.loaded[modname]
		if mod ~= nil then return mod end
		local errs = {"module '"..tostring(modname).."' not found:"}
		local searchers = assert(package.searchers or package.loaders, "couldn't find searchers")
		for _,search in ipairs(searchers) do
			local result = search(modname)
			local resulttype = type(result)
			if resulttype == 'function' then
				local mod = result() or true
				package.loaded[modname] = mod
				return mod
			elseif resulttype == 'string' or resulttype == 'nil' then
				table.insert(errs, result)	-- if I get a nil, should I convert it into an error?
			else
				error("package.searchers result got an unknown type: "..resulttype)
			end
		end
		error(table.concat(errs))
	end
end
