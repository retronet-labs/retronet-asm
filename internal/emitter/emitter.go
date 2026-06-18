// Package emitter assembla gli statement in byte ROM con due passate:
// prima calcola gli indirizzi (e raccoglie le label), poi codifica.
package emitter

import (
	"fmt"

	"github.com/retronet-labs/retronet-asm/arch"
	"github.com/retronet-labs/retronet-asm/internal/parser"
	"github.com/retronet-labs/retronet-asm/internal/symbols"
)

// Assemble traduce gli statement in byte usando l'architettura a.
func Assemble(stmts []parser.Stmt, a arch.Arch) ([]byte, error) {
	syms := symbols.New()

	// Passata 1: assegna gli indirizzi e registra le label.
	pc := 0
	for _, st := range stmts {
		if st.Label != "" {
			if err := syms.Define(st.Label, pc); err != nil {
				return nil, fmt.Errorf("riga %d: %w", st.Line, err)
			}
		}
		if st.Instr != nil {
			sz, err := a.Size(*st.Instr)
			if err != nil {
				return nil, err
			}
			pc += sz
		}
	}

	// Passata 2: codifica, risolvendo le label con la symbol table.
	code := make([]byte, 0, pc)
	pc = 0
	for _, st := range stmts {
		if st.Instr == nil {
			continue // riga di sola label: nessun byte
		}
		b, err := a.Encode(*st.Instr, pc, syms.Lookup)
		if err != nil {
			return nil, err
		}
		code = append(code, b...)
		pc += len(b)
	}
	return code, nil
}
