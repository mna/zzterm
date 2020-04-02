package zzterm

import (
	"encoding/json"
	"strings"
)

// FromTerminfo returns a terminfo map that can be used in the call to
// NewInput. The value v should be a tcell/terminfo.Terminfo struct, a
// pointer to such a struct, or a value that marshals to JSON with an
// equivalent structure.
//
// It first marshals v to JSON and then unmarshals it in a map.  It makes no
// validation that v is a valid terminfo, and it returns nil if there is any
// error when converting to and from the intermediate JSON representations.
func FromTerminfo(v interface{}) map[string]string {
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}

var defaultEsc = map[string]Key{
	"\x1b[A":     keyFromTypeMod(KeyUp, ModNone),
	"\x1b[B":     keyFromTypeMod(KeyDown, ModNone),
	"\x1b[C":     keyFromTypeMod(KeyRight, ModNone),
	"\x1b[D":     keyFromTypeMod(KeyLeft, ModNone),
	"\x1b[2~":    keyFromTypeMod(KeyInsert, ModNone),
	"\x1b[3~":    keyFromTypeMod(KeyDelete, ModNone),
	"\x1b[Z":     keyFromTypeMod(KeyBacktab, ModNone),
	"\x1bOH":     keyFromTypeMod(KeyHome, ModNone),
	"\x1bOF":     keyFromTypeMod(KeyEnd, ModNone),
	"\x1b[5~":    keyFromTypeMod(KeyPgUp, ModNone),
	"\x1b[6~":    keyFromTypeMod(KeyPgDn, ModNone),
	"\x1bOP":     keyFromTypeMod(KeyF1, ModNone),
	"\x1bOQ":     keyFromTypeMod(KeyF2, ModNone),
	"\x1bOR":     keyFromTypeMod(KeyF3, ModNone),
	"\x1bOS":     keyFromTypeMod(KeyF4, ModNone),
	"\x1b[15~":   keyFromTypeMod(KeyF5, ModNone),
	"\x1b[17~":   keyFromTypeMod(KeyF6, ModNone),
	"\x1b[18~":   keyFromTypeMod(KeyF7, ModNone),
	"\x1b[19~":   keyFromTypeMod(KeyF8, ModNone),
	"\x1b[20~":   keyFromTypeMod(KeyF9, ModNone),
	"\x1b[21~":   keyFromTypeMod(KeyF10, ModNone),
	"\x1b[23~":   keyFromTypeMod(KeyF11, ModNone),
	"\x1b[24~":   keyFromTypeMod(KeyF12, ModNone),
	"\x1b[1;2P":  keyFromTypeMod(KeyF13, ModNone),
	"\x1b[1;2Q":  keyFromTypeMod(KeyF14, ModNone),
	"\x1b[1;2R":  keyFromTypeMod(KeyF15, ModNone),
	"\x1b[1;2S":  keyFromTypeMod(KeyF16, ModNone),
	"\x1b[15;2~": keyFromTypeMod(KeyF17, ModNone),
	"\x1b[17;2~": keyFromTypeMod(KeyF18, ModNone),
	"\x1b[18;2~": keyFromTypeMod(KeyF19, ModNone),
	"\x1b[19;2~": keyFromTypeMod(KeyF20, ModNone),
	"\x1b[1;2D":  keyFromTypeMod(KeyLeft, ModShift),
	"\x1b[1;2C":  keyFromTypeMod(KeyRight, ModShift),
}

func cloneEscMap(m map[string]Key) map[string]Key {
	mm := make(map[string]Key)
	for k, v := range m {
		mm[k] = v
	}
	return mm
}

func addFocusESCSeq(m map[string]Key) {
	m["\x1b[I"] = keyFromTypeMod(KeyFocusIn, ModNone)
	m["\x1b[O"] = keyFromTypeMod(KeyFocusOut, ModNone)
}

func escFromTerminfo(tinfo map[string]string) map[string]Key {
	if tinfo == nil {
		return cloneEscMap(defaultEsc)
	}

	m := make(map[string]Key)
	for k, v := range tinfo {
		if !strings.HasPrefix(k, "Key") || !strings.HasPrefix(v, "\x1b") {
			continue
		}
		switch k {
		case "KeyBackspace":
			m[v] = keyFromTypeMod(KeyBS, ModNone)
		case "KeyF1":
			m[v] = keyFromTypeMod(KeyF1, ModNone)
		case "KeyF2":
			m[v] = keyFromTypeMod(KeyF2, ModNone)
		case "KeyF3":
			m[v] = keyFromTypeMod(KeyF3, ModNone)
		case "KeyF4":
			m[v] = keyFromTypeMod(KeyF4, ModNone)
		case "KeyF5":
			m[v] = keyFromTypeMod(KeyF5, ModNone)
		case "KeyF6":
			m[v] = keyFromTypeMod(KeyF6, ModNone)
		case "KeyF7":
			m[v] = keyFromTypeMod(KeyF7, ModNone)
		case "KeyF8":
			m[v] = keyFromTypeMod(KeyF8, ModNone)
		case "KeyF9":
			m[v] = keyFromTypeMod(KeyF9, ModNone)
		case "KeyF10":
			m[v] = keyFromTypeMod(KeyF10, ModNone)
		case "KeyF11":
			m[v] = keyFromTypeMod(KeyF11, ModNone)
		case "KeyF12":
			m[v] = keyFromTypeMod(KeyF12, ModNone)
		case "KeyF13":
			m[v] = keyFromTypeMod(KeyF13, ModNone)
		case "KeyF14":
			m[v] = keyFromTypeMod(KeyF14, ModNone)
		case "KeyF15":
			m[v] = keyFromTypeMod(KeyF15, ModNone)
		case "KeyF16":
			m[v] = keyFromTypeMod(KeyF16, ModNone)
		case "KeyF17":
			m[v] = keyFromTypeMod(KeyF17, ModNone)
		case "KeyF18":
			m[v] = keyFromTypeMod(KeyF18, ModNone)
		case "KeyF19":
			m[v] = keyFromTypeMod(KeyF19, ModNone)
		case "KeyF20":
			m[v] = keyFromTypeMod(KeyF20, ModNone)
		case "KeyF21":
			m[v] = keyFromTypeMod(KeyF21, ModNone)
		case "KeyF22":
			m[v] = keyFromTypeMod(KeyF22, ModNone)
		case "KeyF23":
			m[v] = keyFromTypeMod(KeyF23, ModNone)
		case "KeyF24":
			m[v] = keyFromTypeMod(KeyF24, ModNone)
		case "KeyF25":
			m[v] = keyFromTypeMod(KeyF25, ModNone)
		case "KeyF26":
			m[v] = keyFromTypeMod(KeyF26, ModNone)
		case "KeyF27":
			m[v] = keyFromTypeMod(KeyF27, ModNone)
		case "KeyF28":
			m[v] = keyFromTypeMod(KeyF28, ModNone)
		case "KeyF29":
			m[v] = keyFromTypeMod(KeyF29, ModNone)
		case "KeyF30":
			m[v] = keyFromTypeMod(KeyF30, ModNone)
		case "KeyF31":
			m[v] = keyFromTypeMod(KeyF31, ModNone)
		case "KeyF32":
			m[v] = keyFromTypeMod(KeyF32, ModNone)
		case "KeyF33":
			m[v] = keyFromTypeMod(KeyF33, ModNone)
		case "KeyF34":
			m[v] = keyFromTypeMod(KeyF34, ModNone)
		case "KeyF35":
			m[v] = keyFromTypeMod(KeyF35, ModNone)
		case "KeyF36":
			m[v] = keyFromTypeMod(KeyF36, ModNone)
		case "KeyF37":
			m[v] = keyFromTypeMod(KeyF37, ModNone)
		case "KeyF38":
			m[v] = keyFromTypeMod(KeyF38, ModNone)
		case "KeyF39":
			m[v] = keyFromTypeMod(KeyF39, ModNone)
		case "KeyF40":
			m[v] = keyFromTypeMod(KeyF40, ModNone)
		case "KeyF41":
			m[v] = keyFromTypeMod(KeyF41, ModNone)
		case "KeyF42":
			m[v] = keyFromTypeMod(KeyF42, ModNone)
		case "KeyF43":
			m[v] = keyFromTypeMod(KeyF43, ModNone)
		case "KeyF44":
			m[v] = keyFromTypeMod(KeyF44, ModNone)
		case "KeyF45":
			m[v] = keyFromTypeMod(KeyF45, ModNone)
		case "KeyF46":
			m[v] = keyFromTypeMod(KeyF46, ModNone)
		case "KeyF47":
			m[v] = keyFromTypeMod(KeyF47, ModNone)
		case "KeyF48":
			m[v] = keyFromTypeMod(KeyF48, ModNone)
		case "KeyF49":
			m[v] = keyFromTypeMod(KeyF49, ModNone)
		case "KeyF50":
			m[v] = keyFromTypeMod(KeyF50, ModNone)
		case "KeyF51":
			m[v] = keyFromTypeMod(KeyF51, ModNone)
		case "KeyF52":
			m[v] = keyFromTypeMod(KeyF52, ModNone)
		case "KeyF53":
			m[v] = keyFromTypeMod(KeyF53, ModNone)
		case "KeyF54":
			m[v] = keyFromTypeMod(KeyF54, ModNone)
		case "KeyF55":
			m[v] = keyFromTypeMod(KeyF55, ModNone)
		case "KeyF56":
			m[v] = keyFromTypeMod(KeyF56, ModNone)
		case "KeyF57":
			m[v] = keyFromTypeMod(KeyF57, ModNone)
		case "KeyF58":
			m[v] = keyFromTypeMod(KeyF58, ModNone)
		case "KeyF59":
			m[v] = keyFromTypeMod(KeyF59, ModNone)
		case "KeyF60":
			m[v] = keyFromTypeMod(KeyF60, ModNone)
		case "KeyF61":
			m[v] = keyFromTypeMod(KeyF61, ModNone)
		case "KeyF62":
			m[v] = keyFromTypeMod(KeyF62, ModNone)
		case "KeyF63":
			m[v] = keyFromTypeMod(KeyF63, ModNone)
		case "KeyF64":
			m[v] = keyFromTypeMod(KeyF64, ModNone)
		case "KeyInsert":
			m[v] = keyFromTypeMod(KeyInsert, ModNone)
		case "KeyDelete":
			m[v] = keyFromTypeMod(KeyDelete, ModNone)
		case "KeyHome":
			m[v] = keyFromTypeMod(KeyHome, ModNone)
		case "KeyEnd":
			m[v] = keyFromTypeMod(KeyEnd, ModNone)
		case "KeyHelp":
			m[v] = keyFromTypeMod(KeyHelp, ModNone)
		case "KeyPgUp":
			m[v] = keyFromTypeMod(KeyPgUp, ModNone)
		case "KeyPgDn":
			m[v] = keyFromTypeMod(KeyPgDn, ModNone)
		case "KeyUp":
			m[v] = keyFromTypeMod(KeyUp, ModNone)
		case "KeyDown":
			m[v] = keyFromTypeMod(KeyDown, ModNone)
		case "KeyLeft":
			m[v] = keyFromTypeMod(KeyLeft, ModNone)
		case "KeyRight":
			m[v] = keyFromTypeMod(KeyRight, ModNone)
		case "KeyBacktab":
			m[v] = keyFromTypeMod(KeyBacktab, ModNone)
		case "KeyExit":
			m[v] = keyFromTypeMod(KeyExit, ModNone)
		case "KeyClear":
			m[v] = keyFromTypeMod(KeyClear, ModNone)
		case "KeyPrint":
			m[v] = keyFromTypeMod(KeyPrint, ModNone)
		case "KeyCancel":
			m[v] = keyFromTypeMod(KeyCancel, ModNone)
		case "KeyShfRight":
			m[v] = keyFromTypeMod(KeyRight, ModShift)
		case "KeyShfLeft":
			m[v] = keyFromTypeMod(KeyLeft, ModShift)
		case "KeyShfHome":
			m[v] = keyFromTypeMod(KeyHome, ModShift)
		case "KeyShfEnd":
			m[v] = keyFromTypeMod(KeyEnd, ModShift)
		case "KeyShfUp":
			m[v] = keyFromTypeMod(KeyUp, ModShift)
		case "KeyShfDown":
			m[v] = keyFromTypeMod(KeyDown, ModShift)
		case "KeyShfPgUp":
			m[v] = keyFromTypeMod(KeyPgUp, ModShift)
		case "KeyShfPgDn":
			m[v] = keyFromTypeMod(KeyPgDn, ModShift)
		case "KeyCtrlUp":
			m[v] = keyFromTypeMod(KeyUp, ModCtrl)
		case "KeyCtrlDown":
			m[v] = keyFromTypeMod(KeyDown, ModCtrl)
		case "KeyCtrlRight":
			m[v] = keyFromTypeMod(KeyRight, ModCtrl)
		case "KeyCtrlLeft":
			m[v] = keyFromTypeMod(KeyLeft, ModCtrl)
		case "KeyMetaUp":
			m[v] = keyFromTypeMod(KeyUp, ModMeta)
		case "KeyMetaDown":
			m[v] = keyFromTypeMod(KeyDown, ModMeta)
		case "KeyMetaRight":
			m[v] = keyFromTypeMod(KeyRight, ModMeta)
		case "KeyMetaLeft":
			m[v] = keyFromTypeMod(KeyLeft, ModMeta)
		case "KeyAltUp":
			m[v] = keyFromTypeMod(KeyUp, ModAlt)
		case "KeyAltDown":
			m[v] = keyFromTypeMod(KeyDown, ModAlt)
		case "KeyAltRight":
			m[v] = keyFromTypeMod(KeyRight, ModAlt)
		case "KeyAltLeft":
			m[v] = keyFromTypeMod(KeyLeft, ModAlt)
		case "KeyCtrlHome":
			m[v] = keyFromTypeMod(KeyHome, ModCtrl)
		case "KeyCtrlEnd":
			m[v] = keyFromTypeMod(KeyEnd, ModCtrl)
		case "KeyMetaHome":
			m[v] = keyFromTypeMod(KeyHome, ModMeta)
		case "KeyMetaEnd":
			m[v] = keyFromTypeMod(KeyEnd, ModMeta)
		case "KeyAltHome":
			m[v] = keyFromTypeMod(KeyHome, ModAlt)
		case "KeyAltEnd":
			m[v] = keyFromTypeMod(KeyEnd, ModAlt)
		case "KeyAltShfUp":
			m[v] = keyFromTypeMod(KeyUp, ModAlt|ModShift)
		case "KeyAltShfDown":
			m[v] = keyFromTypeMod(KeyDown, ModAlt|ModShift)
		case "KeyAltShfLeft":
			m[v] = keyFromTypeMod(KeyLeft, ModAlt|ModShift)
		case "KeyAltShfRight":
			m[v] = keyFromTypeMod(KeyRight, ModAlt|ModShift)
		case "KeyMetaShfUp":
			m[v] = keyFromTypeMod(KeyUp, ModMeta|ModShift)
		case "KeyMetaShfDown":
			m[v] = keyFromTypeMod(KeyDown, ModMeta|ModShift)
		case "KeyMetaShfLeft":
			m[v] = keyFromTypeMod(KeyLeft, ModMeta|ModShift)
		case "KeyMetaShfRight":
			m[v] = keyFromTypeMod(KeyRight, ModMeta|ModShift)
		case "KeyCtrlShfUp":
			m[v] = keyFromTypeMod(KeyUp, ModCtrl|ModShift)
		case "KeyCtrlShfDown":
			m[v] = keyFromTypeMod(KeyDown, ModCtrl|ModShift)
		case "KeyCtrlShfLeft":
			m[v] = keyFromTypeMod(KeyLeft, ModCtrl|ModShift)
		case "KeyCtrlShfRight":
			m[v] = keyFromTypeMod(KeyRight, ModCtrl|ModShift)
		case "KeyCtrlShfHome":
			m[v] = keyFromTypeMod(KeyHome, ModCtrl|ModShift)
		case "KeyCtrlShfEnd":
			m[v] = keyFromTypeMod(KeyEnd, ModCtrl|ModShift)
		case "KeyAltShfHome":
			m[v] = keyFromTypeMod(KeyHome, ModAlt|ModShift)
		case "KeyAltShfEnd":
			m[v] = keyFromTypeMod(KeyEnd, ModAlt|ModShift)
		case "KeyMetaShfHome":
			m[v] = keyFromTypeMod(KeyHome, ModMeta|ModShift)
		case "KeyMetaShfEnd":
			m[v] = keyFromTypeMod(KeyEnd, ModMeta|ModShift)
		}
	}
	return m
}
