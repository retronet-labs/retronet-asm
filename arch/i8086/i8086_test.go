package i8086

import (
	"bytes"
	"testing"

	"github.com/retronet-labs/retronet-asm/arch"
)

func enc(t *testing.T, mnem string, ops []string, pc int, resolve arch.Resolver) []byte {
	t.Helper()
	b, err := I8086{}.Encode(arch.Instruction{Mnemonic: mnem, Operands: ops, Line: 1}, pc, resolve)
	if err != nil {
		t.Fatalf("%s %v: %v", mnem, ops, err)
	}
	return b
}

func TestEncodeRegAndImm(t *testing.T) {
	cases := []struct {
		mnem string
		ops  []string
		want []byte
	}{
		{"XOR", []string{"AX", "AX"}, []byte{0x31, 0xC0}},
		{"MOV", []string{"DS", "AX"}, []byte{0x8E, 0xD8}},
		{"MOV", []string{"AX", "DS"}, []byte{0x8C, 0xD8}},
		{"MOV", []string{"BX", "AX"}, []byte{0x89, 0xC3}},
		{"MOV", []string{"SI", "0x7C16"}, []byte{0xBE, 0x16, 0x7C}},
		{"MOV", []string{"AH", "0x0E"}, []byte{0xB4, 0x0E}},
		{"ADD", []string{"AX", "BX"}, []byte{0x01, 0xD8}},
		{"CMP", []string{"AL", "0"}, []byte{0x3C, 0x00}},
		{"CMP", []string{"BL", "5"}, []byte{0x80, 0xFB, 0x05}},
		{"ADD", []string{"CX", "0x1234"}, []byte{0x81, 0xC1, 0x34, 0x12}},
		{"INT", []string{"0x10"}, []byte{0xCD, 0x10}},
		{"INC", []string{"AX"}, []byte{0x40}},
		{"DEC", []string{"BX"}, []byte{0x4B}},
		{"INC", []string{"AL"}, []byte{0xFE, 0xC0}},
		{"PUSH", []string{"AX"}, []byte{0x50}},
		{"POP", []string{"DS"}, []byte{0x1F}},
		{"PUSH", []string{"CS"}, []byte{0x0E}},
		{"SHL", []string{"AX", "1"}, []byte{0xD1, 0xE0}},
		{"SHR", []string{"BL", "CL"}, []byte{0xD2, 0xEB}},
		{"NEG", []string{"AX"}, []byte{0xF7, 0xD8}},
		{"MUL", []string{"BL"}, []byte{0xF6, 0xE3}},
		{"LODSB", nil, []byte{0xAC}},
		{"HLT", nil, []byte{0xF4}},
		{"CLD", nil, []byte{0xFC}},
	}
	for _, c := range cases {
		got := enc(t, c.mnem, c.ops, 0, nil)
		if !bytes.Equal(got, c.want) {
			t.Errorf("%s %v = % X, atteso % X", c.mnem, c.ops, got, c.want)
		}
	}
}

func TestEncodeRelativeJumps(t *testing.T) {
	resolve := func(name string) (int, bool) {
		if name == "L" {
			return 0x110, true
		}
		return 0, false
	}
	// JE rel8 da pc 0x100 a 0x110: rel = 0x110-(0x100+2) = 0x0E.
	if got := enc(t, "JE", []string{"L"}, 0x100, resolve); !bytes.Equal(got, []byte{0x74, 0x0E}) {
		t.Errorf("JE = % X", got)
	}
	// JMP near rel16 da pc 0x100 a 0x110: rel = 0x110-(0x100+3) = 0x0D.
	if got := enc(t, "JMP", []string{"L"}, 0x100, resolve); !bytes.Equal(got, []byte{0xE9, 0x0D, 0x00}) {
		t.Errorf("JMP = % X", got)
	}
	// JMP SHORT rel8.
	if got := enc(t, "JMP", []string{"SHORT", "L"}, 0x100, resolve); !bytes.Equal(got, []byte{0xEB, 0x0E}) {
		t.Errorf("JMP SHORT = % X", got)
	}
}

func TestSize(t *testing.T) {
	cases := []struct {
		mnem string
		ops  []string
		want int
	}{
		{"MOV", []string{"AX", "0x1234"}, 3},
		{"MOV", []string{"AL", "5"}, 2},
		{"MOV", []string{"AX", "BX"}, 2},
		{"INT", []string{"0x10"}, 2},
		{"JMP", []string{"loop"}, 3},
		{"JE", []string{"loop"}, 2},
		{"INC", []string{"AX"}, 1},
		{"INC", []string{"AL"}, 2},
		{"ADD", []string{"AX", "5"}, 3},
		{"ADD", []string{"BX", "5"}, 4},
		{"CMP", []string{"AL", "5"}, 2},
		{"HLT", nil, 1},
	}
	for _, c := range cases {
		got, err := (I8086{}).Size(arch.Instruction{Mnemonic: c.mnem, Operands: c.ops, Line: 1})
		if err != nil {
			t.Fatalf("%s %v: %v", c.mnem, c.ops, err)
		}
		if got != c.want {
			t.Errorf("Size(%s %v) = %d, atteso %d", c.mnem, c.ops, got, c.want)
		}
	}
}
