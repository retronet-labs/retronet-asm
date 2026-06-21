package parser

import (
	"testing"

	"github.com/retronet-labs/retronet-asm/internal/lexer"
)

func TestParseEqu(t *testing.T) {
	st := mustParse(t, ".equ COUNT 5\n.equ BASE 0x100\n")
	if len(st) != 2 {
		t.Fatalf("statement = %d, atteso 2", len(st))
	}
	if st[0].Equ == nil || st[0].Equ.Name != "COUNT" || st[0].Equ.Value != 5 {
		t.Errorf("Equ[0] = %+v, atteso {COUNT 5}", st[0].Equ)
	}
	if st[1].Equ == nil || st[1].Equ.Name != "BASE" || st[1].Equ.Value != 0x100 {
		t.Errorf("Equ[1] = %+v, atteso {BASE 256}", st[1].Equ)
	}
}

func TestParseEquErrors(t *testing.T) {
	// senza nome, senza valore, nome non-identificatore, valore non numerico.
	for _, src := range []string{".equ\n", ".equ COUNT\n", ".equ 5 5\n", ".equ COUNT x\n"} {
		toks, err := lexer.Tokenize(src)
		if err != nil {
			continue
		}
		if _, err := Parse(toks); err == nil {
			t.Errorf("Parse(%q): atteso errore", src)
		}
	}
}
