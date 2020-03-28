package zzterm

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

type testcase struct {
	in  string
	r   rune
	typ KeyType
	m   Mod
}

func TestInput_ReadKey_DefaultTinfo(t *testing.T) {
	cases := []testcase{
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
	}

	input := NewInput(nil)
	for _, c := range cases {
		runTestcase(t, c, input)
	}
}

func TestInput_ReadKey_VT100Tinfo(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/vt100.json")
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}

	cases := []testcase{
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
		{"\x1bOA", -1, KeyUp, ModNone},
		{"\x1bOP", -1, KeyF1, ModNone},
		{"\x1bOD", -1, KeyLeft, ModNone},
	}

	input := NewInput(FromTerminfo(m))
	for _, c := range cases {
		runTestcase(t, c, input)
	}
}

func runTestcase(t *testing.T, c testcase, input *Input) {
	t.Helper()

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

var BenchmarkResult Key

func BenchmarkInput_ReadKey(b *testing.B) {
	cases := []string{
		"a", "B", "1", "\x00", "ø", "👪", "平",
		"\x1b[B", "\x1b[1;2C",
	}
	for _, c := range cases {
		input := NewInput(nil)
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
