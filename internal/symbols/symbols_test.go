package symbols

import "testing"

func TestDefineAndLookup(t *testing.T) {
	tab := New()
	if err := tab.Define("loop", 0x09); err != nil {
		t.Fatal(err)
	}
	addr, ok := tab.Lookup("loop")
	if !ok || addr != 0x09 {
		t.Errorf("Lookup(loop) = %d,%v; atteso 9,true", addr, ok)
	}
}

func TestDefineDuplicate(t *testing.T) {
	tab := New()
	if err := tab.Define("x", 1); err != nil {
		t.Fatal(err)
	}
	if err := tab.Define("x", 2); err == nil {
		t.Error("Define duplicato: atteso errore, ottenuto nil")
	}
}

func TestLookupMissing(t *testing.T) {
	tab := New()
	if _, ok := tab.Lookup("manca"); ok {
		t.Error("Lookup di label inesistente: atteso ok=false")
	}
}
