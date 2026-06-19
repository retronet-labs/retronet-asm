package emitter

import (
	"bytes"
	"testing"
)

func TestByteEmitsLiterals(t *testing.T) {
	code := mustAsm(t, ".byte 0x41, 0x42, 0x43\n")
	if want := []byte{0x41, 0x42, 0x43}; !bytes.Equal(code, want) {
		t.Errorf("code = % X, want % X", code, want)
	}
}

// Una label davanti a .byte punta all'indirizzo dei dati; un salto la risolve.
func TestByteLabelResolvesToDataAddress(t *testing.T) {
	// NOP @0x000; JUN tab @0x001-0x002 (tab=0x003); dati @0x003.
	code := mustAsm(t, "NOP\nJUN tab\ntab: .byte 0xAA, 0xBB\n")
	if want := []byte{0x00, 0x40, 0x03, 0xAA, 0xBB}; !bytes.Equal(code, want) {
		t.Errorf("code = % X, want % X", code, want)
	}
}
