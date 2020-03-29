package zzterm

import (
	"errors"
	"io"
	"unicode/utf8"
)

// Input reads input keys from a reader and returns the key pressed.
type Input struct {
	buf   []byte
	lastn int
	esc   map[string]Key
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
		i.esc = defaultEsc
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
