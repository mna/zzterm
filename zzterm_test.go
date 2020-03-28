package zzterm

import (
	"strings"
	"testing"
)

func TestInput_ReadKey(t *testing.T) {
	cases := []struct {
		in  string
		r   rune
		typ KeyType
		m   Mod
	}{
		{"a", 'a', KeyRune, ModNone},
		{"B", 'B', KeyRune, ModNone},
		{"1", '1', KeyRune, ModNone},
		{"\x00", -1, KeyNUL, ModNone},
		{"\x01", -1, KeySOH, ModNone},
		{"\n", -1, KeyLF, ModNone},
		{"•", '•', KeyRune, ModNone},
		{"👪", '👪', KeyRune, ModNone},
		{"🤡", '🤡', KeyRune, ModNone},
		{"𐰧", '𐰧', KeyRune, ModNone},
		{"⺜", '⺜', KeyRune, ModNone},
		{"\u007f", -1, KeyDEL, ModNone},
		{"\x1b[A", -1, KeyUp, ModNone},
		{"\x1b[3~", -1, KeyDelete, ModNone},
		{"\x1b[1;2D", -1, KeyLeft, ModShift},
		{"\x1b[1;2C", -1, KeyRight, ModShift},
		{"\x1bb", -1, KeyLeft, ModAlt},
		{"\x1bf", -1, KeyRight, ModAlt},
	}

	input := NewInput()
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			k, err := input.ReadKey(strings.NewReader(c.in))
			if err != nil {
				t.Fatal(err)
			}

			if k.Rune() != c.r {
				t.Errorf("want rune %c (%[1]U), got %c (%[2]U)", c.r, k.Rune())
			}
			if k.Type() != c.typ {
				t.Errorf("want key type %d, got %d", c.typ, k.Type())
			}
			if k.Mod() != c.m {
				t.Errorf("want modifier flags %04b, got %04b", c.m, k.Mod())
			}
		})
	}
}

var BenchmarkResult Key

func BenchmarkInput_ReadKey(b *testing.B) {
	cases := []string{
		"a", "B", "1", "\x00", "ø", "👪", "平",
		"\x1b[B", "\x1b[1;2C",
	}
	for _, c := range cases {
		input := NewInput()
		b.Run(c, func(b *testing.B) {
			r := strings.NewReader(c)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				k, err := input.ReadKey(r)
				if err != nil {
					b.Fatal(err)
				}
				BenchmarkResult = k
				r.Reset(c)
			}
		})
	}
}
