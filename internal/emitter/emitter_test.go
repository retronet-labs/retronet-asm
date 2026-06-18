package emitter

import (
	"bytes"
	"testing"

	"github.com/retronet-labs/retronet-asm/arch/i4004"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
	"github.com/retronet-labs/retronet-asm/internal/parser"
)

// asm esegue la pipeline restituendo (byte, err): err per i casi negativi.
func asm(t *testing.T, src string) ([]byte, error) {
	t.Helper()
	toks, err := lexer.Tokenize(src)
	if err != nil {
		t.Fatalf("lexer: %v", err)
	}
	stmts, err := parser.Parse(toks)
	if err != nil {
		t.Fatalf("parser: %v", err)
	}
	return Assemble(stmts, i4004.New())
}

func mustAsm(t *testing.T, src string) []byte {
	t.Helper()
	code, err := asm(t, src)
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	return code
}

// Riferimento incrociato con retronet-4004/testdata: gli stessi byte che
// l'emulatore esegue correttamente devono uscire dall'assembler.

func TestAssembleBCDGolden(t *testing.T) {
	src := `        FIM R0, 0x09
        LDM 8
        ADD R1
        DAA
halt:   JUN halt
`
	want := []byte{0x20, 0x09, 0xD8, 0x81, 0xFB, 0x40, 0x05} // == bcd-add.rom
	if got := mustAsm(t, src); !bytes.Equal(got, want) {
		t.Errorf("bcd =\n % X\natteso\n % X", got, want)
	}
}

func TestAssembleMoltiplicazioneGolden(t *testing.T) {
	src := `        LDM 0
        DCL
        FIM R0, 0x03
        FIM R2, 0x00
        SRC R2
        LDM 12
        XCH R4
loop:   ADD R1
        ISZ R4, loop
        WRM
halt:   JUN halt
`
	want := []byte{
		0xD0, 0xFD, 0x20, 0x03, 0x22, 0x00, 0x23, 0xDC,
		0xB4, 0x81, 0x74, 0x09, 0xE0, 0x40, 0x0D,
	} // == moltiplicazione.rom
	if got := mustAsm(t, src); !bytes.Equal(got, want) {
		t.Errorf("moltiplicazione =\n % X\natteso\n % X", got, want)
	}
}

func TestAssembleForwardLabel(t *testing.T) {
	// 'end' è riferita PRIMA di essere definita: la risolve la passata 1.
	src := `        JUN end
        NOP
end:    JUN end
`
	want := []byte{0x40, 0x03, 0x00, 0x40, 0x03}
	if got := mustAsm(t, src); !bytes.Equal(got, want) {
		t.Errorf("forward label = % X, atteso % X", got, want)
	}
}

func TestAssembleErrors(t *testing.T) {
	cases := map[string]string{
		"label duplicata":    "x: NOP\nx: NOP\n",
		"label non definita": "JUN manca\n",
		"mnemonico ignoto":   "PIPPO\n",
		"arità sbagliata":    "ADD\n",
	}
	for name, src := range cases {
		if _, err := asm(t, src); err == nil {
			t.Errorf("%s: atteso errore, ottenuto nil", name)
		}
	}
}
