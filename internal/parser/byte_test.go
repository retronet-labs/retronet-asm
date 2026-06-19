package parser

import (
	"testing"

	"github.com/retronet-labs/retronet-asm/internal/lexer"
)

func TestParseByte(t *testing.T) {
	st := mustParse(t, ".byte 1, 0x02, 255\n")
	if len(st) != 1 {
		t.Fatalf("statement = %d, atteso 1", len(st))
	}
	if !sameBytes(st[0].Data, []byte{1, 2, 255}) {
		t.Fatalf("Data = %v, atteso [1 2 255]", st[0].Data)
	}
	if st[0].Instr != nil || st[0].Org != nil {
		t.Errorf(".byte non deve avere Instr/Org: %+v", st[0])
	}
}

func TestParseByteWithLabel(t *testing.T) {
	st := mustParse(t, "tab: .byte 0x41, 0x42\n")
	if st[0].Label != "tab" || !sameBytes(st[0].Data, []byte{0x41, 0x42}) {
		t.Fatalf("statement = %+v, atteso label tab + [41 42]", st[0])
	}
}

func TestParseByteErrors(t *testing.T) {
	// senza valori, valore fuori range, operando non numerico.
	for _, src := range []string{".byte\n", ".byte 256\n", ".byte R1\n"} {
		toks, err := lexer.Tokenize(src)
		if err != nil {
			continue
		}
		if _, err := Parse(toks); err == nil {
			t.Errorf("Parse(%q): atteso errore", src)
		}
	}
}

func sameBytes(a, b []byte) bool {
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
