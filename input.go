package zzterm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

// Input reads input keys from a reader and returns the key pressed.
type Input struct {
	buf   []byte
	lastn int
	lastm MouseEvent

	// immutable after NewInput
	esc   map[string]Key
	mouse bool
	focus bool // only required to add the focus-related escape sequences in esc map
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
// key with type KeyMouse, and the mouse information can be retrieved by
// calling Input.Mouse before the next call to Input.RedKey. It is the
// responsibility of the caller to enable mouse tracking for the terminal
// represented by the io.Reader passed to ReadKey, and SGR Mouse Mode must be
// enabled. Not all tracking modes are supported, see MouseEventType constants
// for supported modes. As a convenience, the package provides the EnableMouse
// and DisableMouse functions to enable and disable mouse tracking on a
// terminal represented by an io.Writer.
//
// Only X11 xterm mouse protocol in SGR mouse mode is supported. This should
// be widely supported by any recent terminal with mouse support.  See
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
func WithMouse() Option {
	return func(i *Input) error {
		i.mouse = true
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

// Option defines the function signatures for options to apply when
// creating a new Input.
type Option func(*Input) error

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

// Mouse returns the mouse event corresponding to the last key of type KeyMouse.
// It should be called only after a key of type KeyMouse has been received from
// ReadKey, and before any other call to ReadKey.
func (i *Input) Mouse() MouseEvent {
	return i.lastm
}

const sgrMouseEventPrefix = "\x1b[<"

// ReadKey reads a key from r.
func (i *Input) ReadKey(r io.Reader) (Key, error) {
	i.lastn = 0
	n, err := r.Read(i.buf)
	if err != nil || n == 0 {
		return 0, err
	}
	i.lastn = n
	buf := i.buf[:n]

	c, sz := utf8.DecodeRune(buf)
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
		if i.mouse && bytes.HasPrefix(buf, []byte(sgrMouseEventPrefix)) {
			if k := i.decodeMouseEvent(); k.Type() == KeyMouse {
				return k, nil
			}
		}
		// NOTE: important to use the string conversion exactly like that,
		// inside the brackets of the map key - the Go compiler optimizes
		// this to avoid any allocation.
		if key, ok := i.esc[string(buf)]; ok {
			return key, nil
		}
		// if this is an unknown escape sequence, return KeyESCSeq and the
		// caller may get the uninterpreted sequence from i.Bytes.
		return keyFromTypeMod(KeyESCSeq, ModNone), nil
	}
	return Key(c), nil
}

// returns either a KeyMouse key, or a KeyESCSeq if it can't properly decode
// the mouse event.
func (i *Input) decodeMouseEvent() Key {
	// the prefix has already been validated, strip it from the working buffer
	buf := i.buf[len(sgrMouseEventPrefix):i.lastn]
	if len(buf) < 6 {
		// 2 semicolons, trailing m/M, at least one byte in each section
		return keyFromTypeMod(KeyESCSeq, ModNone)
	}

	// the final character must be M (key press) or m (key release)
	var pressed bool
	switch buf[len(buf)-1] {
	case 'M':
		pressed = true
	case 'm':
	default:
		return keyFromTypeMod(KeyESCSeq, ModNone)
	}
	buf = buf[:len(buf)-1]

	// extract the 3 parameter numbers
	var nums [3]uint16
	for i := 0; i < 2; i++ {
		// must have 3 semicolon-separated parts, so 2 semicolons
		ix := bytes.IndexByte(buf, ';')
		if ix < 0 {
			return keyFromTypeMod(KeyESCSeq, ModNone)
		}
		num, err := parseUintBytes(buf[:ix])
		if err != nil {
			return keyFromTypeMod(KeyESCSeq, ModNone)
		}
		nums[i] = num
		buf = buf[ix+1:]
	}
	// process the 3rd (remaining) number
	num, err := parseUintBytes(buf)
	if err != nil {
		return keyFromTypeMod(KeyESCSeq, ModNone)
	}
	nums[2] = num

	// decode the button event (first number)
	mod := Mod(nums[0]) & modMouseEvent
	btn := int(nums[0] & 0b_0000_0011) // this gives a number between 0-3, but 3 is not a button
	add := int((nums[0] & 0b_1100_0000) >> 4)
	btn += add // button is between 0-11
	// detect if it is a mouse move only - i.e. no button pressed
	if (btn == 0b_0011 && (nums[0]&0b_0010_0000 != 0)) || btn == 3 {
		btn = 0
	} else if btn < 3 {
		btn++ // because 0-1-2 values are for IDs 1-2-3
	}

	i.lastm = MouseEvent{byte(btn), pressed, nums[1], nums[2]}

	//fmt.Printf("%d - %d - %d (pressed? %t; modifier: %s)\r\n", nums[0], nums[1], nums[2], !btnRelease, mod)
	return keyFromTypeMod(KeyMouse, mod)
}

var (
	errInvalidUint = errors.New("invalid uint number")
)

// parse a uint16 number in base 10 from the provided bytes. If the value is
// greater than maxUint16, it returns maxUint16 (not an error).
func parseUintBytes(b []byte) (uint16, error) {
	const (
		maxUint16 = 1<<16 - 1
	)

	if len(b) == 0 {
		return 0, errInvalidUint
	}

	var n uint32
	for i := 0; i < len(b); i++ {
		var v byte
		d := b[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		default:
			return 0, errInvalidUint
		}

		n *= 10
		n += uint32(v)

		if n > maxUint16 {
			return maxUint16, nil
		}
	}
	return uint16(n), nil
}
