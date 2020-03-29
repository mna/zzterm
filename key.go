package zzterm

import "fmt"

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

// String returns the string representation of k.
func (k Key) String() string {
	if k.Type() == KeyRune {
		return fmt.Sprintf("Key(%#U)", k.Rune())
	}

	var flags string
	mod := k.Mod()
	if mod&ModCtrl != 0 {
		flags += "⌃"
	}
	if mod&ModShift != 0 {
		flags += "⇧"
	}
	if mod&ModAlt != 0 {
		flags += "⎇"
	}
	if mod&ModMeta != 0 {
		flags += "⌥"
	}
	if flags != "" {
		flags += " "
	}
	return fmt.Sprintf("Key(%sx%02x)", flags, k.Type())
}

// Rune returns the rune corresponding to this key. It returns -1
// if the KeyType is not KeyRune.
func (k Key) Rune() rune {
	r := rune(k)
	if r < 0 {
		return -1
	}
	return rune(k)
}

// Type returns the KeyType for this key.
func (k Key) Type() KeyType {
	if r := rune(k); r >= 0 {
		return KeyRune
	}
	// otherwise extract type from the first 8 bytes
	return KeyType(k & 0xFF)
}

// Mod returns the key modifier flags set for this key.
func (k Key) Mod() Mod {
	if r := rune(k); r >= 0 {
		return ModNone
	}
	return Mod((k >> 8) & 0xFF)
}

// Mod represents a key modifier such as pressing alt or ctrl.
// Detection of such flags is limited.
type Mod byte

// List of modifier flags.
const (
	ModShift Mod = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
	ModNone Mod = 0
)

// KeyType represents the type of key.
type KeyType byte

// List of supported key types.
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
	KeyRune // covers ASCII 32-126 + any other unicode code point - from this point the key type does not match ASCII values
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
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyF26
	KeyF27
	KeyF28
	KeyF29
	KeyF30
	KeyF31
	KeyF32
	KeyF33
	KeyF34
	KeyF35
	KeyF36
	KeyF37
	KeyF38
	KeyF39
	KeyF40
	KeyF41
	KeyF42
	KeyF43
	KeyF44
	KeyF45
	KeyF46
	KeyF47
	KeyF48
	KeyF49
	KeyF50
	KeyF51
	KeyF52
	KeyF53
	KeyF54
	KeyF55
	KeyF56
	KeyF57
	KeyF58
	KeyF59
	KeyF60
	KeyF61
	KeyF62
	KeyF63
	KeyF64
	KeyHelp
	KeyExit
	KeyClear
	KeyCancel
	KeyPrint // 112

	KeyDEL KeyType = 127
)

// List of some aliases to the key types. The KeyCtrl... constants
// match the ASCII keys at the same position (e.g. KeyCtrlSpace is
// KeyNUL, KeyCtrlLeftSq is KeyESC, etc.).
const (
	KeyCtrlSpace KeyType = iota
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
	KeyCtrlLeftSq
	KeyCtrlBackslash
	KeyCtrlRightSq
	KeyCtrlCarat
	KeyCtrlUnderscore

	KeyBackspace = KeyBS
	KeyEscape    = KeyESC
	KeyEnter     = KeyCR
)
