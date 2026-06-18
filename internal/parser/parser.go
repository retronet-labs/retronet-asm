// Package parser trasforma la sequenza di token del lexer in una lista di
// statement: ogni statement è una riga del sorgente, con una label opzionale
// e/o un'istruzione.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
)

// Stmt è una riga del programma: può definire una label, contenere
// un'istruzione, o essere una direttiva ".org" (le tre cose sono esclusive,
// tranne label + istruzione che possono coesistere, es. "loop: ADD R1").
type Stmt struct {
	Label string            // label definita qui (vuota se assente)
	Instr *arch.Instruction // istruzione (nil se la riga ha solo una label)
	Org   *int              // se non-nil: ".org <Org>" posiziona il codice qui
	Line  int               // riga sorgente (1-based)
}

// Parse converte i token in statement. Una riga ha la forma:
//
//	[ label ':' ] [ mnemonico [ operando { ',' operando } ] ]
//
// Le virgole tra operandi sono separatori opzionali. I mnemonici vengono
// normalizzati in MAIUSCOLO; gli operandi restano col testo originale (le
// label sono case-sensitive).
func Parse(toks []lexer.Token) ([]Stmt, error) {
	var stmts []Stmt
	i := 0

	for i < len(toks) {
		switch toks[i].Type {
		case lexer.Newline:
			i++ // riga vuota
			continue
		case lexer.EOF:
			return stmts, nil
		}

		line := toks[i].Line
		st := Stmt{Line: line}

		// Direttiva (es. ".org <indirizzo>"): statement a sé, niente label/istruzione.
		if toks[i].Type == lexer.Directive {
			if strings.ToLower(toks[i].Text) != ".org" {
				return nil, fmt.Errorf("riga %d: direttiva sconosciuta %q", line, toks[i].Text)
			}
			i++
			if i >= len(toks) || toks[i].Type != lexer.Number {
				return nil, fmt.Errorf("riga %d: sintassi: .org <indirizzo>", line)
			}
			addr, err := parseNum(toks[i].Text)
			if err != nil {
				return nil, fmt.Errorf("riga %d: %w", line, err)
			}
			i++
			st.Org = &addr
			if i < len(toks) && toks[i].Type != lexer.Newline && toks[i].Type != lexer.EOF {
				return nil, fmt.Errorf("riga %d: token inatteso %q dopo .org", toks[i].Line, toks[i].Text)
			}
			stmts = append(stmts, st)
			if i < len(toks) && toks[i].Type == lexer.Newline {
				i++
			}
			continue
		}

		// Label opzionale: Ident seguito da ':'
		if toks[i].Type == lexer.Ident && i+1 < len(toks) && toks[i+1].Type == lexer.Colon {
			st.Label = toks[i].Text
			i += 2
		}

		// Istruzione opzionale: Ident (mnemonico) + operandi fino a fine riga
		if i < len(toks) && toks[i].Type == lexer.Ident {
			mnem := strings.ToUpper(toks[i].Text)
			i++
			var ops []string
			for i < len(toks) && toks[i].Type != lexer.Newline && toks[i].Type != lexer.EOF {
				switch toks[i].Type {
				case lexer.Ident, lexer.Number:
					ops = append(ops, toks[i].Text)
				case lexer.Comma:
					// separatore, ignorato
				default:
					return nil, fmt.Errorf("riga %d: token inatteso %q nell'istruzione",
						toks[i].Line, toks[i].Text)
				}
				i++
			}
			st.Instr = &arch.Instruction{Mnemonic: mnem, Operands: ops, Line: line}
		}

		// Dopo label e istruzione deve esserci fine riga.
		if i < len(toks) && toks[i].Type != lexer.Newline && toks[i].Type != lexer.EOF {
			return nil, fmt.Errorf("riga %d: token inatteso %q", toks[i].Line, toks[i].Text)
		}
		if st.Label == "" && st.Instr == nil {
			return nil, fmt.Errorf("riga %d: riga non valida", line)
		}

		stmts = append(stmts, st)

		if i < len(toks) && toks[i].Type == lexer.Newline {
			i++ // consuma il fine riga
		}
	}
	return stmts, nil
}

// parseNum interpreta un numero decimale (`256`) o esadecimale (`0x100`).
func parseNum(s string) (int, error) {
	t := strings.TrimSpace(s)
	var n int64
	var err error
	if strings.HasPrefix(strings.ToLower(t), "0x") {
		n, err = strconv.ParseInt(t[2:], 16, 32)
	} else {
		n, err = strconv.ParseInt(t, 10, 32)
	}
	if err != nil {
		return 0, fmt.Errorf("numero non valido %q", s)
	}
	return int(n), nil
}
