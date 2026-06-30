// Package source pre-elabora il sorgente .asm prima del lexer: estrae la
// direttiva di architettura ".arch <nome>", attesa come prima riga "di codice".
package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandIncludes legge path ed espande direttive ".include \"file.asm\"".
// Gli include sono locali e relativi al file che li dichiara; path assoluti e
// riferimenti che escono dalla directory del sorgente principale sono rifiutati.
func ExpandIncludes(path string) (string, error) {
	root, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	root = filepath.Clean(root)
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return expandFile(filepath.Clean(abs), root, map[string]bool{})
}

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
		case ".org", ".orgbase", ".com", ".byte", ".word", ".equ":
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

func expandFile(path string, root string, stack map[string]bool) (string, error) {
	if !insideRoot(path, root) {
		return "", fmt.Errorf("include fuori dalla directory sorgente: %s", path)
	}
	if stack[path] {
		return "", fmt.Errorf("include ciclico: %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	stack[path] = true
	defer delete(stack, path)

	dir := filepath.Dir(path)
	var out strings.Builder
	lines := strings.SplitAfter(string(data), "\n")
	for i, line := range lines {
		target, ok, err := includeTarget(line)
		if err != nil {
			return "", fmt.Errorf("%s:%d: %w", filepath.Base(path), i+1, err)
		}
		if !ok {
			out.WriteString(line)
			continue
		}
		includePath, err := resolveInclude(dir, root, target)
		if err != nil {
			return "", fmt.Errorf("%s:%d: %w", filepath.Base(path), i+1, err)
		}
		expanded, err := expandFile(includePath, root, stack)
		if err != nil {
			return "", err
		}
		out.WriteString(expanded)
		if expanded != "" && !strings.HasSuffix(expanded, "\n") && strings.HasSuffix(line, "\n") {
			out.WriteByte('\n')
		}
	}
	return out.String(), nil
}

func includeTarget(line string) (target string, ok bool, err error) {
	code := strings.TrimSpace(stripComment(line))
	if code == "" {
		return "", false, nil
	}
	if !strings.HasPrefix(code, ".include") {
		return "", false, nil
	}
	fields := strings.Fields(code)
	if len(fields) != 2 || fields[0] != ".include" {
		return "", true, fmt.Errorf("sintassi: .include \"file.asm\"")
	}
	value := fields[1]
	if len(value) < 2 || value[0] != '"' || value[len(value)-1] != '"' {
		return "", true, fmt.Errorf("sintassi: .include \"file.asm\"")
	}
	value = strings.Trim(value, `"`)
	if value == "" {
		return "", true, fmt.Errorf("include vuoto")
	}
	return value, true, nil
}

func resolveInclude(dir string, root string, target string) (string, error) {
	if filepath.IsAbs(target) {
		return "", fmt.Errorf("include assoluto non consentito: %s", target)
	}
	path := filepath.Clean(filepath.Join(dir, filepath.Clean(target)))
	if !insideRoot(path, root) {
		return "", fmt.Errorf("include fuori dalla directory sorgente: %s", target)
	}
	return path, nil
}

func insideRoot(path string, root string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}
