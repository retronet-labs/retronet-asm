package parser

import (
	"testing"

	"github.com/retronet-labs/retronet-asm/internal/lexer"
)

func TestParseOrg(t *testing.T) {
	st := mustParse(t, ".org 0x100\nLDM 1\n")
	if len(st) != 2 {
		t.Fatalf("statement = %d, atteso 2", len(st))
	}
	if st[0].Org == nil || *st[0].Org != 0x100 {
		t.Errorf("Org = %v, atteso 0x100", st[0].Org)
	}
	if st[0].Instr != nil || st[0].Label != "" {
		t.Errorf(".org non deve avere label/istruzione: %+v", st[0])
	}
	if st[1].Instr == nil || st[1].Instr.Mnemonic != "LDM" {
		t.Errorf("secondo statement = %+v, atteso LDM", st[1])
	}
}

func TestParseOrgDecimal(t *testing.T) {
	st := mustParse(t, ".org 256\n")
	if st[0].Org == nil || *st[0].Org != 256 {
		t.Errorf("Org = %v, atteso 256", st[0].Org)
	}
}

func TestParseOrgBaseAndCOM(t *testing.T) {
	st := mustParse(t, ".orgbase 0x100\n.com\n")
	if len(st) != 2 {
		t.Fatalf("statement = %d, atteso 2", len(st))
	}
	if st[0].OrgBase == nil || *st[0].OrgBase != 0x100 {
		t.Fatalf("OrgBase = %v, atteso 0x100", st[0].OrgBase)
	}
	if st[1].OrgBase == nil || *st[1].OrgBase != 0x100 {
		t.Fatalf(".com OrgBase = %v, atteso 0x100", st[1].OrgBase)
	}
}

func TestParseOrgErrors(t *testing.T) {
	// .org senza operando, con operando non numerico, direttiva sconosciuta,
	// operandi in eccesso.
	for _, src := range []string{".org\n", ".org R1\n", ".bad 1\n", ".org 1 2\n", ".orgbase\n", ".orgbase R1\n", ".com 1\n"} {
		toks, err := lexer.Tokenize(src)
		if err != nil {
			continue // se l'errore è già nel lexer va bene
		}
		if _, err := Parse(toks); err == nil {
			t.Errorf("Parse(%q): atteso errore", src)
		}
	}
}
