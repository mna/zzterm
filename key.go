package zzterm

import (
	"fmt"
	"strconv"
)

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

	flags := k.Mod().String()
	if flags != "" {
		flags += " "
	}
	return fmt.Sprintf("Key(%s%s)", flags, k.Type())
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

// String returns the string representation of m.
func (m Mod) String() string {
	var flags string
	if m&ModCtrl != 0 {
		flags += "⌃"
	}
	if m&ModShift != 0 {
		flags += "⇧"
	}
	if m&ModAlt != 0 {
		flags += "⎇"
	}
	if m&ModMeta != 0 {
		flags += "⌥"
	}
	return flags
}

// List of modifier flags. Values of Shift, Meta and Ctrl are the same
// as for the xterm mouse tracking.
// See https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Normal-tracking-mode
const (
	_        Mod = 1 << iota
	ModAlt       // 2
	ModShift     // 4
	ModMeta      // 8
	ModCtrl      // 16
	ModNone  Mod = 0

	modMouseEvent = ModShift | ModMeta | ModCtrl // 0b_0001_1100
)

// MouseEvent describes a KeyMouse key type. While the Key returned
// by Input.ReadKey has the modifier flags information, the mouse-related
// properties are defined by the MouseEvent type.
type MouseEvent struct {
	buttonID byte
	pressed  bool
	x, y     uint16
}

// String returns the string representation of a mouse event.
func (m MouseEvent) String() string {
	var btn string
	if m.ButtonID() > 0 {
		state := "⇑"
		if m.ButtonPressed() {
			state = "⇓"
		}
		btn = fmt.Sprintf("%s%02d ", state, m.ButtonID())
	}
	x, y := m.Coords()
	return fmt.Sprintf("Mouse(%sx:%d y:%d)", btn, x, y)
}

// ButtonID returns the button pressed during the mouse event, starting
// at 1. A ButtonID of 0 means that no button was pressed - i.e. this is
// a mouse move event without any button pressed. Up to 11 buttons are
// supported by the X11 mouse protocol.
func (m MouseEvent) ButtonID() int {
	return int(m.buttonID)
}

// ButtonPressed returns true if the button identified by ButtonID was
// pressed during the event. It returns false if instead it was released.
// If ButtonID is 0 (no button for this mouse event), then ButtonPressed
// returns true as this is how the xterm X11 mouse tracking reports it.
func (m MouseEvent) ButtonPressed() bool {
	return m.pressed
}

// Coords returns the screen coordinates of the mouse for this event.
// The upper left character position on the terminal is denoted as 1,1.
func (m MouseEvent) Coords() (x, y int) {
	return int(m.x), int(m.y)
}

// KeyType represents the type of key.
type KeyType byte

// String returns the string representation of the key type.
func (k KeyType) String() string {
	if int(k) < len(keyNames) && keyNames[k] != "" {
		return keyNames[k]
	}
	return strconv.Itoa(int(k))
}

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
	KeyRune // covers ASCII 32-126 + any other unicode code point - from this point on the key type does not match ASCII values
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
	KeyPrint
	KeyESCSeq
	KeyMouse
	KeyFocusIn
	KeyFocusOut // 116

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

var keyNames = [...]string{
	KeyNUL:      "NUL",
	KeySOH:      "SOH",
	KeySTX:      "STX",
	KeyETX:      "ETX",
	KeyEOT:      "EOT",
	KeyENQ:      "ENQ",
	KeyACK:      "ACK",
	KeyBEL:      "BEL",
	KeyBS:       "BS",
	KeyTAB:      "TAB",
	KeyLF:       "LF",
	KeyVT:       "VT",
	KeyFF:       "FF",
	KeyCR:       "CR",
	KeySO:       "SO",
	KeySI:       "SI",
	KeyDLE:      "DLE",
	KeyDC1:      "DC1",
	KeyDC2:      "DC2",
	KeyDC3:      "DC3",
	KeyDC4:      "DC4",
	KeyNAK:      "NAK",
	KeySYN:      "SYN",
	KeyETB:      "ETB",
	KeyCAN:      "CAN",
	KeyEM:       "EM",
	KeySUB:      "SUB",
	KeyESC:      "ESC",
	KeyFS:       "FS",
	KeyGS:       "GS",
	KeyRS:       "RS",
	KeyUS:       "US",
	KeyLeft:     "Left",
	KeyRight:    "Right",
	KeyUp:       "Up",
	KeyDown:     "Down",
	KeyInsert:   "Insert",
	KeyBacktab:  "Backtab",
	KeyDelete:   "Delete",
	KeyHome:     "Home",
	KeyEnd:      "End",
	KeyPgUp:     "PgUp",
	KeyPgDn:     "PgDn",
	KeyF1:       "F1",
	KeyF2:       "F2",
	KeyF3:       "F3",
	KeyF4:       "F4",
	KeyF5:       "F5",
	KeyF6:       "F6",
	KeyF7:       "F7",
	KeyF8:       "F8",
	KeyF9:       "F9",
	KeyF10:      "F10",
	KeyF11:      "F11",
	KeyF12:      "F12",
	KeyF13:      "F13",
	KeyF14:      "F14",
	KeyF15:      "F15",
	KeyF16:      "F16",
	KeyF17:      "F17",
	KeyF18:      "F18",
	KeyF19:      "F19",
	KeyF20:      "F20",
	KeyF21:      "F21",
	KeyF22:      "F22",
	KeyF23:      "F23",
	KeyF24:      "F24",
	KeyF25:      "F25",
	KeyF26:      "F26",
	KeyF27:      "F27",
	KeyF28:      "F28",
	KeyF29:      "F29",
	KeyF30:      "F30",
	KeyF31:      "F31",
	KeyF32:      "F32",
	KeyF33:      "F33",
	KeyF34:      "F34",
	KeyF35:      "F35",
	KeyF36:      "F36",
	KeyF37:      "F37",
	KeyF38:      "F38",
	KeyF39:      "F39",
	KeyF40:      "F40",
	KeyF41:      "F41",
	KeyF42:      "F42",
	KeyF43:      "F43",
	KeyF44:      "F44",
	KeyF45:      "F45",
	KeyF46:      "F46",
	KeyF47:      "F47",
	KeyF48:      "F48",
	KeyF49:      "F49",
	KeyF50:      "F50",
	KeyF51:      "F51",
	KeyF52:      "F52",
	KeyF53:      "F53",
	KeyF54:      "F54",
	KeyF55:      "F55",
	KeyF56:      "F56",
	KeyF57:      "F57",
	KeyF58:      "F58",
	KeyF59:      "F59",
	KeyF60:      "F60",
	KeyF61:      "F61",
	KeyF62:      "F62",
	KeyF63:      "F63",
	KeyF64:      "F64",
	KeyHelp:     "Help",
	KeyExit:     "Exit",
	KeyClear:    "Clear",
	KeyCancel:   "Cancel",
	KeyPrint:    "Print",
	KeyESCSeq:   "ESCSeq",
	KeyMouse:    "Mouse",
	KeyFocusIn:  "FocusIn",
	KeyFocusOut: "FocusOut",
	KeyDEL:      "DEL",
}
