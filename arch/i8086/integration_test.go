package i8086_test

import (
	"testing"

	"github.com/retronet-labs/retronet-asm/arch/i8086"
	"github.com/retronet-labs/retronet-asm/internal/emitter"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
	"github.com/retronet-labs/retronet-asm/internal/parser"
)

// Assembla un boot sector completo attraverso l'intera pipeline e verifica
// dimensione (512), firma di boot e i primi byte (codice e indirizzo del
// messaggio risolto in spazio .orgbase 0x7C00).
func TestAssembleBootSector(t *testing.T) {
	src := `.orgbase 0x7C00
        xor ax, ax
        mov ds, ax
        mov si, msg
        mov ah, 0x0E
print:  lodsb
        cmp al, 0
        je halt
        int 0x10
        jmp print
halt:   jmp halt
msg:    .byte "Hi", 0
        .org 0x7DFE
        .byte 0x55, 0xAA
`
	toks, err := lexer.Tokenize(src)
	if err != nil {
		t.Fatal(err)
	}
	stmts, err := parser.Parse(toks)
	if err != nil {
		t.Fatal(err)
	}
	code, err := emitter.Assemble(stmts, i8086.New())
	if err != nil {
		t.Fatal(err)
	}

	if len(code) != 512 {
		t.Fatalf("dimensione = %d, attesi 512", len(code))
	}
	if code[510] != 0x55 || code[511] != 0xAA {
		t.Errorf("firma di boot mancante: %#02x %#02x", code[510], code[511])
	}
	// xor ax,ax ; mov ds,ax ; mov si,0x7C16 ...
	want := []byte{0x31, 0xC0, 0x8E, 0xD8, 0xBE, 0x16, 0x7C, 0xB4, 0x0E, 0xAC}
	for i, w := range want {
		if code[i] != w {
			t.Fatalf("byte %d = %#02x, atteso %#02x", i, code[i], w)
		}
	}
	// Il messaggio "Hi\0" deve trovarsi all'offset 0x16 (22).
	if string(code[0x16:0x19]) != "Hi\x00" {
		t.Errorf("messaggio non all'offset atteso: % X", code[0x16:0x19])
	}
}
