package emitter

import (
	"bytes"
	"testing"

	"github.com/retronet-labs/retronet-asm/arch/i8080"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
	"github.com/retronet-labs/retronet-asm/internal/parser"
)

func asm8080(t *testing.T, src string) ([]byte, error) {
	t.Helper()
	toks, err := lexer.Tokenize(src)
	if err != nil {
		t.Fatalf("lexer: %v", err)
	}
	stmts, err := parser.Parse(toks)
	if err != nil {
		t.Fatalf("parser: %v", err)
	}
	return Assemble(stmts, i8080.New())
}

func mustAsm8080(t *testing.T, src string) []byte {
	t.Helper()
	code, err := asm8080(t, src)
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	return code
}

func TestOrgBaseResolvesLabelsWithoutPadding(t *testing.T) {
	code := mustAsm8080(t, ".orgbase 0x0100\nLXI D,msg\nMVI C,9\nCALL 0x0005\nmsg: .byte 0x48, 0x49, 0x24\n")
	want := []byte{0x11, 0x08, 0x01, 0x0E, 0x09, 0xCD, 0x05, 0x00, 0x48, 0x49, 0x24}
	if !bytes.Equal(code, want) {
		t.Fatalf("code=% X want=% X", code, want)
	}
}

func TestCOMAliasUsesCPMOrigin(t *testing.T) {
	code := mustAsm8080(t, ".com\nJMP start\nstart: HLT\n")
	want := []byte{0xC3, 0x03, 0x01, 0x76}
	if !bytes.Equal(code, want) {
		t.Fatalf("code=% X want=% X", code, want)
	}
}

func TestOrgBaseAfterCodeError(t *testing.T) {
	if _, err := asm8080(t, "NOP\n.orgbase 0x0100\n"); err == nil {
		t.Fatal("atteso errore per .orgbase dopo codice")
	}
}
