# zzterm [![builds.sr.ht status](https://builds.sr.ht/~mna/zzterm.svg)](https://builds.sr.ht/~mna/zzterm?) [![GoDoc](https://godoc.org/git.sr.ht/~mna/zzterm?status.svg)](http://godoc.org/git.sr.ht/~mna/zzterm) [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/git.sr.ht/~mna/zzterm)

Package zzterm efficiently reads and decodes terminal input keys and mouse events
without any memory allocation. It is intended to be used with a terminal set in
raw mode as its `io.Reader`. See the [package documentation][godoc] for details,
API reference and usage example (alternatively, on [pkg.go.dev][pgd]).

You can also check out [zztermtest][zztt] for more usage examples and how to
efficiently print output to an `io.Writer` (with a zero-allocation "echo"
program example).

* Canonical repository: https://git.sr.ht/~mna/zzterm
* Issues: https://todo.sr.ht/~mna/zzterm
* Builds: https://builds.sr.ht/~mna/zzterm

Note that at the moment, zzterm is only tested on macOS and linux. Mouse support
is through the Xterm X11 mouse protocol with SGR enabled, the terminal has to
support that mode for mouse (and focus) key events to be emitted.

## See Also

Similar Go packages:

* [tj/go-terminput](https://github.com/tj/go-terminput): similar small scope
with a focus on input key decoding, no mouse support.
* [gdamore/tcell](https://github.com/gdamore/tcell): larger scope, handles
output too, colors, raw mode, etc.
* [nsf/termbox-go](https://github.com/nsf/termbox-go): larger scope like tcell.

## Benchmarks

The input processing is typically in the hot path of a terminal application. Zzterm
is quite fast and does not allocate - not when decoding standard keys, not when
decoding escape sequences, neither when decoding mouse events.

```
benchmark                                    iter     time/iter   bytes alloc        allocs
---------                                    ----     ---------   -----------        ------
BenchmarkInput_ReadKey/a-4               61804756   18.40 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/B-4               66716232   17.90 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/1-4               62950432   18.60 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/\x00-4            65492827   18.20 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/ø-4               60368734   19.90 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/👪-4              57783043   20.60 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/平-4              57067489   20.80 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/\x1b[B-4          26063134   45.90 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/\x1b[1;2C-4       26355364   45.40 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/\x1b[I-4          26530273   44.40 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey/\x1b[<35;1;2M-4   21740397   55.30 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey_Bytes-4           49141444   24.40 ns/op        0 B/op   0 allocs/op
BenchmarkInput_ReadKey_Mouse-4           19961526   60.70 ns/op        0 B/op   0 allocs/op
```

## License

The [BSD 3-Clause license][bsd].

[bsd]: http://opensource.org/licenses/BSD-3-Clause
[godoc]: http://godoc.org/git.sr.ht/~mna/zzterm
[pgd]: https://pkg.go.dev/git.sr.ht/~mna/zzterm
[zztt]: https://git.sr.ht/~mna/zztermtest
