// Package lexer trasforma il sorgente .asm in una sequenza di token.
// È indipendente dall'architettura: distingue solo numeri, identificatori,
// ':' , ',' e fine riga, saltando spazi e commenti (';').
package lexer

import "fmt"

// Type è il tipo di un token.
type Type int

const (
	EOF     Type = iota // fine input
	Newline             // fine riga (separa gli statement)
	Ident               // mnemonico, registro o label: inizia con lettera o '_'
	Number              // numero decimale o esadecimale: inizia con una cifra
	Colon               // ':' (definizione di label)
	Comma               // ',' (separatore di operandi)
	Directive           // direttiva: '.' seguito da lettere (es. ".org")
)

// Token è un'unità lessicale con il testo originale e la riga (1-based).
type Token struct {
	Type Type
	Text string
	Line int
}

// Tokenize scompone src in token. Restituisce errore su un carattere inatteso.
func Tokenize(src string) ([]Token, error) {
	var toks []Token
	line := 1
	i, n := 0, len(src)

	for i < n {
		c := src[i]
		switch {
		case c == '\n':
			toks = append(toks, Token{Newline, "", line})
			line++
			i++
		case c == '\r' || c == ' ' || c == '\t':
			i++ // spazi ignorati
		case c == ';':
			for i < n && src[i] != '\n' { // commento fino a fine riga
				i++
			}
		case c == ':':
			toks = append(toks, Token{Colon, ":", line})
			i++
		case c == ',':
			toks = append(toks, Token{Comma, ",", line})
			i++
		case c == '.':
			j := i + 1
			for j < n && isIdent(src[j]) { // ".org", ".arch", ...
				j++
			}
			if j == i+1 {
				return nil, fmt.Errorf("riga %d: direttiva vuota dopo '.'", line)
			}
			toks = append(toks, Token{Directive, src[i:j], line})
			i = j
		case isDigit(c):
			j := i
			for j < n && isAlnum(src[j]) { // include "0x0C"
				j++
			}
			toks = append(toks, Token{Number, src[i:j], line})
			i = j
		case isIdentStart(c):
			j := i
			for j < n && isIdent(src[j]) {
				j++
			}
			toks = append(toks, Token{Ident, src[i:j], line})
			i = j
		default:
			return nil, fmt.Errorf("riga %d: carattere inatteso %q", line, string(c))
		}
	}
	toks = append(toks, Token{EOF, "", line})
	return toks, nil
}

func isDigit(c byte) bool      { return c >= '0' && c <= '9' }
func isAlpha(c byte) bool      { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isAlnum(c byte) bool      { return isAlpha(c) || isDigit(c) }
func isIdentStart(c byte) bool { return isAlpha(c) || c == '_' }
func isIdent(c byte) bool      { return isIdentStart(c) || isDigit(c) }
