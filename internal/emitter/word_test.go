package emitter

import (
	"bytes"
	"testing"

	"github.com/retronet-labs/retronet-asm/arch/i6502"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
	"github.com/retronet-labs/retronet-asm/internal/parser"
)

func mustAsm6502(t *testing.T, src string) []byte {
	t.Helper()
	toks, err := lexer.Tokenize(src)
	if err != nil {
		t.Fatalf("lexer: %v", err)
	}
	stmts, err := parser.Parse(toks)
	if err != nil {
		t.Fatalf("parser: %v", err)
	}
	code, err := Assemble(stmts, i6502.New())
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	return code
}

func TestWordDirectiveLittleEndianAndLabels(t *testing.T) {
	code := mustAsm6502(t, ".orgbase $8000\nreset: NOP\n.word reset, $1234\n")
	want := []byte{0xEA, 0x00, 0x80, 0x34, 0x12}
	if !bytes.Equal(code, want) {
		t.Fatalf("code=% X want=% X", code, want)
	}
}

func TestAssemble6502MOSOperands(t *testing.T) {
	code := mustAsm6502(t, ".orgbase $8000\nstart: LDA #$01\nSTA $0200\nBNE start\n")
	want := []byte{0xA9, 0x01, 0x8D, 0x00, 0x02, 0xD0, 0xF9}
	if !bytes.Equal(code, want) {
		t.Fatalf("code=% X want=% X", code, want)
	}
}
