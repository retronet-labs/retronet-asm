package parser

import (
	"bytes"
	"testing"

	"github.com/retronet-labs/retronet-asm/internal/lexer"
)

func TestByteStringDirective(t *testing.T) {
	toks, err := lexer.Tokenize(`msg: .byte "Hi", 0`)
	if err != nil {
		t.Fatal(err)
	}
	stmts, err := Parse(toks)
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) != 1 {
		t.Fatalf("attesi 1 statement, trovati %d", len(stmts))
	}
	if stmts[0].Label != "msg" {
		t.Errorf("label = %q, attesa msg", stmts[0].Label)
	}
	if !bytes.Equal(stmts[0].Data, []byte{'H', 'i', 0}) {
		t.Errorf(".byte data = % X, atteso 48 69 00", stmts[0].Data)
	}
}
