package zzterm

import "testing"

func TestKey_String(t *testing.T) {
	cases := []struct {
		key Key
		out string
	}{
		{keyFromTypeMod(KeyNUL, ModNone), `Key(NUL)`},
		{keyFromTypeMod(KeyESC, ModNone), `Key(ESC)`},
		{keyFromTypeMod(KeyUp, ModShift), `Key(⇧ Up)`},
		{Key('a'), `Key(U+0061 'a')`},
		{Key('👪'), `Key(U+1F46A '👪')`},
		{Key('\u202f'), `Key(U+202F)`},
		{keyFromTypeMod(KeyHome, ModCtrl|ModShift), `Key(⌃⇧ Home)`},
		{keyFromTypeMod(KeyLeft, ModAlt), `Key(⎇ Left)`},
		{keyFromTypeMod(KeyLeft, ModMeta), `Key(⌥ Left)`},
	}
	for _, c := range cases {
		t.Run(c.key.String(), func(t *testing.T) {
			s := c.key.String()
			if s != c.out {
				t.Errorf("want %s, got %s", c.out, s)
			}
		})
	}
}
