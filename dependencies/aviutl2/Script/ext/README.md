## Lua Extension Library

[![Donate via Stripe](https://img.shields.io/badge/Donate-Stripe-green.svg)](https://buy.stripe.com/00gbJZ0OdcNs9zi288)<br>

I found myself recreating so many extensions to the base lua classes over and over again.
I thought I'd just put them in one place.
I'm sure this will grow out of hand.

Note to users: The structure of the source code doesn't exactly match the structure in the rock install destination.
This is because I personally use a `LUA_PATH` pattern of "?/?.lua" in addition to the typical "?.lua".
To work with compatability of everyone else who does not use this convention, I have the rockspec install `ext/ext.lua` into `ext.lua` and keep `ext/everything_else.lua` at `ext/everything_else.lua`.

TLDR, how to set up your environment, choose one:
- Install the rock and don't worry.
- Move ext.lua into the parent directory.
- Add "?/?.lua" to your `LUA_PATH`.


Descriptions of the Files:

### assert.lua

Returns an object that can be used as a typical `assert()` call, but also contains several functions that can be used for testing and outputting specific assert cases:

`assert.eq(a, b, [msg, ...])` = If `a == b` returns all arguments.  Otherwise error with message specifying `a`, `b`, and optionally `msg`.

`assert.eqeps(a, b, [eps, msg, ...])` = If `|a-b|<=eps` returns all arguments.  Otherwise error with message specifying `a`, `b`, `eps`, and optionally `msg`.  Default `eps` is `1e-7`.

`assert.eqepsnorm(a, b, [eps, norm, msg, ...])` = If `norm(a,b)<=eps` returns all arguments.  Otherwise error with message specifying `a`, `b`, `eps`, and optionally `msg`.  Default `eps` is `1e-7`, and default `norm(a,b)` is `math.abs(a - b)`.

`assert.ne(a, b, [msg, ...])` = If `a ~= b` returns all arguments.  Otherwise error.

`assert.ge(a, b, [msg, ...])` = If `a >= b` returns all arguments.  Otherwise error.

`assert.gt(a, b, [msg, ...])` = If `a > b` returns all arguments.  Otherwise error.

`assert.le(a, b, [msg, ...])` = If `a <= b` returns all arguments.  Otherwise error.

`assert.lt(a, b, [msg, ...])` = If `a < b` returns all arguments.  Otherwise error.

`assert.len(t, n, [msg, ...])` = If `#t == n` returns all arguments.  Otherwise error.

`assert.type(x, t, [msg, ...])` = If `type(x) == t` returns all arguments.  Otherwise error.

`assert.types(msg, n, ...)` = Assumes `...` holds `n` types and then `n` variables.  If all the `[n+1,2*n]` variables' types match the `[1,n]` types then it returns all the variables.  Otherwise error.

`assert.index(t, k, [msg, ...])` = If `t[k]` evaluates to true then returns `t[k], msg, ...`.  Otherwise error.

`assert.tableieq(t1, t2, [msg, ...])` = If `#t1 == #t2` and `t1[i] == t2[i]` then returns all arguments.  Otherwise error.

`assert.error(f, [msg, ...])` = If calling `f(...)` does not produce an error then this errors.  If calling `f(...)` does error then this returns the error message.

`assert.is(o, cl, [msg, ...])` = If object `o` is an instance of class `cl` then returns all arguments.  Otherwise error.

### ext.lua

`require`ing this file sets up global environment via ext.env and overrides metatables to all default Lua types via ext.meta.
I would not recommend `require`ing this directly in any large scale projects, as it changes many default operations of Lua.
It is much more useful for small-scale scripts.

### env.lua

This file adds all tables to the specified Lua environment table (`_G` by default).
It also sets _ to os.execute for shorthand shell scripting.

### class.lua:

`class(parent1, parent2, ...)` = create a 'class' table, fill it with the union of the table 'parent1', 'parent2', etc.

Class tables can instanciate object tables using the `\_\_call` operator:

``` lua
Cl = class()
obj = Cl(arg1, arg2, ...)
```

Upon construction of an object, the `Cl:init()` method is called, with arguments `arg1, arg2, ...` passed to it.

Example:

``` lua
Cl = class()
function Cl:init(arg)
	self.x = arg
end
obj = Cl(42)
assert(obj.x == 42)
```

Within this method, `self` is the object being constructed, so no values need to be returned.
If `Cl:init()` does return any values then it will be returned by the constructor after the object.

Example:

``` lua
Cl = class()
function Cl:init()
	return 123
end
obj, extra = Cl()
assert(extra == 123)
```

Class tables have the following fields:

- `Cl.super` = specifies whatever table was passed to parent1, nil if none was.
- `Cl.supers` = specifies the table of all tables passed to parent1.

Class tables have the following methods:

- `Cl:isa(obj)` = returns `true` if obj is a table, and is an instance of class `Cl`.

- `Cl:subclass(...)` = create a subclass of `Cl` and of any other tables specified in `...`.

Notice that the object's metatable is the class table.  This means that any metamethods defined in the class table are functional in the object table:

``` lua
Cl = class()
function Cl.__add(a,b)
	return 'abc'
end
obj1, obj2 = Cl(), Cl()
assert(obj1 + obj2 == 'abc')
```

### coroutine.lua

`coroutine.assertresume(thread, ...)` = This is just like `coroutine.resume` except, upon failure, it automatically includes a stack trace in the error message.

### io.lua

`io.readfile(path)` = Returns the contents of the file as a string.  If an error exists, returns `false` and the error.

`io.writefile(path, data)` = Writes the string in `data` to the file.  If an error occurs, returns `false` and the error, otherwise returns true.

`io.appendfile(path, data)` = Appends the string in `data` to the file.  If an error occurs, returns `false` and the error, otherwise returns true.

`io.readproc(command)` = Runs a process, reads the entirety of its output, returns the output.  If an error occurs, returns `false` and the error.

`local dirname, filename = io.getfiledir(path)` = Returns the directory and the file name of the file at the specified path.

`local pathWithoutExtension, extension = io.getfileext(path)` = Returns the filename up to the extension (without the dot) and the extension.

`file:lock()` = Shorthand for `lfs.lock(filehandle)`.

`file:unlock()` = Shorthand for `lfs.unlock(filehandle)`.

### math.lua

`math.nan` = not-a-number.  `math.nan ~= math.nan`.

`math.e` = Euler's natural base.

`math.atan2` = for Lua version compatability, if this isn't defined then it will be mapped to `math.atan`.

`math.clamp(x, min, max)` = Returns min if x is less than min, max if x is less than max, or x if x is between the two.

`math.sign(x)` = Returns the sign of x.  Note that this is using the convention that `math.sign(0) = 0`.

`math.trunc(x)` = Truncates x, rounding it towards zero.

`math.round(x)` = Rounds x towards the nearest integer.

`math.isnan(x)` = Returns true if x is not-a-number.

`math.isinf(x)`
`math.isinfinite(x)` = Returns true if x is infinite.

`math.isprime(x)` = Returns true if x is a prime number.

`math.factors(x)` = Returns a table of the distinct factors of x.

`math.primeFactorization(x)` = Returns a table of the prime factorization of x.

`math.cbrt(x)` = Returns the cube root of x.

### os.lua

`os.execute(...)` = Attempts to fix the compatability of `os.execute` for Lua 5.1 to match that of Lua 5.2+.

`os.exec(...)` = Prints the command executed, then performs `os.execute`.

`os.fileexists(path)` = Returns true/false whether the associated file exists.

`os.isdir(path)` = Returns true if the file at `path` is a directory.

`os.listdir(path)` = Return an iterator that iterates through all files in a directory.

Example:
``` lua
for fn in os.listdir('.') do
	print(fn)
end
```

`os.rlistdir(path, callback)` = Returns an iterator that recursively iterates through all files in a directory tree.  If `callback` is specified then it is called for each file, and if the callback returns false then the file is not returned, or the sub-directory is not traversed.

`os.mkdir(path[, createParents])` = Create directory.  Set createParents to `true` to create parents as well.

`os.rmdir(path)` = Removes directory.

`os.move(src, dst)` = Move file from src to dst.

`os.home()` = Returns the home directory path.  Queries environment variable HOME or USERPROFILE.

### string.lua

Don't forget that - just as with vanilla Lua - all of these are operable via Lua string metamethods: `("a b c"):split(" ")` gives you `{'a', 'b', 'c'}`.

`string.concat(...)` = `tostring`'s and then `concat`'s all the arguments.

`string.split(str, sep)` = Splits a string, returning a table of the pieces separated by `sep`.  If `sep` is not provided then the string is split into individual characters.

`string.trim(str)` = Returns the string with the whitespace at the beginning and end removed.

`string.bytes(s)` = Returns a table containing the numeric byte values of the characters of the string.

`string.load(str)` = Shorthand for `load` or `loadstring`.  Returns a function of the compiled Lua code, or false and any errors.

`string.csub(str, start, size)` = Returns a substring, where `start` is 0-based and `size` is the length of the substring.

`string.hex(s, [uppercase])` = Return a hex encoding of the string.

`string.unhex(s)` = Returns a string from a hex-encoding.

`string.hexdump(str, columnLength, hexWordSize, spaceEveryNColumns)` = Returns a hex-dump of the string.
- `str` = The string to be hex-dumped.
- `columnLength` = How many columns wide. Default 32.
- `hexWordSize` = How many bytes per word in the hex dump.  Default 1.
- `spaceEveryNColumns` = How often to insert an extra 1-character space.  Default every 8 bytes.

Example of the output:

``` lua
> io.readfile'8x8-24bpp-solidwhite.bmp':hexdump()
00000000  42 4d 3a 01 00 00 00 00  00 00 7a 00 00 00 6c 00  00 00 08 00 00 00 08 00  00 00 01 00 18 00 00 00  BM:.......z...l.................
00000020  00 00 c0 00 00 00 23 2e  00 00 23 2e 00 00 00 00  00 00 00 00 00 00 42 47  52 73 00 00 00 00 00 00  ..�...#...#...........BGRs......
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................................
00000060  00 00 00 00 00 00 00 00  00 00 02 00 00 00 00 00  00 00 00 00 00 00 00 00  00 00 ff ff ff ff ff ff  ..........................������
00000080  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ��������������������������������
000000a0  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ��������������������������������
000000c0  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ��������������������������������
000000e0  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ��������������������������������
00000100  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ��������������������������������
00000120  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ff ff ff ff ff ff  ff ff ��������������������������
```

### table.lua: extensions to the Lua builtin tables

`table.new([table1, table2, ...])`
`table([table1, table2, ...])` = Returns a new table with `table` as its metatable.  If any tables are passed as arguments then this performs a shallow-union of them into the resulting table.

Notice that tables created with `table()` / `table.new()`, i.e. tables with `table` as a metatable, also have `table` as their metatable's `\_\_index` and can use infix notation to call any `table` methods.

`table.unpack(t[, start[, end]])` is assigned to `unpack` for compatability with Lua <5.2.

`table.pack(...)` is made compatible for Lua <5.2.

`table.maxn(t)` is brought back.  Returns the maximum of all number keys of a table.

`table.sort(t, comparator)` is modified to return itself, for ease of chaining infix operations.

`table.append(t, [t1, t2, ...])` = For each additional table provided in the argument (t1, t2, ...), inserts its numeric keys on the end of the first table `t`.

`table.removeKeys(t, k1, k2, ...)` = Remove the specified keys from the table.

`table.map(t, callback)` = Creates and returns a new table. The new table is formed by iterating through the elements of the table, calling `callback(value, key, newtable)`.  If callback returns zero or one values then this is assumed to be the new value for the associated key in the new table.  If callback returns two values then this is assumed to be the new value and key in the new table.

`table.mapi(t, callback)` = Same as `table.map`, but only iterates through the `ipairs` keys.

`table.filter(t, callback)` = Returns a new table formed by iterating through all keys of table `t` using `pairs` and calling the callback.  If the callback returns false then the key/value is excluded from the new table, otherwise it is included.  If the key is numeric then it is inserted using `table.insert`, otherwise it is directly assigned to the new table.

`table.filteri(t, callback)` = Returns a new table formed by iterating through all keys of table `t` using `ipairs` and calling the callback.  If the callback returns false then the key/value is excluded from the new table, otherwise it is inserted using `table.insert`.

`table.keys(t)` = Returns a new table of the enumerated keys of `t` in the order that `pairs` provides them.

`table.values(t)` = Returns a new table of the enumerated values of `t` in the order that `pairs` provides them.

`table.kvpairs(t)` = Returns a new table containing a list of `{[key] = value}` for each key/value pair within `t`.  Intended for use with `next()`.

`table.find(t, value, callback)` = Returns the key and value of the results of finding a value in a table `t` using `pairs`.  If `callback` is not provided then the table is searched for `value`, and matching is tested with ==.  If `callback` is provided then it is called for each i'th value in the table as `callback(ithValue, value)`.  If it returns true then this element is considered to be matching.

`table.findi(t, value, callback)` = Returns the key and value of the results of finding a value in a table `t` using `ipairs`.  If `callback` is not provided then the table is searched for `value`, and matching is tested with ==.  If `callback` is provided then it is called for each i'th value in the table as `callback(ithValue, value)`.  If it returns true then this element is considered to be matching.

`table.insertUnique(t, value, callback)` = Inserts a value into a vable only if it does not already exist in the table (using `table.find`).  If `callback` is provided then it is used as the callback in `table.find`.

`table.removeObject(t, value, callback)` = Removes the specified value from the table, using `table.find`.  Returns a list of the keys that were removed.

`table.sup(t, comparator)` = Returns the maximum value in the table.  `comparator(a,b)` is `a > b` by default.

`table.inf(t, comparator)` = Returns the minimum value in the table.  `comparator(a,b)` is `a < b` by default.

`table.combine(t, callback)` = Combines elements in a table.  The accumulated value is initialized to the first value, and all subsequent values are combined using `callback(accumulatedValue, newValue)`.

`table.sum(t)` = Returns the sum of all values in the table.

`table.product(t)` = Returns the product of all values in the table.

`table.last(t)` = Shorthand for `t[#t]`.

`table.sub(t, start, finish)` = Returns a subset of the table with elements from start to finish.  Uses negative indexes and inclusive ranges, just like `string.sub`.

`table.reverse(t)` = Returns the integer elements of table `t` reversed.

`table.rep(t, n)` = Returns a table of `t` repeated `n` times, just like `string.rep`.

`table.shuffle(t)` = Returns a new table with the integer keys of t shuffled randomly.

`table.pickRandom(t)` = Return a random value in an indexed table.

`table.pickWeighted(t)` = Asssumes collection to pick from are in keys and probabilities are in values.  Returns a random key based.

`table.wrapfor(f, s, var)` = Iterates over a Lua iterator, collecting results in a table.

`table.permutations(t)` = Return an iterator across all permutations of indexed table values.

`table.setmetatable(t)` = `setmetatable`

### number.lua: holds some extra number metatable functionality

This is an extension off of the `ext/math.lua` file, with additional modifications for use as a replacement for the Lua number metatable.

By requiring `ext.number` it will not set up the metatable, don't worry, that all only takes place via `ext.meta`.

The major contribution of this file is `number.tostring`:

`number.tostring(n, base, maxDigits)` = Converts a number to a string, where `base` is the specified base (default 10, non-integers are valid) and `maxDigits` is the maximum number of digits (default 50).
`number.base` = A number value specifying the default base.  This is initialized to 10.
`number.maxdigits` = A number value specifying the default max digts.  This is initialized to 50.

This library doesn't assign `tostring` to `\_\_tostring` by default, but you can accomplish this using the following code:
``` lua
-- assign number to the Lua number metatable
require 'ext.meta'

-- assign number's tostring to __tostring
(function(m) m.__tostring = m.tostring end)(debug.getmetatable(0))
```
Here are some useful shorthands:

`number.hex(x)` = Returns the hexadecimal value of this number.

`number.oct(x)` = Returns the octal value of this number.

`number.bin(x)` = Returns the binary value of this number.

`number.charfor(digit)` = Returns a character for the specified digit.  This function assumes the digit is within the range of its desired base.

`number.todigit(char)` = Returns the numerical value associated with the character.  This is the opposite of `number.charfor`.

`number.alphabets` = A table of tables containing the start and ending of ranges of unicode characters to use for the number to string conversion.

`number.char` = Shorthand assigned to `string.char` for easy infix notation.  Example, when `number` is set to the number metatable, `(97):char()` produces `'a'`.

Example with `number` as Lua number metatables:
``` lua
> require 'ext.meta'
> (10):tostring(3)
101.00000000000000000000000000000000211012220021010120
> (10):tostring(2)
1010.
> (10):tostring(math.pi)
100.01022122221121122001111210201211021100300011111010
```

Yes, go ahead and sum up those digits scaled by their respective powers of pi, see that it does come out to be 10.

Example with `number`, changing the default
``` lua
> require 'ext.meta'
> debug.getmetatable(0).base = 3
> print(10)
10			-- nothing new
> print(tostring(10))
10			-- also nothing new.  tostring is not assigned to __tostring by default.
> print((10):tostring())
101.00000000000000000000000000000000211012220021010120	-- subject to roundoff error
```

### meta.lua: extends off of builtin metatables for all types.

This file modifies most primitive Lua value metatables to add more functionality.

I would not recommend using this file unless it is for small-scale scripts.  The modifications it does are pretty overarching in Lua's behavior, and therefore could have unforeseen side-effects on librarys (though fwiw I haven't seen any yet).
More importantly, it will definitely have implications on the interoperability of any Lua code that you write that is dependent on this behavior.

Metatable changes:

`nil` metatable now supports concat that converts values `tostring` before concatenating them.

`boolean` metatable also now supports concatenation.  Some infix functions are added to represent primitive boolean operations:

- `and_` = `and`, such that `(true):and_(false)` produces `false`.
- `or_` = `or`
- `not_` = `not`
- `xor` = logical XOR, equivalent to `a ~= b` for booleans a and b.
- `implies` = logical 'implies', i.e. `not a or b`.

`number` metatable is assigned to the table in `ext.number` (which is an extension of `ext.math` with some string-serialization additions).

`string` metatable is assigned to the table in `ext.string` and given `tostring` concatenation.

`coroutine` metatable is assigned to the table in `ext.coroutine`.

`function` metatable is given the following operations:

	Default `tostring` concatenation.

	The following binary operators will now work on functions: + - * / % ^
	Example:

``` lua
> function f(x) return 2*x end
> function g(x) return x+1 end
> (f+g)(3)	-- produces (function(x) return 2*x + x+1 end)
10			-- ...and calls it
```

	The following unary operators will now work on functions: - #

``` lua
> function f() return 2 end
> (-f)()
-2
> function f() return 'foo' end
> (#f)()
3
```

The following infix methods are added to functions:

`f:dump()` = Infix shorthand for `string.dump`.

`f:wrap()` = Infix shorthand for `corountine.wrap`.

`f:co()` = Infix shorthand for `coroutine.create`.

`f:index(key)` = Returns a function `g` such that `g(...) == f(...)[k]`.

`f:assign(key, value)` = Returns a function `g` such that `g(...)` executes `f(...)[k] = v`.

`f:compose(g1[, g2, ..., gN])`
`f:o(g1[, g2, ..., gN])` = Returns a function `h` such that `h(x1, x2, ...)` executes `f(g1(g2(...gN( x1, x2, ... ))))`.

`f:compose_n(n, g1[, g2, ..., gM])`
`f:o_n(n, g1[, g2, ..., gM])` = Returns a function 'h' that only replaces the n'th argument with a concatenation of subsequent functions g1...gN.

`f:bind(arg1[, arg2, ..., argN])` = Function curry. Returns a function 'g' that already has arguments arg1 ... argN bound to the first 'n' arguments of 'f'.

`f:bind_n(n, arg1[, arg2, ..., argN])` = Same as above, but start the binding at argument index 'n'.

`f:uncurry(n)` = Uncurry's a function 'n' functions deep.  (a1 -> (a2 -> ... (an -> b))) -> (a1, a2, ..., an -> b).

`f:nargs(n)` = Returns a function that is a duplicate of 'f' but only accepts 'n' arguments.  Very useful with passing builtin functions as callbacks to `table.map` when you don't want extra return values to mess up your resulting table's keys..

`f:swap()` = Returns a new function with the first two parameters swapped.

### range.lua:

`range = require 'ext.range'`
`range(a[, b[, c]])` = A very simple function for creating tables of numeric for-loops.

### reload.lua:

`reload = require 'ext.reload'`
`reload(packagename)` = Removes the package from package.loaded, re-requires it, and returns its result.  Useful for live testing of newly developed features.

### tolua.lua:

`tolua = require 'ext.tolua'`
`tolua(obj[, args])` = Serialization from any Lua value to a string.

args can be any of the following:
- `indent` = Default to 'true', set to 'false' to make results concise.
- `pairs` = The `pairs()` operator to use when iterating over tables.  This defaults to a form of pairs() which iterates over all fields using next().
Set this to your own custom pairs function, or 'pairs' if you would like serialization to respect the `_ _ pairs` metatable (which it does not by default).
- `serializeForType` = A table with keys of lua types and values of callbacks for serializing those types.
- `serializeMetatables` = Set to 'true' to include serialization of metatables.
- `serializeMetatableFunc` = Function to override the default serialization of metatables.
- `skipRecursiveReferences` = Default to 'false', set this to 'true' to not include serialization of recursive references.

### fromlua.lua:

`fromlua = require 'ext.fromlua'`
`fromlua(str)` = De-Serialization from a string to any Lua value.  This is just a small wrapper using `load`.

### cmdline.lua:

`getCmdline = require 'cmdline'`
`cmdline = getCmdline(...)` = builds the table `cmdline` from all command-line arguments.  Here are the rules it follows:
- `cmdline[i]` is the `i`th command-line argument.
- If a command-line argument `k` has no equals sign then `cmdline[k]` is assigned to `true`.
- If a command-line argument has an equals sign present, i.e. is of the form `k=v`, then `cmdline[k]` is assigned the value of `v` if it was evaluated in Lua.
- If evaluating it in Lua produces an error or nil then `cmdline[k]` is assigned the literal string of `v`.
- Don't forget to wrap your complicated assignments in quotations marks.

Notice that, calling `require 'ext'` will also call `getCmdline` on `arg`, producing the global `cmdline`.

### path.lua: path wrapper

`path = require 'ext.path'` = represents an object representing the cwd.

`p = path(pathstr)` returns a new `path` object representing the path at `pathstr`.  Relative paths are appended from the previous `path` object's path.


`p:open(...)` is an alias of `io.open(p.path, ...)`.

`p:read(...)` is an alias of `io.readfile(p.path, ...)`.

`p:write(...)` is an alias of `io.writefile(p.path, ...)`.

`p:append(...)` is an alias of `io.appendfile(p.path, ...)`.

`p:getdir(...)` is an alias of `io.getfiledir(p.path, ...)`, except that it wraps the arguments in a `Path` object.

`p:getext(...)` is an alias of `io.getfileext(p.path, ...)`, except that it wraps the 1st argument in a `Path` object.

`p:setext(newext)` = Returns a path matching `p`s path but with the last extension replaced with `newext`.  If `newext` is `nil` then the last extension is removed.


`p:remove(...)` is an alias of `os.remove(p.path, ...)`.

`p:mkdir(...)` is an alias of `os.mkdir(p.path, ...)`.

`p:dir(...)` is an alias of `os.listdir(p.path, ...)`, except that it wraps the arguments in a `Path` object.

`p:exists(...)` is an alias of `os.fileexists(p.path, ...)`.

`p:isdir(...)` is an alias of `os.isdir(p.path, ...)`.

`p:rdir(...)` is an alias of `os.rlistdir(p.path, ...)`, except that it wraps the arguments in a `Path` object.


`p:attr(...)` is an alias of `lfs.attributes(p.path, ...)`.

`p:symattr(...)` is an alias of `lfs.symlinkattributes(p.path, ...)`.

`p:cd(...)` is an alias of `lfs.chdir(p.path, ...)`.

`p:link(...)` is an alias of `lfs.link(p.path, ...)`.

`p:setmode(...)` is an alias of `lfs.setmode(p.path, ...)`.

`p:touch(...)` is an alias of `lfs.touch(p.path, ...)`.

`p:lockdir(...)` is an alias of `lfs.lock_dir(p.path, ...)`.


`p:cwd()` returns the absolute cwd path, as a Path object.  This is an alias of `lfs.currendir()` if available, or `io.readproc'pwd'` on Linux or `io.readproc'cd'` on Windows.

`p:abs()` returns the absolute form of the path, as a Path object.

### debug.lua

use this like so:

Add debug code to your file like so:

`--DEBUG:print"This only runs in debug"`

then run your code like so:

`lua -lext.debug ...`

or

`lua -e "require'ext.debug'" ...`

...and any single-line comments in any code that start with `--DEBUG:` will be uncommented.

You can also add specific runtime tags:

`--DEBUG(mytag):print"This only runs with debug'mytag'`

`lua -e "require'ext.debug''tag==[[mytag]]'`

... or log-levels (the default level is 1):

`--DEBUG(@2):print"This only runs with debug'level<=2'`

`lua -e "require'ext.debug''level<=2'`

... or both:

`--DEBUG(mytag@3):print"This only runs with debug'tag=="mytag" and level<=3'`

`lua -e "require'ext.debug''tag==[[mytag]] and level<=3'`

Debug queries can be any Lua expression using the variables `source`, `line`, `tag`, and `log`.

### gcmem.lua:

Provides FFI-based functions for manually or automatically allocating and freeing memory.

WIP due to crashing in LuaJIT when you run `ptr = ffi.cast('T*', ptr)` and haven't bound `ptr` anywhere else.

# NOTICE:
- path.lua will optionally use luafilesystem, if available.  My own luafilesystem fork is preferred: https://github.com/thenumbernine/luafilesystem since it maintains a single copy of ffi cdefs within my lua-ffi-bindings library.
- gcmem.lua depends on ffi, particularly some ffi headers of stdio found in my lua-ffi-bindings project: https://github.com/thenumbernine/lua-ffi-bindings
