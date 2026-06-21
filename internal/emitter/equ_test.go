package emitter

import (
	"bytes"
	"testing"
)

// Una costante .equ è usabile come immediato (qui LDM i4004: 7 -> 0xD7).
func TestEquConstantInImmediate(t *testing.T) {
	code := mustAsm(t, ".equ COUNT 7\nLDM COUNT\n")
	if want := []byte{0xD7}; !bytes.Equal(code, want) {
		t.Errorf("code = % X, want % X", code, want)
	}
}

// Costante usata prima della definizione e come indirizzo di salto.
func TestEquForwardReferenceAndAddress(t *testing.T) {
	code := mustAsm(t, "JUN TARGET\n.equ TARGET 0x010\n")
	if want := []byte{0x40, 0x10}; !bytes.Equal(code, want) {
		t.Errorf("code = % X, want % X", code, want)
	}
}

// Nome duplicato (costante che collide con una label) -> errore.
func TestEquDuplicateSymbol(t *testing.T) {
	if _, err := asm(t, ".equ FOO 1\nFOO: NOP\n"); err == nil {
		t.Error("simbolo duplicato (FOO costante e label): atteso errore")
	}
}
