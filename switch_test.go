package classfileparser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"testing"
)

func switchBytecode(nops int, opcode byte, operands ...int32) []byte {
	var b bytes.Buffer
	for i := 0; i < nops; i++ {
		b.WriteByte(0x00) // nop
	}
	b.WriteByte(opcode)
	for b.Len()%4 != 0 {
		b.WriteByte(0x00) // padding
	}
	binary.Write(&b, binary.BigEndian, operands)
	b.WriteByte(0xB1) // return
	return b.Bytes()
}

func TestSwitchPadding(t *testing.T) {
	cases := map[string]struct {
		opcode   byte
		operands []int32
		want     interface{}
	}{
		"tableswitch": {
			0xAA,
			[]int32{40, -1, 1, 10, 20, 30},
			Tableswitch{Default: 40, Low: -1, High: 1, Offsets: []int32{10, 20, 30}},
		},
		"lookupswitch": {
			0xAB,
			[]int32{99, 2, -5, 11, 1000, 22},
			Lookupswitch{Default: 99, Pairs: []struct{ Match, Offset int32 }{{Match: -5, Offset: 11}, {Match: 1000, Offset: 22}}},
		},
		"lookupswitch without pairs": {
			0xAB,
			[]int32{8, 0},
			Lookupswitch{Default: 8, Pairs: []struct{ Match, Offset int32 }{}},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// The opcode at addresses 0 to 3 covers the four possible padding sizes.
			for nops := 0; nops < 4; nops++ {
				got, err := parseInstructions(switchBytecode(nops, c.opcode, c.operands...))
				if err != nil {
					t.Fatalf("nops=%d: parse failed: %v", nops, err)
				}
				want := []interface{}{}
				for i := 0; i < nops; i++ {
					want = append(want, Nop{})
				}
				want = append(want, c.want, Return{})
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("nops=%d: got %#v, want %#v", nops, got, want)
				}
			}
		})
	}
}

func TestInvalidSwitches(t *testing.T) {
	cases := map[string]struct {
		bytecode []byte
		message  string
	}{
		"tableswitch low > high":      {switchBytecode(0, 0xAA, 0, 5, 1), "invalid tableswitch: low (5) > high (1)"},
		"tableswitch too many":        {switchBytecode(0, 0xAA, 0, 0, 1<<30), "invalid tableswitch: 1073741825 offsets exceed the remaining code"},
		"lookupswitch negative pairs": {switchBytecode(0, 0xAB, 0, -1), "invalid lookupswitch: bad number of pairs (-1)"},
		"lookupswitch too many pairs": {switchBytecode(0, 0xAB, 0, 1<<30), "invalid lookupswitch: bad number of pairs (1073741824)"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := parseInstructions(c.bytecode)
			want := "failed to read attribute Code: failed to read instruction at offset 0: malformed class file: " + c.message
			if !errors.Is(err, ErrMalformed) || err.Error() != want {
				t.Fatalf("got %v, want %s", err, want)
			}
		})
	}
}
