package i4004

import (
	"bytes"
	"testing"

	"github.com/retronet-labs/retronet-asm/arch"
)

// ins costruisce una arch.Instruction compatta per i test.
func ins(mn string, ops ...string) arch.Instruction {
	return arch.Instruction{Mnemonic: mn, Operands: ops, Line: 1}
}

func TestName(t *testing.T) {
	if got := (I4004{}).Name(); got != "i4004" {
		t.Errorf("Name() = %q, atteso \"i4004\"", got)
	}
}

func TestSize(t *testing.T) {
	a := I4004{}
	cases := []struct {
		in   arch.Instruction
		want int
	}{
		{ins("NOP"), 1},
		{ins("ADD", "R1"), 1},
		{ins("LDM", "5"), 1},
		{ins("SRC", "R2"), 1},
		{ins("JUN", "0x100"), 2},
		{ins("JCN", "2", "0x50"), 2},
		{ins("ISZ", "R4", "0x09"), 2},
		{ins("FIM", "R0", "0x35"), 2},
	}
	for _, c := range cases {
		got, err := a.Size(c.in)
		if err != nil {
			t.Errorf("Size(%s) errore inatteso: %v", c.in.Mnemonic, err)
			continue
		}
		if got != c.want {
			t.Errorf("Size(%s) = %d, atteso %d", c.in.Mnemonic, got, c.want)
		}
	}
}

func TestSizeErrors(t *testing.T) {
	a := I4004{}
	bad := []arch.Instruction{
		ins("PIPPO"),           // mnemonico sconosciuto
		ins("ADD"),             // arità: manca il registro
		ins("ADD", "R1", "R2"), // arità: troppi operandi
		ins("NOP", "R1"),       // arità: NOP non vuole operandi
	}
	for _, in := range bad {
		if _, err := a.Size(in); err == nil {
			t.Errorf("Size(%s) atteso errore, ottenuto nil", in.Mnemonic)
		}
	}
}

func TestEncodeOK(t *testing.T) {
	a := I4004{}
	cases := []struct {
		in   arch.Instruction
		want []byte
	}{
		// simple
		{ins("NOP"), []byte{0x00}},
		{ins("DAA"), []byte{0xFB}},
		{ins("WRM"), []byte{0xE0}},
		// reg
		{ins("ADD", "R1"), []byte{0x81}},
		{ins("LD", "R5"), []byte{0xA5}},
		// imm
		{ins("LDM", "7"), []byte{0xD7}},
		{ins("LDM", "0x0C"), []byte{0xDC}},
		{ins("BBL", "0"), []byte{0xC0}},
		// regPair (forza il registro pari)
		{ins("SRC", "R2"), []byte{0x23}},
		{ins("SRC", "R3"), []byte{0x23}}, // R3 → coppia pari R2
		{ins("FIN", "R0"), []byte{0x30}},
		// addr12
		{ins("JUN", "0x123"), []byte{0x41, 0x23}},
		{ins("JMS", "0x005"), []byte{0x50, 0x05}},
		// condAddr
		{ins("JCN", "2", "0x50"), []byte{0x12, 0x50}},
		// regAddr
		{ins("ISZ", "R4", "0x09"), []byte{0x74, 0x09}},
		// regImm
		{ins("FIM", "R0", "0x35"), []byte{0x20, 0x35}},
		{ins("FIM", "R2", "0x10"), []byte{0x22, 0x10}},
	}
	for _, c := range cases {
		got, err := a.Encode(c.in, 0, nil)
		if err != nil {
			t.Errorf("Encode(%s %v) errore inatteso: %v", c.in.Mnemonic, c.in.Operands, err)
			continue
		}
		if !bytes.Equal(got, c.want) {
			t.Errorf("Encode(%s %v) = % X, atteso % X", c.in.Mnemonic, c.in.Operands, got, c.want)
		}
	}
}

func TestEncodeResolvesLabel(t *testing.T) {
	a := I4004{}
	resolve := func(name string) (int, bool) {
		switch name {
		case "loop":
			return 0x123, true
		case "halt":
			return 0x005, true
		}
		return 0, false
	}
	cases := []struct {
		in   arch.Instruction
		want []byte
	}{
		{ins("JUN", "loop"), []byte{0x41, 0x23}},
		{ins("JUN", "halt"), []byte{0x40, 0x05}},
		{ins("ISZ", "R6", "loop"), []byte{0x76, 0x23}}, // di una label conta solo il byte basso
	}
	for _, c := range cases {
		got, err := a.Encode(c.in, 0, resolve)
		if err != nil {
			t.Errorf("Encode(%s %v) errore inatteso: %v", c.in.Mnemonic, c.in.Operands, err)
			continue
		}
		if !bytes.Equal(got, c.want) {
			t.Errorf("Encode(%s %v) = % X, atteso % X", c.in.Mnemonic, c.in.Operands, got, c.want)
		}
	}
}

func TestEncodeErrors(t *testing.T) {
	a := I4004{}
	resolveNone := func(string) (int, bool) { return 0, false }
	cases := []struct {
		name    string
		in      arch.Instruction
		resolve arch.Resolver
	}{
		{"mnemonico ignoto", ins("PIPPO"), nil},
		{"registro fuori range", ins("ADD", "R16"), nil},
		{"registro non valido", ins("ADD", "X"), nil},
		{"numero non valido", ins("LDM", "abc"), nil},
		{"immediato fuori range", ins("LDM", "20"), nil},
		{"dato FIM fuori range", ins("FIM", "R0", "0x1FF"), nil},
		{"indirizzo fuori 12 bit", ins("JUN", "0x1000"), nil},
		{"label non definita", ins("JUN", "manca"), resolveNone},
		{"arità sbagliata", ins("ADD"), nil},
	}
	for _, c := range cases {
		if _, err := a.Encode(c.in, 0, c.resolve); err == nil {
			t.Errorf("%s: atteso errore, ottenuto nil", c.name)
		}
	}
}
