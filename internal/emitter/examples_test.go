package emitter

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/retronet-labs/retronet-asm/arch"
	"github.com/retronet-labs/retronet-asm/arch/i4004"
	"github.com/retronet-labs/retronet-asm/arch/i6502"
	"github.com/retronet-labs/retronet-asm/arch/i8008"
	"github.com/retronet-labs/retronet-asm/arch/i8080"
	"github.com/retronet-labs/retronet-asm/arch/i8086"
	"github.com/retronet-labs/retronet-asm/internal/lexer"
	"github.com/retronet-labs/retronet-asm/internal/parser"
	"github.com/retronet-labs/retronet-asm/internal/source"
)

var exampleArches = map[string]func() arch.Arch{
	"i4004": i4004.New,
	"i6502": i6502.New,
	"i8008": i8008.New,
	"i8080": i8080.New,
	"i8086": i8086.New,
}

// assembleFile riproduce la pipeline della CLI sui .asm di examples/.
func assembleFile(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "examples", name)
	src, err := source.ExpandIncludes(path)
	if err != nil {
		t.Fatalf("lettura %s: %v", name, err)
	}
	archName, body, err := source.SplitDirective(src)
	if err != nil {
		t.Fatalf("%s: direttiva: %v", name, err)
	}
	if archName == "" {
		archName = "i4004"
	}
	mk, ok := exampleArches[archName]
	if !ok {
		t.Fatalf("%s: architettura %q non registrata nel test", name, archName)
	}
	toks, err := lexer.Tokenize(body)
	if err != nil {
		t.Fatalf("%s: lexer: %v", name, err)
	}
	stmts, err := parser.Parse(toks)
	if err != nil {
		t.Fatalf("%s: parser: %v", name, err)
	}
	code, err := Assemble(stmts, mk())
	if err != nil {
		t.Fatalf("%s: assemble: %v", name, err)
	}
	return code
}

func TestExamplesAssembleAndHaveDocs(t *testing.T) {
	root := filepath.Join("..", "..", "examples")
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".asm") {
			return nil
		}
		doc := strings.TrimSuffix(path, ".asm") + ".md"
		if _, err := os.Stat(doc); err != nil {
			t.Errorf("%s: documentazione mancante %s", path, doc)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		assembleFile(t, rel)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestExamplesGolden protegge i .asm dell'aritmetica BCD (Step 13) da regressioni
// dell'assembler. I byte attesi sono quelli prodotti dalla pipeline e verificati
// eseguendo le ROM su retronet-4004.
func TestExamplesGolden(t *testing.T) {
	cases := map[string][]byte{
		"i4004/06-sottrazione-bcd.asm": {
			0xD0, 0xFD, 0x20, 0x00, 0x21, 0xD5, 0xB1, 0xD7,
			0xB2, 0xFA, 0xF9, 0x91, 0x82, 0xFB, 0xE0, 0x40, 0x0F,
		},
		"i4004/07-sottrazione-multicifra.asm": {
			0xD0, 0xFD, 0x20, 0x00, 0x21, 0xD2, 0xE0, 0x20, 0x01, 0x21, 0xD5, 0xE0, 0x22, 0x10, 0x23, 0xD7,
			0xE0, 0x22, 0x11, 0x23, 0xD2, 0xE0, 0x20, 0x00, 0x22, 0x10, 0x24, 0x20, 0x26, 0xE0, 0xFA, 0xF9,
			0x23, 0xE8, 0x21, 0xEB, 0xFB, 0x25, 0xE0, 0x61, 0x63, 0x65, 0x76, 0x1F, 0x40, 0x2C,
		},
		"i4004/08-moltiplicazione-bcd.asm": {
			0xD0, 0xFD, 0x20, 0x00, 0x21, 0xD5, 0xE0, 0x20, 0x01, 0x21, 0xD2, 0xE0, 0x26, 0x0B, 0x20, 0x00,
			0x22, 0x10, 0x24, 0xD0, 0xF1, 0x23, 0xE9, 0x21, 0xEB, 0xFB, 0x23, 0xE0, 0x61, 0x63, 0x74, 0x15,
			0x77, 0x0E, 0x40, 0x22,
		},
		"i4004/09-divisione-bcd.asm": {
			0xD0, 0xFD, 0xD2, 0xB1, 0xD7, 0xB2, 0xD0, 0xB3, 0xFA, 0xF9, 0x91, 0x82, 0xFB, 0x12, 0x11, 0x40,
			0x15, 0xB2, 0x63, 0x40, 0x08, 0x20, 0x00, 0x21, 0xA3, 0xE0, 0x20, 0x01, 0x21, 0xA2, 0xE0, 0x40,
			0x1F,
		},
	}
	for name, want := range cases {
		if got := assembleFile(t, name); !bytes.Equal(got, want) {
			t.Errorf("%s =\n % X\natteso\n % X", name, got, want)
		}
	}
}
