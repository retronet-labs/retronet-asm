package source

import "testing"

func TestSplitDirectivePresent(t *testing.T) {
	arch, body, err := SplitDirective(".arch i8080\nLDM 5\n")
	if err != nil {
		t.Fatal(err)
	}
	if arch != "i8080" {
		t.Errorf("arch = %q, atteso i8080", arch)
	}
	// la direttiva è rimossa ma la riga resta vuota (numeri di riga preservati)
	if body != "\nLDM 5\n" {
		t.Errorf("body = %q, atteso \"\\nLDM 5\\n\"", body)
	}
}

func TestSplitDirectiveAfterCommentsAndBlanks(t *testing.T) {
	arch, _, err := SplitDirective("; intestazione\n\n.arch i4004 ; commento\nNOP\n")
	if err != nil {
		t.Fatal(err)
	}
	if arch != "i4004" {
		t.Errorf("arch = %q, atteso i4004", arch)
	}
}

func TestSplitDirectiveAbsent(t *testing.T) {
	src := "LDM 5\nDAA\n"
	arch, body, err := SplitDirective(src)
	if err != nil {
		t.Fatal(err)
	}
	if arch != "" {
		t.Errorf("arch = %q, atteso vuoto", arch)
	}
	if body != src {
		t.Errorf("body modificato: %q", body)
	}
}

func TestSplitDirectiveErrors(t *testing.T) {
	for _, src := range []string{".arch\n", ".arch a b\n", ".foo i4004\n"} {
		if _, _, err := SplitDirective(src); err == nil {
			t.Errorf("SplitDirective(%q): atteso errore, ottenuto nil", src)
		}
	}
}
