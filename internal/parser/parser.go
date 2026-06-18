// Package parser trasforma la sequenza di token del lexer in una lista di
// statement: ogni statement è una riga del sorgente, con una label opzionale
// e/o un'istruzione.
package parser

import (
	"fmt"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
)

// Stmt è una riga del programma: può definire una label, contenere
// un'istruzione, o entrambe (es. "loop: ADD R1").
type Stmt struct {
	Label string            // label definita qui (vuota se assente)
	Instr *arch.Instruction // istruzione (nil se la riga ha solo una label)
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
