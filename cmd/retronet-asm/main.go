// Comando retronet-asm: assembla un sorgente .asm in una ROM .rom.
//
// Uso:
//
//	retronet-asm build <file.asm> [-o <out.rom>]
//
// L'architettura si sceglie con una direttiva ".arch <nome>" sulla prima riga
// di codice del sorgente (default: i4004 se assente).
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
	"github.com/retronet-labs/retronet-asm/arch/i4004"
	"github.com/retronet-labs/retronet-asm/arch/i8008"
	"github.com/retronet-labs/retronet-asm/arch/i8080"
	"github.com/retronet-labs/retronet-asm/internal/emitter"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
	"github.com/retronet-labs/retronet-asm/internal/parser"
	"github.com/retronet-labs/retronet-asm/internal/source"
)

const usage = "uso: retronet-asm build <file.asm> [-o <out.rom>]"

// defaultArch è l'architettura usata se il sorgente non ha la direttiva .arch.
const defaultArch = "i4004"

// arches è il registro delle architetture supportate (nome → costruttore).
var arches = map[string]func() arch.Arch{
	"i4004": i4004.New,
	"i8008": i8008.New,
	"i8080": i8080.New,
}

func main() {
	if len(os.Args) < 2 || os.Args[1] != "build" {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	// Separa l'argomento posizionale (il file .asm) dai flag, così l'ordine non
	// conta: "build f.asm -o o.rom" e "build -o o.rom f.asm" funzionano entrambi
	// (il package flag, da solo, si ferma al primo argomento posizionale).
	input, flagArgs, err := splitArgs(os.Args[2:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	fs := flag.NewFlagSet("build", flag.ExitOnError)
	out := fs.String("o", "", "file ROM di output (default: <input>.rom)")
	_ = fs.Parse(flagArgs)

	src, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "errore lettura %q: %v\n", input, err)
		os.Exit(1)
	}

	// La direttiva .arch (prima riga di codice) sceglie l'architettura.
	archName, body, err := source.SplitDirective(string(src))
	if err != nil {
		fmt.Fprintf(os.Stderr, "errore: %v\n", err)
		os.Exit(1)
	}
	if archName == "" {
		archName = defaultArch
	}
	mk, ok := arches[archName]
	if !ok {
		fmt.Fprintf(os.Stderr, "architettura sconosciuta %q (disponibili: %s)\n", archName, available())
		os.Exit(1)
	}

	code, err := assemble(body, mk())
	if err != nil {
		fmt.Fprintf(os.Stderr, "errore di assemblaggio: %v\n", err)
		os.Exit(1)
	}

	outPath := *out
	if outPath == "" {
		outPath = strings.TrimSuffix(input, filepath.Ext(input)) + ".rom"
	}
	if err := os.WriteFile(outPath, code, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "errore scrittura %q: %v\n", outPath, err)
		os.Exit(1)
	}
	fmt.Printf("assemblato %s (%s) -> %s (%d byte)\n", input, archName, outPath, len(code))
}

// splitArgs estrae l'unico argomento posizionale (il file .asm) dai flag,
// riconoscendo "-o valore" e "-o=valore". Errore se i posizionali non sono uno.
func splitArgs(args []string) (input string, flagArgs []string, err error) {
	var positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-o":
			flagArgs = append(flagArgs, a)
			if i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
		case strings.HasPrefix(a, "-"):
			flagArgs = append(flagArgs, a) // -o=..., -h, ecc.
		default:
			positional = append(positional, a)
		}
	}
	if len(positional) != 1 {
		return "", nil, fmt.Errorf("atteso esattamente un file .asm, trovati %d", len(positional))
	}
	return positional[0], flagArgs, nil
}

// available elenca le architetture registrate (per i messaggi d'errore).
func available() string {
	names := make([]string, 0, len(arches))
	for n := range arches {
		names = append(names, n)
	}
	return strings.Join(names, ", ")
}

// assemble esegue l'intera pipeline: lexer → parser → emitter (architettura a).
func assemble(src string, a arch.Arch) ([]byte, error) {
	toks, err := lexer.Tokenize(src)
	if err != nil {
		return nil, err
	}
	stmts, err := parser.Parse(toks)
	if err != nil {
		return nil, err
	}
	return emitter.Assemble(stmts, a)
}
