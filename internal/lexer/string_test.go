package lexer

import "testing"

func TestStringToken(t *testing.T) {
	toks, err := Tokenize(`.byte "AB\0"`)
	if err != nil {
		t.Fatal(err)
	}
	var got string
	found := false
	for _, tk := range toks {
		if tk.Type == String {
			got = tk.Text
			found = true
		}
	}
	if !found {
		t.Fatal("nessun token String prodotto")
	}
	if got != "AB\x00" {
		t.Errorf("stringa decodificata = %q, attesa \"AB\\x00\"", got)
	}
}

func TestUnterminatedString(t *testing.T) {
	if _, err := Tokenize(`.byte "manca la chiusura`); err == nil {
		t.Error("stringa non terminata dovrebbe dare errore")
	}
}
