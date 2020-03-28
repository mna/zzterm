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
	"\x1bOQ":     keyFromTypeMod(KeyF1, ModNone),
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

func escFromTerminfo(tinfo map[string]string) map[string]Key {
	if tinfo == nil {
		return defaultEsc
	}

	m := make(map[string]Key)
	for k, v := range tinfo {
		if !strings.HasPrefix(k, "Key") || !strings.HasPrefix(v, "\x1b") {
			continue
		}
		switch k {
		case "KeyBackspace":
			m[k] = keyFromTypeMod(KeyBS, ModNone)
		case "KeyF1":
			m[k] = keyFromTypeMod(KeyF1, ModNone)
		case "KeyF2":
			m[k] = keyFromTypeMod(KeyF2, ModNone)
		case "KeyF3":
			m[k] = keyFromTypeMod(KeyF3, ModNone)
		case "KeyF4":
			m[k] = keyFromTypeMod(KeyF4, ModNone)
		case "KeyF5":
			m[k] = keyFromTypeMod(KeyF5, ModNone)
		case "KeyF6":
			m[k] = keyFromTypeMod(KeyF6, ModNone)
		case "KeyF7":
			m[k] = keyFromTypeMod(KeyF7, ModNone)
		case "KeyF8":
			m[k] = keyFromTypeMod(KeyF8, ModNone)
		case "KeyF9":
			m[k] = keyFromTypeMod(KeyF9, ModNone)
		case "KeyF10":
			m[k] = keyFromTypeMod(KeyF10, ModNone)
		case "KeyF11":
			m[k] = keyFromTypeMod(KeyF11, ModNone)
		case "KeyF12":
			m[k] = keyFromTypeMod(KeyF12, ModNone)
		case "KeyF13":
			m[k] = keyFromTypeMod(KeyF13, ModNone)
		case "KeyF14":
			m[k] = keyFromTypeMod(KeyF14, ModNone)
		case "KeyF15":
			m[k] = keyFromTypeMod(KeyF15, ModNone)
		case "KeyF16":
			m[k] = keyFromTypeMod(KeyF16, ModNone)
		case "KeyF17":
			m[k] = keyFromTypeMod(KeyF17, ModNone)
		case "KeyF18":
			m[k] = keyFromTypeMod(KeyF18, ModNone)
		case "KeyF19":
			m[k] = keyFromTypeMod(KeyF19, ModNone)
		case "KeyF20":
			m[k] = keyFromTypeMod(KeyF20, ModNone)
		case "KeyF21":
			m[k] = keyFromTypeMod(KeyF21, ModNone)
		case "KeyF22":
			m[k] = keyFromTypeMod(KeyF22, ModNone)
		case "KeyF23":
			m[k] = keyFromTypeMod(KeyF23, ModNone)
		case "KeyF24":
			m[k] = keyFromTypeMod(KeyF24, ModNone)
		case "KeyF25":
			m[k] = keyFromTypeMod(KeyF25, ModNone)
		case "KeyF26":
			m[k] = keyFromTypeMod(KeyF26, ModNone)
		case "KeyF27":
			m[k] = keyFromTypeMod(KeyF27, ModNone)
		case "KeyF28":
			m[k] = keyFromTypeMod(KeyF28, ModNone)
		case "KeyF29":
			m[k] = keyFromTypeMod(KeyF29, ModNone)
		case "KeyF30":
			m[k] = keyFromTypeMod(KeyF30, ModNone)
		case "KeyF31":
			m[k] = keyFromTypeMod(KeyF31, ModNone)
		case "KeyF32":
			m[k] = keyFromTypeMod(KeyF32, ModNone)
		case "KeyF33":
			m[k] = keyFromTypeMod(KeyF33, ModNone)
		case "KeyF34":
			m[k] = keyFromTypeMod(KeyF34, ModNone)
		case "KeyF35":
			m[k] = keyFromTypeMod(KeyF35, ModNone)
		case "KeyF36":
			m[k] = keyFromTypeMod(KeyF36, ModNone)
		case "KeyF37":
			m[k] = keyFromTypeMod(KeyF37, ModNone)
		case "KeyF38":
			m[k] = keyFromTypeMod(KeyF38, ModNone)
		case "KeyF39":
			m[k] = keyFromTypeMod(KeyF39, ModNone)
		case "KeyF40":
			m[k] = keyFromTypeMod(KeyF40, ModNone)
		case "KeyF41":
			m[k] = keyFromTypeMod(KeyF41, ModNone)
		case "KeyF42":
			m[k] = keyFromTypeMod(KeyF42, ModNone)
		case "KeyF43":
			m[k] = keyFromTypeMod(KeyF43, ModNone)
		case "KeyF44":
			m[k] = keyFromTypeMod(KeyF44, ModNone)
		case "KeyF45":
			m[k] = keyFromTypeMod(KeyF45, ModNone)
		case "KeyF46":
			m[k] = keyFromTypeMod(KeyF46, ModNone)
		case "KeyF47":
			m[k] = keyFromTypeMod(KeyF47, ModNone)
		case "KeyF48":
			m[k] = keyFromTypeMod(KeyF48, ModNone)
		case "KeyF49":
			m[k] = keyFromTypeMod(KeyF49, ModNone)
		case "KeyF50":
			m[k] = keyFromTypeMod(KeyF50, ModNone)
		case "KeyF51":
			m[k] = keyFromTypeMod(KeyF51, ModNone)
		case "KeyF52":
			m[k] = keyFromTypeMod(KeyF52, ModNone)
		case "KeyF53":
			m[k] = keyFromTypeMod(KeyF53, ModNone)
		case "KeyF54":
			m[k] = keyFromTypeMod(KeyF54, ModNone)
		case "KeyF55":
			m[k] = keyFromTypeMod(KeyF55, ModNone)
		case "KeyF56":
			m[k] = keyFromTypeMod(KeyF56, ModNone)
		case "KeyF57":
			m[k] = keyFromTypeMod(KeyF57, ModNone)
		case "KeyF58":
			m[k] = keyFromTypeMod(KeyF58, ModNone)
		case "KeyF59":
			m[k] = keyFromTypeMod(KeyF59, ModNone)
		case "KeyF60":
			m[k] = keyFromTypeMod(KeyF60, ModNone)
		case "KeyF61":
			m[k] = keyFromTypeMod(KeyF61, ModNone)
		case "KeyF62":
			m[k] = keyFromTypeMod(KeyF62, ModNone)
		case "KeyF63":
			m[k] = keyFromTypeMod(KeyF63, ModNone)
		case "KeyF64":
			m[k] = keyFromTypeMod(KeyF64, ModNone)
		case "KeyInsert":
			m[k] = keyFromTypeMod(KeyInsert, ModNone)
		case "KeyDelete":
			m[k] = keyFromTypeMod(KeyDelete, ModNone)
		case "KeyHome":
			m[k] = keyFromTypeMod(KeyHome, ModNone)
		case "KeyEnd":
			m[k] = keyFromTypeMod(KeyEnd, ModNone)
		case "KeyHelp":
			m[k] = keyFromTypeMod(KeyHelp, ModNone)
		case "KeyPgUp":
			m[k] = keyFromTypeMod(KeyPgUp, ModNone)
		case "KeyPgDn":
			m[k] = keyFromTypeMod(KeyPgDn, ModNone)
		case "KeyUp":
			m[k] = keyFromTypeMod(KeyUp, ModNone)
		case "KeyDown":
			m[k] = keyFromTypeMod(KeyDown, ModNone)
		case "KeyLeft":
			m[k] = keyFromTypeMod(KeyLeft, ModNone)
		case "KeyRight":
			m[k] = keyFromTypeMod(KeyRight, ModNone)
		case "KeyBacktab":
			m[k] = keyFromTypeMod(KeyBacktab, ModNone)
		case "KeyExit":
			m[k] = keyFromTypeMod(KeyExit, ModNone)
		case "KeyClear":
			m[k] = keyFromTypeMod(KeyClear, ModNone)
		case "KeyPrint":
			m[k] = keyFromTypeMod(KeyPrint, ModNone)
		case "KeyCancel":
			m[k] = keyFromTypeMod(KeyCancel, ModNone)
		case "KeyShfRight":
			m[k] = keyFromTypeMod(KeyRight, ModShift)
		case "KeyShfLeft":
			m[k] = keyFromTypeMod(KeyLeft, ModShift)
		case "KeyShfHome":
			m[k] = keyFromTypeMod(KeyHome, ModShift)
		case "KeyShfEnd":
			m[k] = keyFromTypeMod(KeyEnd, ModShift)
		case "KeyShfUp":
			m[k] = keyFromTypeMod(KeyUp, ModShift)
		case "KeyShfDown":
			m[k] = keyFromTypeMod(KeyDown, ModShift)
		case "KeyShfPgUp":
			m[k] = keyFromTypeMod(KeyPgUp, ModShift)
		case "KeyShfPgDn":
			m[k] = keyFromTypeMod(KeyPgDn, ModShift)
		case "KeyCtrlUp":
			m[k] = keyFromTypeMod(KeyUp, ModCtrl)
		case "KeyCtrlDown":
			m[k] = keyFromTypeMod(KeyDown, ModCtrl)
		case "KeyCtrlRight":
			m[k] = keyFromTypeMod(KeyRight, ModCtrl)
		case "KeyCtrlLeft":
			m[k] = keyFromTypeMod(KeyLeft, ModCtrl)
		case "KeyMetaUp":
			m[k] = keyFromTypeMod(KeyUp, ModMeta)
		case "KeyMetaDown":
			m[k] = keyFromTypeMod(KeyDown, ModMeta)
		case "KeyMetaRight":
			m[k] = keyFromTypeMod(KeyRight, ModMeta)
		case "KeyMetaLeft":
			m[k] = keyFromTypeMod(KeyLeft, ModMeta)
		case "KeyAltUp":
			m[k] = keyFromTypeMod(KeyUp, ModAlt)
		case "KeyAltDown":
			m[k] = keyFromTypeMod(KeyDown, ModAlt)
		case "KeyAltRight":
			m[k] = keyFromTypeMod(KeyRight, ModAlt)
		case "KeyAltLeft":
			m[k] = keyFromTypeMod(KeyLeft, ModAlt)
		case "KeyCtrlHome":
			m[k] = keyFromTypeMod(KeyHome, ModCtrl)
		case "KeyCtrlEnd":
			m[k] = keyFromTypeMod(KeyEnd, ModCtrl)
		case "KeyMetaHome":
			m[k] = keyFromTypeMod(KeyHome, ModMeta)
		case "KeyMetaEnd":
			m[k] = keyFromTypeMod(KeyEnd, ModMeta)
		case "KeyAltHome":
			m[k] = keyFromTypeMod(KeyHome, ModAlt)
		case "KeyAltEnd":
			m[k] = keyFromTypeMod(KeyEnd, ModAlt)
		case "KeyAltShfUp":
			m[k] = keyFromTypeMod(KeyUp, ModAlt|ModShift)
		case "KeyAltShfDown":
			m[k] = keyFromTypeMod(KeyDown, ModAlt|ModShift)
		case "KeyAltShfLeft":
			m[k] = keyFromTypeMod(KeyLeft, ModAlt|ModShift)
		case "KeyAltShfRight":
			m[k] = keyFromTypeMod(KeyRight, ModAlt|ModShift)
		case "KeyMetaShfUp":
			m[k] = keyFromTypeMod(KeyUp, ModMeta|ModShift)
		case "KeyMetaShfDown":
			m[k] = keyFromTypeMod(KeyDown, ModMeta|ModShift)
		case "KeyMetaShfLeft":
			m[k] = keyFromTypeMod(KeyLeft, ModMeta|ModShift)
		case "KeyMetaShfRight":
			m[k] = keyFromTypeMod(KeyRight, ModMeta|ModShift)
		case "KeyCtrlShfUp":
			m[k] = keyFromTypeMod(KeyUp, ModCtrl|ModShift)
		case "KeyCtrlShfDown":
			m[k] = keyFromTypeMod(KeyDown, ModCtrl|ModShift)
		case "KeyCtrlShfLeft":
			m[k] = keyFromTypeMod(KeyLeft, ModCtrl|ModShift)
		case "KeyCtrlShfRight":
			m[k] = keyFromTypeMod(KeyRight, ModCtrl|ModShift)
		case "KeyCtrlShfHome":
			m[k] = keyFromTypeMod(KeyHome, ModCtrl|ModShift)
		case "KeyCtrlShfEnd":
			m[k] = keyFromTypeMod(KeyEnd, ModCtrl|ModShift)
		case "KeyAltShfHome":
			m[k] = keyFromTypeMod(KeyHome, ModAlt|ModShift)
		case "KeyAltShfEnd":
			m[k] = keyFromTypeMod(KeyEnd, ModAlt|ModShift)
		case "KeyMetaShfHome":
			m[k] = keyFromTypeMod(KeyHome, ModMeta|ModShift)
		case "KeyMetaShfEnd":
			m[k] = keyFromTypeMod(KeyEnd, ModMeta|ModShift)
		}
	}
	return m
}
