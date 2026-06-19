package i8008

import (
	"testing"

	"github.com/retronet-labs/retronet-asm/arch"
)

func enc(t *testing.T, mnem string, ops ...string) []byte {
	t.Helper()
	b, err := New().Encode(arch.Instruction{Mnemonic: mnem, Operands: ops, Line: 1}, 0, nil)
	if err != nil {
		t.Fatalf("Encode(%s) = %v", mnem, err)
	}
	return b
}

// Le codifiche attese coincidono con i pattern di bit dell'8008 (e con i byte
// che il disassembler dell'emulatore retronet-8008 produce per gli stessi mnemonici).
func TestEncodeSimple(t *testing.T) {
	tests := []struct {
		mnem string
		want byte
	}{
		{"HLT", 0x00}, {"RET", 0x07}, {"NOP", 0xC0},
		{"RLC", 0x02}, {"RRC", 0x0A}, {"RAL", 0x12}, {"RAR", 0x1A},
		{"LAB", 0xC1}, {"LBA", 0xC8}, {"LAM", 0xC7}, {"LMA", 0xF8}, {"LLH", 0xF5},
		{"ADA", 0x80}, {"ADB", 0x81}, {"ADM", 0x87}, {"CPM", 0xBF}, {"XRA", 0xA8},
		{"INB", 0x08}, {"INL", 0x30}, {"DCB", 0x09}, {"DCL", 0x31},
		{"RFC", 0x03}, {"RTC", 0x23}, {"RFZ", 0x0B}, {"RTP", 0x3B},
	}
	for _, tt := range tests {
		b := enc(t, tt.mnem)
		if len(b) != 1 || b[0] != tt.want {
			t.Errorf("%s = % X, want %02X", tt.mnem, b, tt.want)
		}
	}
}

func TestSizeOperandsAndErrors(t *testing.T) {
	a := New()

	if a.Name() != "i8008" {
		t.Fatalf("Name = %q, want i8008", a.Name())
	}
	if n, err := a.Size(arch.Instruction{Mnemonic: "LAB", Line: 1}); err != nil || n != 1 {
		t.Fatalf("Size(LAB) = %d, %v; want 1, nil", n, err)
	}
	if _, err := a.Size(arch.Instruction{Mnemonic: "ZZZ", Line: 1}); err == nil {
		t.Error("Size(ZZZ): atteso errore mnemonico sconosciuto")
	}
	if _, err := a.Size(arch.Instruction{Mnemonic: "HLT", Operands: []string{"1"}, Line: 1}); err == nil {
		t.Error("Size(HLT con operando): atteso errore di arita'")
	}
}

// dst == src non genera un move (sarebbe NOP/HLT): solo "NOP" e' accettato.
func TestSameRegisterMoveIsNotDefined(t *testing.T) {
	if _, err := New().Size(arch.Instruction{Mnemonic: "LAA", Line: 1}); err == nil {
		t.Error("LAA non deve essere un mnemonico valido (usa NOP)")
	}
}
