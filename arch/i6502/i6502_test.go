package i6502

import (
	"bytes"
	"testing"

	"github.com/retronet-labs/retronet-asm/arch"
)

func enc(t *testing.T, m string, ops []string, pc int, resolve arch.Resolver) []byte {
	t.Helper()
	b, err := (I6502{}).Encode(arch.Instruction{Mnemonic: m, Operands: ops, Line: 1}, pc, resolve)
	if err != nil {
		t.Fatalf("%s %v: %v", m, ops, err)
	}
	return b
}

func TestEncodeAddressingModes(t *testing.T) {
	cases := []struct {
		m    string
		ops  []string
		want []byte
	}{
		{"LDA", []string{"#$01"}, []byte{0xA9, 0x01}},
		{"LDA", []string{"$20"}, []byte{0xA5, 0x20}},
		{"LDA", []string{"$20", "X"}, []byte{0xB5, 0x20}},
		{"LDA", []string{"$2000", "X"}, []byte{0xBD, 0x00, 0x20}},
		{"LDA", []string{"($44,X)"}, []byte{0xA1, 0x44}},
		{"LDA", []string{"($44),Y"}, []byte{0xB1, 0x44}},
		{"STA", []string{"$0200"}, []byte{0x8D, 0x00, 0x02}},
		{"JMP", []string{"($12FF)"}, []byte{0x6C, 0xFF, 0x12}},
		{"ASL", []string{"A"}, []byte{0x0A}},
		{"ASL", []string{"$20"}, []byte{0x06, 0x20}},
		{"CLC", nil, []byte{0x18}},
	}
	for _, c := range cases {
		got := enc(t, c.m, c.ops, 0, nil)
		if !bytes.Equal(got, c.want) {
			t.Errorf("%s %v = % X, atteso % X", c.m, c.ops, got, c.want)
		}
	}
}

func TestLabelsDefaultAbsoluteAndCanForceZeroPage(t *testing.T) {
	resolve := func(name string) (int, bool) {
		if name == "ptr" {
			return 0x0044, true
		}
		return 0, false
	}
	if got := enc(t, "LDA", []string{"ptr"}, 0, resolve); !bytes.Equal(got, []byte{0xAD, 0x44, 0x00}) {
		t.Fatalf("LDA ptr = % X", got)
	}
	if got := enc(t, "LDA", []string{"<ptr"}, 0, resolve); !bytes.Equal(got, []byte{0xA5, 0x44}) {
		t.Fatalf("LDA <ptr = % X", got)
	}
	if got := enc(t, "LDA", []string{"#>ptr"}, 0, resolve); !bytes.Equal(got, []byte{0xA9, 0x00}) {
		t.Fatalf("LDA #>ptr = % X", got)
	}
}

func TestRelativeBranch(t *testing.T) {
	resolve := func(name string) (int, bool) {
		if name == "loop" {
			return 0x1000, true
		}
		return 0, false
	}
	got := enc(t, "BNE", []string{"loop"}, 0x1004, resolve)
	if !bytes.Equal(got, []byte{0xD0, 0xFA}) {
		t.Fatalf("BNE loop = % X", got)
	}
}

func TestSizeDeterministicForSymbols(t *testing.T) {
	cases := []struct {
		m    string
		ops  []string
		want int
	}{
		{"LDA", []string{"$20"}, 2},
		{"LDA", []string{"label"}, 3},
		{"LDA", []string{"<label"}, 2},
		{"BNE", []string{"label"}, 2},
		{"JMP", []string{"($1234)"}, 3},
	}
	for _, c := range cases {
		got, err := (I6502{}).Size(arch.Instruction{Mnemonic: c.m, Operands: c.ops, Line: 1})
		if err != nil {
			t.Fatalf("Size(%s %v): %v", c.m, c.ops, err)
		}
		if got != c.want {
			t.Errorf("Size(%s %v)=%d, atteso %d", c.m, c.ops, got, c.want)
		}
	}
}

func TestInvalidFormRejected(t *testing.T) {
	if _, err := (I6502{}).Encode(arch.Instruction{Mnemonic: "STX", Operands: []string{"$20", "X"}, Line: 1}, 0, nil); err == nil {
		t.Fatal("STX zp,X deve essere rifiutato (solo zp,Y)")
	}
}
