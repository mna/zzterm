package zzterm

import "testing"

func TestKey_String(t *testing.T) {
	cases := []struct {
		key Key
		out string
	}{
		{keyFromTypeMod(KeyNUL, ModNone), `Key(x00)`},
		{keyFromTypeMod(KeyESC, ModNone), `Key(x1b)`},
		{keyFromTypeMod(KeyUp, ModShift), `Key(⇧ x23)`},
		{Key('a'), `Key(U+0061 'a')`},
		{Key('👪'), `Key(U+1F46A '👪')`},
		{Key('\u202f'), `Key(U+202F)`},
		{keyFromTypeMod(KeyHome, ModCtrl|ModShift), `Key(⌃⇧ x28)`},
		{keyFromTypeMod(KeyLeft, ModAlt), `Key(⎇ x21)`},
		{keyFromTypeMod(KeyLeft, ModMeta), `Key(⌥ x21)`},
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
