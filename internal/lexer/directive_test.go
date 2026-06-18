package lexer

import "testing"

func TestTokenizeDirective(t *testing.T) {
	toks, err := Tokenize(".org 0x100\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []tok{{Directive, ".org"}, {Number, "0x100"}, {Newline, ""}, {EOF, ""}}
	if got := collect(toks); !equalToks(got, want) {
		t.Errorf("got %v, atteso %v", got, want)
	}
}

func TestTokenizeLoneDotErrors(t *testing.T) {
	if _, err := Tokenize(". 5\n"); err == nil {
		t.Error("atteso errore su '.' isolato")
	}
}
