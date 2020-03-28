package zzterm

import (
	"errors"
	"io"
	"unicode/utf8"
)

// Input reads input keys from a reader and returns the key pressed.
type Input struct {
	buf []byte
}

// NewInput creates an Input ready to use.
func NewInput() *Input {
	return &Input{buf: make([]byte, 16)}
}

// ReadKey reads a key from r.
func (i *Input) ReadKey(r io.Reader) (Key, error) {
	n, err := r.Read(i.buf)
	if err != nil {
		return 0, err
	}

	c, sz := utf8.DecodeRune(i.buf[:n])
	if c == utf8.RuneError && sz < 2 {
		return 0, errors.New("invalid rune")
	}

	// if c is a control character
	if n == 1 && (KeyType(c) <= KeyUS || KeyType(c) == KeyDEL) {
		return keyFromTypeMod(KeyType(c), ModNone), nil
	}

	// sequences
	key, ok := keySequences[string(i.buf[:n])]
	if ok {
		var mod Mod
		if key > KeyDEL {
			// key with modifiers set
			mod = Mod((key & 0b_0111_0000) >> 4)
			key = (key & 0b_0000_1111) + KeyRune + 1
		}
		return keyFromTypeMod(key, mod), nil
	}

	return Key(c), nil
}

// Key represents a single key. It contains the key type,
// the key modifier flags and the rune itself in a compact
// form. Use the Rune, Type and Mod methods to get information
// on the key.
type Key uint32

// The key format is:
// * if the key is control character or a special key, the sign bit
//   is set to negative and the first (lower) byte is the Type and
//   the second byte is the Mod.
// * otherwise, the (positive) value itself is the rune.
//
// There is no Mod set for a standard rune because generally in a raw
// mode terminal we cannot tell if Shift or Ctrl or some other modifier
// key was pressed to generate the rune.
func keyFromTypeMod(t KeyType, m Mod) Key {
	k := Key(m) << 8
	k |= Key(t)
	k |= (1 << 31)
	return k
}

// TODO: String for Key

func (k Key) Rune() rune {
	r := rune(k)
	if r < 0 {
		return -1
	}
	return rune(k)
}

func (k Key) Type() KeyType {
	if r := rune(k); r >= 0 {
		return KeyRune
	}
	// otherwise extract type from the first 8 bytes
	return KeyType(k & 0xFF)
}

func (k Key) Mod() Mod {
	if r := rune(k); r >= 0 {
		return ModNone
	}
	return Mod((k >> 8) & 0xFF)
}

// Mod represents a key modifier such as pressing alt or ctrl.
type Mod byte

// List of modifier flags.
const (
	ModShift Mod = 1 << iota
	ModCtrl
	ModAlt
	ModNone Mod = 0
)

// KeyType represents the type of key.
type KeyType byte

// Supported key types - the following match the ASCII value.
const (
	KeyNUL KeyType = iota
	KeySOH
	KeySTX
	KeyETX
	KeyEOT
	KeyENQ
	KeyACK
	KeyBEL
	KeyBS
	KeyTAB
	KeyLF
	KeyVT
	KeyFF
	KeyCR
	KeySO
	KeySI
	KeyDLE
	KeyDC1
	KeyDC2
	KeyDC3
	KeyDC4
	KeyNAK
	KeySYN
	KeyETB
	KeyCAN
	KeyEM
	KeySUB
	KeyESC
	KeyFS
	KeyGS
	KeyRS
	KeyUS
	KeyRune // covers ASCII 32-126 + any other unicode code point

	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyInsert
	KeyBacktab
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDn
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20

	KeyDEL KeyType = 0x7f

	// keys with modifiers - values above 127 with the bits 6-5-4 indicating the
	// modifiers and the low 4 bits the actual key (up to 16 values,
	// corresponding to  KeyType(KeyRune + 1 + low 4 bits value).
	keyShiftLeft  = 0b_1001_0000
	keyShiftRight = 0b_1001_0001
	keyAltLeft    = 0b_1100_0000
	keyAltRight   = 0b_1100_0001
)

var keySequences = map[string]KeyType{
	"\x1b[A":     KeyUp,
	"\x1b[B":     KeyDown,
	"\x1b[C":     KeyRight,
	"\x1b[D":     KeyLeft,
	"\x1b[2~":    KeyInsert,
	"\x1b[3~":    KeyDelete,
	"\x1b[Z":     KeyBacktab,
	"\x1bOH":     KeyHome,
	"\x1bOF":     KeyEnd,
	"\x1b[5~":    KeyPgUp,
	"\x1b[6~":    KeyPgDn,
	"\x1bOP":     KeyF1,
	"\x1bOQ":     KeyF2,
	"\x1bOR":     KeyF3,
	"\x1bOS":     KeyF4,
	"\x1b[15~":   KeyF5,
	"\x1b[17~":   KeyF6,
	"\x1b[18~":   KeyF7,
	"\x1b[19~":   KeyF8,
	"\x1b[20~":   KeyF9,
	"\x1b[21~":   KeyF10,
	"\x1b[23~":   KeyF11,
	"\x1b[24~":   KeyF12,
	"\x1b[1;2P":  KeyF13,
	"\x1b[1;2Q":  KeyF14,
	"\x1b[1;2R":  KeyF15,
	"\x1b[1;2S":  KeyF16,
	"\x1b[15;2~": KeyF17,
	"\x1b[17;2~": KeyF18,
	"\x1b[18;2~": KeyF19,
	"\x1b[19;2~": KeyF20,
	"\x1b[1;2D":  keyShiftLeft,
	"\x1b[1;2C":  keyShiftRight,
	"\x1bb":      keyAltLeft,
	"\x1bf":      keyAltRight,
}
