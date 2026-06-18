// Comando retronet-asm: assembla un sorgente .asm in una ROM .rom.
//
// Uso:
//
//	retronet-asm build <file.asm> [-o <out.rom>]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch/i4004"
	"github.com/retronet-labs/retronet-asm/internal/emitter"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
	"github.com/retronet-labs/retronet-asm/internal/parser"
)

const usage = "uso: retronet-asm build <file.asm> [-o <out.rom>]"

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

	code, err := assemble(string(src))
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
	fmt.Printf("assemblato %s -> %s (%d byte)\n", input, outPath, len(code))
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

// assemble esegue l'intera pipeline: lexer → parser → emitter (arch i4004).
func assemble(src string) ([]byte, error) {
	toks, err := lexer.Tokenize(src)
	if err != nil {
		return nil, err
	}
	stmts, err := parser.Parse(toks)
	if err != nil {
		return nil, err
	}
	return emitter.Assemble(stmts, i4004.New())
}
