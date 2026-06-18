package parser

import (
	"testing"

	"github.com/retronet-labs/retronet-asm/internal/lexer"
)

func mustParse(t *testing.T, src string) []Stmt {
	t.Helper()
	toks, err := lexer.Tokenize(src)
	if err != nil {
		t.Fatalf("lexer: %v", err)
	}
	stmts, err := Parse(toks)
	if err != nil {
		t.Fatalf("parser: %v", err)
	}
	return stmts
}

func TestParseInstruction(t *testing.T) {
	st := mustParse(t, "ADD R1\n")
	if len(st) != 1 {
		t.Fatalf("statement = %d, atteso 1", len(st))
	}
	if st[0].Label != "" {
		t.Errorf("Label = %q, attesa vuota", st[0].Label)
	}
	in := st[0].Instr
	if in == nil || in.Mnemonic != "ADD" || len(in.Operands) != 1 || in.Operands[0] != "R1" {
		t.Errorf("Instr = %+v, atteso ADD [R1]", in)
	}
}

func TestParseLabelOnly(t *testing.T) {
	st := mustParse(t, "loop:\n")
	if len(st) != 1 || st[0].Label != "loop" || st[0].Instr != nil {
		t.Fatalf("statement = %+v, attesa label \"loop\" senza istruzione", st)
	}
}

func TestParseLabelAndInstruction(t *testing.T) {
	st := mustParse(t, "loop: ADD R1\n")
	if len(st) != 1 || st[0].Label != "loop" || st[0].Instr == nil || st[0].Instr.Mnemonic != "ADD" {
		t.Fatalf("statement = %+v, atteso loop + ADD", st)
	}
}

func TestParseCommaOptional(t *testing.T) {
	for _, src := range []string{"FIM R0, 0x35\n", "FIM R0 0x35\n"} {
		in := mustParse(t, src)[0].Instr
		if in == nil || in.Mnemonic != "FIM" || len(in.Operands) != 2 ||
			in.Operands[0] != "R0" || in.Operands[1] != "0x35" {
			t.Errorf("%q -> %+v, atteso FIM [R0 0x35]", src, in)
		}
	}
}

func TestParseMnemonicUppercased(t *testing.T) {
	st := mustParse(t, "ldm 5\n")
	if st[0].Instr.Mnemonic != "LDM" {
		t.Errorf("Mnemonic = %q, atteso LDM", st[0].Instr.Mnemonic)
	}
}

func TestParseOperandsVerbatim(t *testing.T) {
	// Le label negli operandi NON vanno normalizzate (sono case-sensitive).
	st := mustParse(t, "JUN Loop\n")
	if st[0].Instr.Operands[0] != "Loop" {
		t.Errorf("operando = %q, atteso \"Loop\" (verbatim)", st[0].Instr.Operands[0])
	}
}

func TestParseSkipsBlankAndComments(t *testing.T) {
	st := mustParse(t, "\n; solo commento\nLDM 5\n\nDAA\n")
	if len(st) != 2 {
		t.Fatalf("statement = %d, atteso 2 (LDM, DAA)", len(st))
	}
	if st[0].Instr.Mnemonic != "LDM" || st[1].Instr.Mnemonic != "DAA" {
		t.Errorf("statement = %+v, %+v, attesi LDM, DAA", st[0].Instr, st[1].Instr)
	}
}

func TestParseProgram(t *testing.T) {
	src := `        LDM 0
loop:   ADD R1
        ISZ R4, loop
halt:   JUN halt
`
	st := mustParse(t, src)
	if len(st) != 4 {
		t.Fatalf("statement = %d, atteso 4", len(st))
	}
	if st[1].Label != "loop" || st[1].Instr.Mnemonic != "ADD" {
		t.Errorf("stmt[1] = %+v, atteso loop + ADD", st[1])
	}
	if st[2].Instr.Mnemonic != "ISZ" || len(st[2].Instr.Operands) != 2 || st[2].Instr.Operands[1] != "loop" {
		t.Errorf("stmt[2] = %+v, atteso ISZ [R4 loop]", st[2].Instr)
	}
	if st[3].Line != 4 {
		t.Errorf("stmt[3].Line = %d, attesa 4", st[3].Line)
	}
}

func TestParseErrors(t *testing.T) {
	for _, src := range []string{"5 ADD\n", ", R1\n"} {
		toks, err := lexer.Tokenize(src)
		if err != nil {
			t.Fatalf("lexer(%q): %v", src, err)
		}
		if _, err := Parse(toks); err == nil {
			t.Errorf("Parse(%q) atteso errore, ottenuto nil", src)
		}
	}
}
