// Package source pre-elabora il sorgente .asm prima del lexer: estrae la
// direttiva di architettura ".arch <nome>", attesa come prima riga "di codice".
package source

import (
	"fmt"
	"strings"
)

// SplitDirective cerca una direttiva ".arch <nome>" come prima riga di codice
// (saltando righe vuote e di solo commento). Se la trova, restituisce il nome
// dell'architettura e il sorgente con quella riga svuotata — la riga resta al
// suo posto (vuota) per preservare i numeri di riga negli errori successivi.
// Se non c'è direttiva, arch è "" e body è il sorgente invariato.
func SplitDirective(src string) (arch string, body string, err error) {
	lines := strings.Split(src, "\n")
	for i, ln := range lines {
		code := strings.TrimSpace(stripComment(ln))
		if code == "" {
			continue // riga vuota o di solo commento
		}
		if !strings.HasPrefix(code, ".") {
			return "", src, nil // la prima riga di codice non è una direttiva
		}
		fields := strings.Fields(code)
		switch fields[0] {
		case ".arch":
			if len(fields) != 2 {
				return "", "", fmt.Errorf("riga %d: sintassi: .arch <nome>", i+1)
			}
			lines[i] = "" // rimuovi la direttiva mantenendo il numero di riga
			return fields[1], strings.Join(lines, "\n"), nil
		case ".org":
			// Direttiva posizionale: la gestiscono lexer/parser/emitter. Nessun
			// .arch in testa → arch vuoto (default), sorgente invariato.
			return "", src, nil
		default:
			return "", "", fmt.Errorf("riga %d: direttiva sconosciuta %q", i+1, fields[0])
		}
	}
	return "", src, nil // sorgente vuoto o di soli commenti
}

// stripComment toglie un eventuale commento (da ';' a fine riga).
func stripComment(ln string) string {
	if i := strings.IndexByte(ln, ';'); i >= 0 {
		return ln[:i]
	}
	return ln
}
