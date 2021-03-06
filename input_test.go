package zzterm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestInput_ReadKey_Multiple(t *testing.T) {
	invalidRuneKey := Key('\x01')

	cases := []struct {
		in    string
		keys  []Key
		bytes []string
	}{
		{"", nil, nil},
		{"a", []Key{Key('a')}, []string{"a"}},
		{"ab", []Key{Key('a'), Key('b')}, []string{"a", "b"}},
		{"\xff", []Key{invalidRuneKey}, []string{"\xff"}},
		{"\xffa", []Key{invalidRuneKey, Key('a')}, []string{"\xff", "a"}},
		{"😿\x1b[abc", []Key{Key('😿'), keyFromTypeMod(KeyESCSeq, ModNone)}, []string{"😿", "\x1b[abc"}},
	}

	input := NewInput(WithMouse(), WithFocus())
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			r := strings.NewReader(c.in)
			for i, wantk := range c.keys {
				wantb := c.bytes[i]
				got, err := input.ReadKey(r)
				if wantk == invalidRuneKey {
					if err.Error() != "invalid rune" {
						t.Fatalf("[%d]: want invalid rune, got %v", i, err)
					}
					wantk = Key(0)
				} else if err != nil {
					t.Fatalf("[%d]: want %s, got error %v", i, wantk, err)
				}

				if got.Type() != wantk.Type() {
					t.Fatalf("[%d]: want key type %s, got %s", i, wantk.Type(), got.Type())
				}
				if gotb := string(input.Bytes()); gotb != wantb {
					t.Fatalf("[%d]: want bytes %q, got %q", i, wantb, gotb)
				}
			}

			// after the loop, must return a "timeout" error (via EOF)
			got, err := input.ReadKey(r)
			if !errors.Is(err, ErrTimeout) {
				t.Fatalf("after loop: want ErrTimeout, got %v (key %v)", err, got)
			}
			if !os.IsTimeout(err) {
				t.Fatal("after loop: want ErrTimeout to be identified as such with os.IsTimeout")
			}
		})
	}
}

func TestInput_ReadKey_BustBuffer(t *testing.T) {
	// this '⬼ ' character is 3 bytes in utf-8, so it ends up crossing
	// over the buffer size (which is an even number). This tests that
	// ReadKey properly tries a Read to get more bytes when it only
	// has an invalid rune to work with.
	want, wantn := "\xE2\xAC\xBC", 100
	r := strings.NewReader(strings.Repeat(want, wantn))
	input := NewInput()
	var count int
	for {
		key, err := input.ReadKey(r)
		if errors.Is(err, ErrTimeout) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		count++
		if got := string(input.Bytes()); got != want {
			t.Fatalf("[%d]: got bytes %q", count, got)
		}
		if key.Type() != KeyRune || key.Rune() != '⬼' {
			t.Fatalf("[%d]: unexpected key %v", count, key)
		}
	}
	if count != 100 {
		t.Fatalf("want %d keys, got %d", wantn, count)
	}
}

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

	input := NewInput()
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

	input := NewInput(WithESCSeq(FromTerminfo(m)))
	for _, c := range cases {
		runTestcase(t, c, input)
	}
}

func TestInput_ReadKey_Focus(t *testing.T) {
	input := NewInput(WithFocus())

	in := "\x1b[I"
	k, err := input.ReadKey(strings.NewReader(in))
	if err != nil {
		log.Fatal(err)
	}
	if k.Type() != KeyFocusIn || k.Mod() != ModNone {
		t.Errorf("invalid modifier flags or key type: %s", k)
	}

	in = "\x1b[O"
	k, err = input.ReadKey(strings.NewReader(in))
	if err != nil {
		log.Fatal(err)
	}
	if k.Type() != KeyFocusOut || k.Mod() != ModNone {
		t.Errorf("invalid modifier flags or key type: %s", k)
	}

	// without focus decoding
	input = NewInput()

	in = "\x1b[O"
	k, err = input.ReadKey(strings.NewReader(in))
	if err != nil {
		log.Fatal(err)
	}
	if k.Type() != KeyESCSeq || k.Mod() != ModNone {
		t.Errorf("invalid modifier flags or key type: %s", k)
	}

	in = "\x1b[I"
	k, err = input.ReadKey(strings.NewReader(in))
	if err != nil {
		log.Fatal(err)
	}
	if k.Type() != KeyESCSeq || k.Mod() != ModNone {
		t.Errorf("invalid modifier flags or key type: %s", k)
	}
}

func TestInput_ReadKey_Mouse(t *testing.T) {
	cases := []struct {
		in      string
		m       Mod
		btn     int
		pressed bool
		x, y    int
	}{
		{"\x1b[<35;1;1M", ModNone, 0, true, 1, 1},
		{"\x1b[<0;21;13m", ModNone, 1, false, 21, 13},
		{"\x1b[<6;123;542M", ModShift, 3, true, 123, 542},
		{"\x1b[<70;1;1m", ModShift, 6, false, 1, 1},
		{"\x1b[<157;65536;65536m", ModShift | ModMeta | ModCtrl, 9, false, 65535, 65535},

		// all button IDs
		{"\x1b[<0;1;1m", ModNone, 1, false, 1, 1},
		{"\x1b[<1;1;1m", ModNone, 2, false, 1, 1},
		{"\x1b[<2;1;1m", ModNone, 3, false, 1, 1},
		{"\x1b[<3;1;1m", ModNone, 0, false, 1, 1}, // AFAICT, this should never happen (no button should be value 35)
		{"\x1b[<64;1;1m", ModNone, 4, false, 1, 1},
		{"\x1b[<65;1;1m", ModNone, 5, false, 1, 1},
		{"\x1b[<66;1;1m", ModNone, 6, false, 1, 1},
		{"\x1b[<67;1;1m", ModNone, 7, false, 1, 1},
		{"\x1b[<128;1;1m", ModNone, 8, false, 1, 1},
		{"\x1b[<129;1;1m", ModNone, 9, false, 1, 1},
		{"\x1b[<130;1;1m", ModNone, 10, false, 1, 1},
		{"\x1b[<131;1;1m", ModNone, 11, false, 1, 1},
		{"\x1b[<132;1;1m", ModShift, 8, false, 1, 1},
	}

	input := NewInput(WithMouse())
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			k, err := input.ReadKey(strings.NewReader(c.in))
			if err != nil {
				t.Fatal(err)
			}
			if k.Type() != KeyMouse {
				t.Fatalf("want key type %d, got %d", KeyMouse, k.Type())
			}
			if k.Mod() != c.m {
				t.Fatalf("want modifier flags %04b, got %04b", c.m, k.Mod())
			}

			mouse := input.Mouse()
			if mouse.ButtonID() != c.btn {
				t.Errorf("want button %d, got %d", c.btn, mouse.ButtonID())
			}
			if mouse.ButtonPressed() != c.pressed {
				t.Errorf("want pressed %t, got %t", c.pressed, mouse.ButtonPressed())
			}
			if x, y := mouse.Coords(); x != c.x || y != c.y {
				t.Errorf("want %d, %d, got %d, %d", c.x, c.y, x, y)
			}
		})
	}
}

func TestInput_ReadKey_Bytes(t *testing.T) {
	input := NewInput(WithESCSeq(make(map[string]string)))

	// before any read, Bytes returns nil
	b := input.Bytes()
	if b != nil {
		t.Fatalf("want nil bytes, got %x", b)
	}

	in := "\x1baBc"
	k, err := input.ReadKey(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if k.Type() != KeyESCSeq || k.Mod() != ModNone {
		t.Fatalf("want escape sequence, got %s", k)
	}

	// bytes return the same bytes
	b = input.Bytes()
	if !bytes.Equal([]byte(in), b) {
		t.Fatalf("want %x, got %x", []byte(in), b)
	}

	// appending to the returned bytes does not impact the internal buffer
	bb, err := ioutil.ReadFile("testdata/vt100.json")
	if err != nil {
		t.Fatal(err)
	}

	oriLen, oriCap := len(input.buf), cap(input.buf)
	_ = append(b, bb...)

	if len(input.buf) != oriLen {
		t.Errorf("appending to returned bytes impacted internal buffer: len %d => %d", oriLen, len(input.buf))
	}
	if cap(input.buf) != oriCap {
		t.Errorf("appending to returned bytes impacted internal buffer: cap %d => %d", oriCap, cap(input.buf))
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

var BenchmarkKey Key

func BenchmarkInput_ReadKey(b *testing.B) {
	cases := []string{
		"a", "B", "1", "\x00", "ø", "👪", "平",
		"\x1b[B", "\x1b[1;2C", "\x1b[I", "\x1b[<35;1;2M",
	}
	for _, c := range cases {
		input := NewInput(WithFocus(), WithMouse())
		b.Run(c, func(b *testing.B) {
			r := strings.NewReader(c)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				k, err := input.ReadKey(r)
				if err != nil {
					b.Fatal(err)
				}
				BenchmarkKey = k
				r.Reset(c)
			}
		})
	}
}

var BenchmarkBytes []byte

func BenchmarkInput_ReadKey_Bytes(b *testing.B) {
	input := NewInput(WithESCSeq(make(map[string]string)))
	data := "\x1baBc"
	r := strings.NewReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k, err := input.ReadKey(r)
		if err != nil {
			b.Fatal(err)
		}
		BenchmarkKey = k
		BenchmarkBytes = input.Bytes()
		r.Reset(data)
	}
}

var BenchmarkMouseEvent MouseEvent

func BenchmarkInput_ReadKey_Mouse(b *testing.B) {
	input := NewInput(WithMouse())
	data := "\x1b[<6;123;542M"
	r := strings.NewReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k, err := input.ReadKey(r)
		if err != nil {
			b.Fatal(err)
		}
		BenchmarkKey = k
		BenchmarkMouseEvent = input.Mouse()
		r.Reset(data)
	}
}

func BenchmarkInput_ReadKey_Multiple(b *testing.B) {
	input := NewInput(WithMouse())
	data := "a⬼\x1b[<6;123;542M"
	r := strings.NewReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var count int
		for j := 0; j < 3; j++ {
			if _, err := input.ReadKey(r); err != nil {
				b.Fatal(err)
			}
			count++
		}
		if count != 3 {
			b.Fatalf("want 3 keys, got %d", count)
		}
		r.Reset(data)
	}
}

func BenchmarkInput_ReadKey_Timeout(b *testing.B) {
	input := NewInput()
	data := "⬼"
	r := strings.NewReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var count int
		for j := 0; j < 2; j++ {
			if _, err := input.ReadKey(r); err != nil {
				if err == ErrTimeout {
					break
				}
				b.Fatal(err)
			}
			count++
		}
		if count != 1 {
			b.Fatalf("want 1 key, got %d", count)
		}
		r.Reset(data)
	}
}
