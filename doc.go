// Package zzterm efficiently reads and decodes terminal input keys and mouse
// events without any memory allocation. It is intended to be used with a
// terminal set in raw mode as its io.Reader. Setting the terminal in raw
// mode is not handled by this package, there are a number of Go packages that
// can do this (see the example).
//
// Basic usage
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
//         	   	   log.Panic(err)
//         	   }
//
//         	   switch k.Type() {
//         	   case zzterm.KeyRune:
//                 // k.Rune() returns the rune
//         	   case zzterm.KeyESC, zzterm.KeyCtrlC:
//                 // quit on ESC or Ctrl-C
//         	   	   return
//         	   }
//         }
//     }
//
// Mouse and focus events
//
// Mouse events are supported through the Xterm X11 mouse protocol in SGR
// mode, which is a complex way to call the "modern" handling of mouse events
// (beyond the limits of 223 for mouse position coordinates in the old protocol).
// This should be widely supported by modern terminals, but the tracking of mouse
// events must be enabled on the terminal so that the escape sequences get sent
// to zzterm. It is the responsibility of the caller to enable this (with SGR
// mode) before using Input.ReadKey, but as a convenience zzterm provides the
// EnableMouse and DisableMouse functions:
//
//     t, err := term.Open("/dev/tty", term.RawMode)
//     ...
//     defer t.Restore()
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
//         ...
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
//     ...
//     case zzterm.KeyFocusIn, zzterm.KeyFocusOut:
//         // terminal has gained/lost focus
//     ...
//
package zzterm // import "git.sr.ht/~mna/zzterm"
