package parser

import "testing"

func TestParseWordDirective(t *testing.T) {
	st := mustParse(t, "vec: .word $8000, reset\n")
	if len(st) != 1 || st[0].Label != "vec" {
		t.Fatalf("stmt = %+v", st)
	}
	if got := st[0].Words; len(got) != 2 || got[0] != "$8000" || got[1] != "reset" {
		t.Fatalf("Words = %#v", got)
	}
}
