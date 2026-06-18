package emitter

import (
	"bytes"
	"testing"
)

// .org riempie il vuoto con NOP e posiziona la label all'indirizzo richiesto.
func TestOrgPadsAndPlacesLabel(t *testing.T) {
	code := mustAsm(t, ".org 0x004\nLDM 7\nhalt: JUN halt\n")
	// 0x000-0x003: NOP; 0x004: LDM 7 (D7); 0x005: JUN halt (40 05).
	want := []byte{0x00, 0x00, 0x00, 0x00, 0xD7, 0x40, 0x05}
	if !bytes.Equal(code, want) {
		t.Errorf("code = % X, atteso % X", code, want)
	}
}

// Una label dopo .org viene risolta correttamente da un JUN che la precede.
func TestOrgLabelResolvesAcrossPad(t *testing.T) {
	code := mustAsm(t, "JUN sub\n.org 0x100\nsub: LDM 1\n")
	if len(code) != 0x101 {
		t.Fatalf("len = 0x%X, atteso 0x101", len(code))
	}
	if code[0] != 0x41 || code[1] != 0x00 { // JUN 0x100
		t.Errorf("JUN sub = % X, atteso 41 00", code[0:2])
	}
	if code[0x100] != 0xD1 { // LDM 1
		t.Errorf("byte a 0x100 = 0x%02X, atteso 0xD1 (LDM 1)", code[0x100])
	}
}

// .org all'indietro: errore (sovrapporrebbe codice già emesso).
func TestOrgBackwardError(t *testing.T) {
	if _, err := asm(t, "LDM 0\n.org 0x000\n"); err == nil {
		t.Error("atteso errore: .org all'indietro")
	}
}

// .org oltre lo spazio ROM (12 bit): errore.
func TestOrgOutOfRange(t *testing.T) {
	if _, err := asm(t, ".org 0x1000\nLDM 0\n"); err == nil {
		t.Error("atteso errore: .org fuori dallo spazio ROM")
	}
}
