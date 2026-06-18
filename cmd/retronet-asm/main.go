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

func main() {
	if len(os.Args) < 2 || os.Args[1] != "build" {
		fmt.Fprintln(os.Stderr, "uso: retronet-asm build <file.asm> [-o <out.rom>]")
		os.Exit(2)
	}

	fs := flag.NewFlagSet("build", flag.ExitOnError)
	out := fs.String("o", "", "file ROM di output (default: <input>.rom)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "uso: retronet-asm build <file.asm> [-o <out.rom>]")
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[2:])

	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}
	inPath := fs.Arg(0)

	src, err := os.ReadFile(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "errore lettura %q: %v\n", inPath, err)
		os.Exit(1)
	}

	code, err := assemble(string(src))
	if err != nil {
		fmt.Fprintf(os.Stderr, "errore di assemblaggio: %v\n", err)
		os.Exit(1)
	}

	outPath := *out
	if outPath == "" {
		outPath = strings.TrimSuffix(inPath, filepath.Ext(inPath)) + ".rom"
	}
	if err := os.WriteFile(outPath, code, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "errore scrittura %q: %v\n", outPath, err)
		os.Exit(1)
	}
	fmt.Printf("assemblato %s -> %s (%d byte)\n", inPath, outPath, len(code))
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
