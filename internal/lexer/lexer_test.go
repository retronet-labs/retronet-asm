package lexer

import "testing"

// tok è una coppia tipo/testo, per confronti compatti nei test.
type tok struct {
	t Type
	s string
}

func collect(toks []Token) []tok {
	out := make([]tok, len(toks))
	for i, k := range toks {
		out[i] = tok{k.Type, k.Text}
	}
	return out
}

func TestTokenizeBasic(t *testing.T) {
	toks, err := Tokenize("loop: ADD R1, 0x0C\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []tok{
		{Ident, "loop"}, {Colon, ":"}, {Ident, "ADD"}, {Ident, "R1"},
		{Comma, ","}, {Number, "0x0C"}, {Newline, ""}, {EOF, ""},
	}
	got := collect(toks)
	if len(got) != len(want) {
		t.Fatalf("numero token = %d, atteso %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("token %d = %+v, atteso %+v", i, got[i], want[i])
		}
	}
}

func TestTokenizeNumbersDecAndHex(t *testing.T) {
	toks, err := Tokenize("12 0x0C 255")
	if err != nil {
		t.Fatal(err)
	}
	want := []tok{{Number, "12"}, {Number, "0x0C"}, {Number, "255"}, {EOF, ""}}
	if got := collect(toks); !equalToks(got, want) {
		t.Fatalf("token = %v, atteso %v", got, want)
	}
}

func TestTokenizeSkipsComments(t *testing.T) {
	toks, err := Tokenize("LDM 5 ; carica 5 in A\n; riga di solo commento\nDAA\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []tok{
		{Ident, "LDM"}, {Number, "5"}, {Newline, ""},
		{Newline, ""}, // la riga di solo commento lascia solo il suo a-capo
		{Ident, "DAA"}, {Newline, ""}, {EOF, ""},
	}
	if got := collect(toks); !equalToks(got, want) {
		t.Fatalf("token = %v, atteso %v", got, want)
	}
}

func TestTokenizeLineNumbers(t *testing.T) {
	toks, err := Tokenize("LDM 5\nDAA\n")
	if err != nil {
		t.Fatal(err)
	}
	// LDM(1) 5(1) NL(1) DAA(2) NL(2) EOF(3)
	wantLines := []int{1, 1, 1, 2, 2, 3}
	if len(toks) != len(wantLines) {
		t.Fatalf("numero token = %d, atteso %d", len(toks), len(wantLines))
	}
	for i, w := range wantLines {
		if toks[i].Line != w {
			t.Errorf("token %d (%+v) riga = %d, attesa %d", i, toks[i], toks[i].Line, w)
		}
	}
}

func TestTokenizeUnexpectedChar(t *testing.T) {
	if _, err := Tokenize("LDM @5"); err == nil {
		t.Fatal("atteso errore per carattere inatteso, ottenuto nil")
	}
}

func equalToks(a, b []tok) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
