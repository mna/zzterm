// Package zzterm efficiently reads and decodes terminal input keys and mouse
// events without any memory allocation. It is intended to be used with a
// terminal set in raw mode as its io.Reader. Setting the terminal in raw
// mode is not handled by this package, there are a number of Go packages that
// can do this (see the example).
//
// Basic usage
//
// Set the terminal in raw mode, use NewInput to create the input key reader
// and read from the terminal:
//
//     func main() {
//         // set the terminal in raw mode, e.g. with github.com/pkg/term
//         t, err := term.Open("/dev/tty", term.RawMode)
//         if err != nil {
//             log.Panic(err)
//         }
//         defer t.Restore()
//
//         input := zzterm.NewInput()
//         for {
//         	   k, err := input.ReadKey(t)
//         	   if err != nil {
//                 log.Panic(err)
//         	   }
//
//         	   switch k.Type() {
//         	   case zzterm.KeyRune:
//                 // k.Rune() returns the rune
//         	   case zzterm.KeyESC, zzterm.KeyCtrlC:
//                 // quit on ESC or Ctrl-C
//                 return
//             }
//         }
//     }
//
// Mouse and focus events
//
// Mouse events are supported through the Xterm X11 mouse protocol in SGR
// mode, which is a complex way to call the "modern" handling of mouse events [1]
// (beyond the limits of 223 for mouse position coordinates in the old protocol).
// This should be widely supported by modern terminals, but the tracking of mouse
// events must be enabled on the terminal so that the escape sequences get sent
// to zzterm. It is the responsibility of the caller to enable this (with SGR
// mode) before using Input.ReadKey, but as a convenience zzterm provides the
// EnableMouse and DisableMouse functions:
//
//     t, err := term.Open("/dev/tty", term.RawMode)
//     // ...
//     defer t.Restore()
//
//     // Mouse events can be enabled only to report button presses (zzterm.MouseButton)
//     // or any mouse event (including mouse moves, zzterm.MouseAny).
//     zzterm.EnableMouse(t, zzterm.MouseAny)
//     defer zzterm.DisableMouse(t, zzterm.MouseAny)
//
// And then mouse events will be reported (if supported by the terminal):
//
//     // The WithMouse option must be set to decode the mouse events, otherwise
//     // they would be reported as uninterpreted KeyESCSeq (escape sequence).
//     input := zzterm.NewInput(zzterm.WithMouse())
//     for {
//         switch k.Type() {
//         case zzterm.KeyRune:
//             // k.Rune() returns the rune
//         case zzterm.KeyMouse:
//             // k.Mod() returns the modifier flags (e.g. Shift) pressed during the event
//             // input.Mouse() returns the mouse information, coordinates 1,1 is top-left
//         // ...
//         }
//     }
//
// It works similarly to enable reporting focus in/out of the terminal:
//
//     zzterm.EnableFocus(t)
//     defer zzterm.DisableFocus(t)
//
//     // The WithFocus option must be set to decode the focus events, otherwise
//     // they would be reported as uninterpreted KeyESCSeq (escape sequence).
//     input := zzterm.NewInput(zzterm.WithMouse(), zzterm.WithFocus())
//     for {
//         // ...
//         case zzterm.KeyFocusIn, zzterm.KeyFocusOut:
//             // terminal has gained/lost focus
//         // ...
//     }
//
// Terminfo
//
// Different terminals sometimes understand different escape sequences to interpret
// special keys such as function keys (F1, F2, etc.) and arrows. That configuration
// is part of the terminfo database (at least on Unix-like systems). While zzterm does
// not read the terminfo database itself, it supports specifying a map of values where
// the key is the name of the special key and the value is the escape sequence that
// should map to this key.
//
//     escSeq := map[string]string{"KeyDown": "\x1b[B"}
//     input := zzterm.NewInput(zzterm.WithESCSeq(escSeq))
//
// The github.com/gdamore/tcell repository has a good number of terminal configurations
// described in a Go struct and accessible via terminfo.LookupTermInfo [2]. To enable
// reuse of this, zzterm provides the FromTerminfo function to convert from those
// structs to the supported map format. It is the responsibility of the caller to
// detect the right terminfo to use for the terminal.
//
//     ti, err := terminfo.LookupTerminfo("termite")
//     // handle error
//     input := zzterm.NewInput(zzterm.WithESCSeq(zzterm.FromTerminfo(ti)))
//
// Note, however, that the tcell package patches those terminfo descriptions before use
// due to some inconsistencies in behaviour - using the raw terminfo definitions may
// not always work as expected [3].
//
// When no WithESCSeq option is provided (or if a nil map is passed), then a default
// mapping is used. If a non-nil but empty map is provided, then any escape sequence
// translation will be disabled (except for mouse and focus events if enabled), and all
// such sequences will be read as keys of type KeyESCSeq. The input.Bytes method can
// then be called to inspect the raw bytes of the sequence.
//
// [1]: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
// [2]: https://godoc.org/github.com/gdamore/tcell/terminfo#LookupTerminfo
// [3]: https://github.com/gdamore/tcell/blob/8ec73b6fa6c543d5d067722c0444b07f7607ba2f/tscreen.go#L337-L367
//
package zzterm // import "git.sr.ht/~mna/zzterm"
