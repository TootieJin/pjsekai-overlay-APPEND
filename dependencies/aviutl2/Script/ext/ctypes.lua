--[[
put all the ctypes in scope for more c++-ish programming

I'm not sure if this belongs in ext or ffi or anywhere else ...
I probably won't add it to git
Maybe I'll never use this, since it primarily operates on the global namespace

hmm could/should I bother cache these in a module table,
and then in ext.env write to _G?
would it save anything to lookup in require'ext.ctypes'.double versus ffi.typeof'double' ?
--]]

local ffi = require 'ffi'

for _,t in ipairs{
	'char',
	--'byte', = uint8_t right?
	'short',
	'int',
	'long',
	'int8_t',
	'uint8_t',
	'int16_t',
	'uint16_t',
	'int32_t',
	'uint32_t',
	'int64_t',
	'uint64_t',
	'float',
	'double',
	--'float32' = float right
	--'float64' = double right?
	--'float80' = long double ... in most cases ... ? or is long double sometimes float128?
} do
	_G[t] = ffi.typeof(t)
end

_G.byte = ffi.typeof'uint8_t'
_G.float32 = ffi.typeof'float'
_G.float64 = ffi.typeof'double'
_G.float80 = ffi.typeof'long double'

_G.ffi = ffi
_G.sizeof = ffi.sizeof
_G.typeof = ffi.typeof
