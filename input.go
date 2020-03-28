package zzterm

import (
	"errors"
	"io"
	"unicode/utf8"
)

// Input reads input keys from a reader and returns the key pressed.
type Input struct {
	buf []byte
	esc map[string]Key
}

// NewInput creates an Input ready to use. The tinfo map is an optional
// terminfo-like map defining the supported escape codes to translate some
// special keys. The map has the same field names as those used in the
// github.com/gdamore/tcell/terminfo package for the Terminfo struct.
// Only the fields starting with "Key" are supported, and only the
// key sequences starting with ESC (0x1b) are considered.
//
// If nil is passed, common default values are used.
//
// If you want to use tcell's terminfo definitions directly, you can
// use the helper function FromTerminfo that accepts an interface{}
// and returns a map[string]string that can be used here,
// in order to avoid adding tcell as a dependency, and to support any
// value that marshals to JSON the same way as tcell/terminfo.
func NewInput(tinfo map[string]string) *Input {
	return &Input{
		buf: make([]byte, 128),
		esc: escFromTerminfo(tinfo),
	}
}

// ReadKey reads a key from r.
func (i *Input) ReadKey(r io.Reader) (Key, error) {
	n, err := r.Read(i.buf)
	if err != nil || n == 0 {
		return 0, err
	}

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
	}
	return Key(c), nil
}
