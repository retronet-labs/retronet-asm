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

func TestEncodeImmediate(t *testing.T) {
	tests := []struct {
		mnem string
		arg  string
		want []byte
	}{
		{"LAI", "0x2A", []byte{0x06, 0x2A}},
		{"LMI", "0xFF", []byte{0x3E, 0xFF}},
		{"LLI", "5", []byte{0x36, 0x05}},
		{"ADI", "5", []byte{0x04, 0x05}},
		{"CPI", "0x10", []byte{0x3C, 0x10}},
		{"XRI", "0", []byte{0x2C, 0x00}},
	}
	for _, tt := range tests {
		if b := enc(t, tt.mnem, tt.arg); !bytesEqual(b, tt.want) {
			t.Errorf("%s %s = % X, want % X", tt.mnem, tt.arg, b, tt.want)
		}
	}
}

func TestEncodeAddressNumeric(t *testing.T) {
	tests := []struct {
		mnem string
		arg  string
		want []byte
	}{
		{"JMP", "0x100", []byte{0x44, 0x00, 0x01}},
		{"CAL", "0x0010", []byte{0x46, 0x10, 0x00}},
		{"JFZ", "0x0004", []byte{0x48, 0x04, 0x00}},
		{"JTP", "0x3FFF", []byte{0x78, 0xFF, 0x3F}},
		{"CTC", "256", []byte{0x62, 0x00, 0x01}},
	}
	for _, tt := range tests {
		if b := enc(t, tt.mnem, tt.arg); !bytesEqual(b, tt.want) {
			t.Errorf("%s %s = % X, want % X", tt.mnem, tt.arg, b, tt.want)
		}
	}
}

func TestEncodeAddressResolvesLabel(t *testing.T) {
	resolve := func(name string) (int, bool) {
		if name == "loop" {
			return 0x0123, true
		}
		return 0, false
	}
	b, err := New().Encode(arch.Instruction{Mnemonic: "JMP", Operands: []string{"loop"}, Line: 1}, 0, resolve)
	if err != nil {
		t.Fatalf("Encode(JMP loop) = %v", err)
	}
	if want := []byte{0x44, 0x23, 0x01}; !bytesEqual(b, want) {
		t.Errorf("JMP loop = % X, want % X", b, want)
	}
	if _, err := New().Encode(arch.Instruction{Mnemonic: "JMP", Operands: []string{"ignota"}, Line: 1}, 0, resolve); err == nil {
		t.Error("JMP verso label non definita: atteso errore")
	}
}

func TestEncodeRSTAndIO(t *testing.T) {
	tests := []struct {
		mnem string
		arg  string
		want byte
	}{
		{"RST", "0", 0x05}, {"RST", "2", 0x15}, {"RST", "7", 0x3D},
		{"INP", "0", 0x41}, {"INP", "7", 0x4F},
		{"OUT", "8", 0x51}, {"OUT", "31", 0x7F},
	}
	for _, tt := range tests {
		b := enc(t, tt.mnem, tt.arg)
		if len(b) != 1 || b[0] != tt.want {
			t.Errorf("%s %s = % X, want %02X", tt.mnem, tt.arg, b, tt.want)
		}
	}
}

func TestEncodeRangeErrors(t *testing.T) {
	cases := []struct {
		mnem string
		arg  string
	}{
		{"ADI", "256"},    // immediato > 255
		{"LAI", "-1"},     // immediato < 0
		{"JMP", "0x4000"}, // indirizzo > 14 bit
		{"RST", "8"},      // vettore > 7
		{"INP", "8"},      // porta input > 7
		{"OUT", "7"},      // porta output < 8
		{"OUT", "32"},     // porta output > 31
	}
	for _, tt := range cases {
		if _, err := New().Encode(arch.Instruction{Mnemonic: tt.mnem, Operands: []string{tt.arg}, Line: 1}, 0, nil); err == nil {
			t.Errorf("%s %s: atteso errore di range", tt.mnem, tt.arg)
		}
	}
}

func TestSizeMultiByte(t *testing.T) {
	a := New()
	for _, tt := range []struct {
		mnem string
		want int
	}{
		{"LAI", 2}, {"ADI", 2}, {"JMP", 3}, {"CTP", 3}, {"RST", 1}, {"INP", 1}, {"OUT", 1},
	} {
		n, err := a.Size(arch.Instruction{Mnemonic: tt.mnem, Operands: []string{"0"}, Line: 1})
		if err != nil || n != tt.want {
			t.Errorf("Size(%s) = %d, %v; want %d", tt.mnem, n, err, tt.want)
		}
	}
}

// Un immediato può essere un simbolo (costante .equ / label) risolto da resolve.
func TestEncodeImmediateResolvesSymbol(t *testing.T) {
	resolve := func(name string) (int, bool) {
		if name == "COUNT" {
			return 5, true
		}
		return 0, false
	}
	b, err := New().Encode(arch.Instruction{Mnemonic: "LBI", Operands: []string{"COUNT"}, Line: 1}, 0, resolve)
	if err != nil {
		t.Fatalf("Encode(LBI COUNT) = %v", err)
	}
	if !bytesEqual(b, []byte{0x0E, 0x05}) { // LBI = 0x0E, COUNT = 5
		t.Errorf("LBI COUNT = % X, want 0E 05", b)
	}
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
