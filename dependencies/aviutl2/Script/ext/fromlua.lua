return function(str, ...)
	return assert(load('return '..str, ...))()
end
