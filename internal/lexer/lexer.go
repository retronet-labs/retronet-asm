// Package lexer trasforma il sorgente .asm in una sequenza di token.
// È indipendente dall'architettura: distingue solo numeri, identificatori,
// ':' , ',' e fine riga, saltando spazi e commenti (';').
package lexer

import (
	"fmt"
	"strings"
)

// Type è il tipo di un token.
type Type int

const (
	EOF       Type = iota // fine input
	Newline               // fine riga (separa gli statement)
	Ident                 // mnemonico, registro o label: inizia con lettera o '_'
	Number                // numero decimale o esadecimale: inizia con una cifra
	Colon                 // ':' (definizione di label)
	Comma                 // ',' (separatore di operandi)
	Directive             // direttiva: '.' seguito da lettere (es. ".org")
	String                // stringa tra virgolette: "testo" (per .byte)
	Mem                   // operando in memoria tra parentesi quadre: [bx+si+disp]
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
		case c == '"':
			s, end, err := scanString(src, i, line)
			if err != nil {
				return nil, err
			}
			toks = append(toks, Token{String, s, line})
			i = end
		case c == '[':
			j := i + 1
			for j < n && src[j] != ']' && src[j] != '\n' {
				j++
			}
			if j >= n || src[j] != ']' {
				return nil, fmt.Errorf("riga %d: parentesi quadra non chiusa", line)
			}
			// Testo dell'operando in memoria, spazi interni rimossi.
			inner := strings.ReplaceAll(strings.ReplaceAll(src[i+1:j], " ", ""), "\t", "")
			toks = append(toks, Token{Mem, "[" + inner + "]", line})
			i = j + 1
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

// scanString legge una stringa tra virgolette a partire da src[start] ('"'),
// interpretando gli escape \n \r \t \0 \\ \". Restituisce il testo decodificato e
// l'indice del carattere dopo la virgoletta di chiusura.
func scanString(src string, start, line int) (string, int, error) {
	var sb strings.Builder
	i, n := start+1, len(src)
	for i < n {
		c := src[i]
		switch {
		case c == '"':
			return sb.String(), i + 1, nil
		case c == '\n':
			return "", 0, fmt.Errorf("riga %d: stringa non terminata", line)
		case c == '\\' && i+1 < n:
			i++
			switch src[i] {
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case '0':
				sb.WriteByte(0)
			case '\\':
				sb.WriteByte('\\')
			case '"':
				sb.WriteByte('"')
			default:
				return "", 0, fmt.Errorf("riga %d: escape sconosciuto \\%c", line, src[i])
			}
		default:
			sb.WriteByte(c)
		}
		i++
	}
	return "", 0, fmt.Errorf("riga %d: stringa non terminata", line)
}

func isDigit(c byte) bool      { return c >= '0' && c <= '9' }
func isAlpha(c byte) bool      { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isAlnum(c byte) bool      { return isAlpha(c) || isDigit(c) }
func isIdentStart(c byte) bool { return isAlpha(c) || c == '_' }
func isIdent(c byte) bool      { return isIdentStart(c) || isDigit(c) }
