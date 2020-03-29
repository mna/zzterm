package zzterm

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

// Input reads input keys from a reader and returns the key pressed.
type Input struct {
	buf   []byte
	lastn int
	esc   map[string]Key
	mouse MouseEventType // 0=no mouse event
	focus bool           // false=no focus event
}

// Option defines the function signatures for options to apply when
// creating a new Input.
type Option func(*Input) error

// WithESCSeq sets the terminfo-like map that defines the interpretation of
// escape sequences as special keys. The map has the same field names as those
// used in the github.com/gdamore/tcell/terminfo package for the Terminfo
// struct.  Only the fields starting with "Key" are supported, and only the key
// sequences starting with ESC (0x1b) are considered.
//
// If nil is passed (or if the option is not specified), common default values
// are used. To prevent any translation of escape sequences to special keys,
// pass a non-nil empty map. All escape sequences will be returned as KeyESCSeq
// and the raw bytes of the sequence can be retrieved by calling Input.Bytes.
//
// If you want to use tcell's terminfo definitions directly, you can use the
// helper function FromTerminfo that accepts an interface{} and returns a
// map[string]string that can be used here, in order to avoid adding tcell as a
// dependency, and to support any value that marshals to JSON the same way as
// tcell/terminfo. Note, however, that tcell manually patches some escape
// sequences in its code, overriding the terminfo definitions in some cases. It
// is up to the caller to ensure the mappings are correct, zzterm does not
// apply any patching.
//
// See https://github.com/gdamore/tcell/blob/8ec73b6fa6c543d5d067722c0444b07f7607ba2f/tscreen.go#L337-L367
func WithESCSeq(tinfo map[string]string) Option {
	return func(i *Input) error {
		i.esc = escFromTerminfo(tinfo)
		return nil
	}
}

// MouseEventType represents a type of mouse events.
type MouseEventType int

// List of supported mouse event types.
const (
	MouseButton MouseEventType = iota + 1 // CSI ? 1000 h
	_                                     // unsupported but reserved, CSI ? 1001 h
	_                                     // unsupported but reserved, CSI ? 1002 h
	MouseAny                              // CSI ? 1003 h
)

// EnableMouse sends the Control Sequence Introducer (CSI) function to
// w to enable tracking of the specified mouse event type in SGR mode.
func EnableMouse(w io.Writer, eventType MouseEventType) error {
	code := eventType + 1000 - 1
	_, err := fmt.Fprintf(w, "\x1b[?%d;1006h", code)
	return err
}

// DisableMouse sends the Control Sequence Introducer (CSI) function to
// w to disable tracking of the specified mouse event type and to disable
// SGR mode.
func DisableMouse(w io.Writer, eventType MouseEventType) error {
	code := eventType + 1000 - 1
	_, err := fmt.Fprintf(w, "\x1b[?%d;1006l", code)
	return err
}

// EnableFocus sends the Control Sequence Introducer (CSI) function to
// w to enable sending focus escape sequences.
func EnableFocus(w io.Writer) error {
	_, err := fmt.Fprint(w, "\x1b[?1004h")
	return err
}

// DisableFocus sends the Control Sequence Introducer (CSI) function to
// w to disable sending focus escape sequences.
func DisableFocus(w io.Writer) error {
	_, err := fmt.Fprint(w, "\x1b[?1004l")
	return err
}

// WithMouse enables mouse event reporting.  Such events will be reported as a
// key with type KeyMouse. It is the responsibility of the caller to enable
// mouse tracking for the terminal represented by the io.Reader passed to
// ReadKey, and SGR Mouse Mode must be enabled. Not all tracking modes are
// supported, see MouseEventType constants for supported modes. As a
// convenience, the package provides the EnableMouse and DisableMouse
// functions to enable and disable mouse tracking on a terminal represented by
// an io.Writer.
//
// Only X11 xterm mouse protocol in SGR mouse mode is supported. This should
// be widely supported by any recent terminal with mouse support.  See
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
func WithMouse(eventType MouseEventType) Option {
	return func(i *Input) error {
		i.mouse = eventType
		return nil
	}
}

// WithFocus enables reporting of focus in and focus out events when the
// terminal gets and loses focus. Such events will be reported as a key with
// type KeyFocusIn or KeyFocusOut. It is the responsibility of the caller to
// enable focus tracking for the terminal represented by the io.Reader passed
// to ReadKey. As a convenience, the package provides the EnableFocus and
// DisableFocus functions to enable and disable focus tracking on a terminal
// represented by an io.Writer.  See
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-FocusIn_FocusOut
func WithFocus() Option {
	return func(i *Input) error {
		i.focus = true
		return nil
	}
}

// NewInput creates an Input ready to use. Call Input.ReadKey to read a single
// key from an io.Reader - typically a terminal file descriptor set in raw mode.
// The translation of escape sequences to special keys is controlled by the
// WithESCSeq option.
func NewInput(opts ...Option) (*Input, error) {
	i := &Input{
		buf: make([]byte, 128),
	}
	for _, o := range opts {
		if err := o(i); err != nil {
			return nil, err
		}
	}
	if i.esc == nil {
		i.esc = cloneEscMap(defaultEsc)
	}
	if i.focus {
		addFocusESCSeq(i.esc)
	}

	return i, nil
}

// Bytes returns the uninterpreted bytes from the last key read. The bytes
// are valid only until the next call to ReadKey and should not be modified.
func (i *Input) Bytes() []byte {
	if i.lastn <= 0 {
		return nil
	}
	return i.buf[:i.lastn:i.lastn]
}

// ReadKey reads a key from r.
func (i *Input) ReadKey(r io.Reader) (Key, error) {
	i.lastn = 0
	n, err := r.Read(i.buf)
	if err != nil || n == 0 {
		return 0, err
	}
	i.lastn = n

	c, sz := utf8.DecodeRune(i.buf[:n])
	if c == utf8.RuneError && sz < 2 {
		return 0, errors.New("invalid rune")
	}

	// if c is a control character (if n == 1 so that if an escape
	// sequence is read, it does not return immediately with just ESC)
	if n == 1 && (KeyType(c) <= KeyUS || KeyType(c) == KeyDEL) {
		return keyFromTypeMod(KeyType(c), ModNone), nil
	}

	// translate escape sequences
	if KeyType(c) == KeyESC {
		// NOTE: important to use the string conversion exactly like that,
		// inside the brackets of the map key - the Go compiler optimizes
		// this to avoid any allocation.
		if key, ok := i.esc[string(i.buf[:n])]; ok {
			return key, nil
		}
		// if this is an unknown escape sequence, return KeyESCSeq and the
		// caller may get the uninterpreted sequence from i.Bytes.
		return keyFromTypeMod(KeyESCSeq, ModNone), nil
	}
	return Key(c), nil
}
