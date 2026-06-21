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
		if st.Org != nil {
			if *st.Org < pc {
				return nil, fmt.Errorf("riga %d: .org 0x%03X precede la posizione corrente 0x%03X", st.Line, *st.Org, pc)
			}
			if *st.Org > 0xFFF {
				return nil, fmt.Errorf("riga %d: .org 0x%X fuori dallo spazio ROM (max 0xFFF)", st.Line, *st.Org)
			}
			pc = *st.Org
		}
		// La label viene registrata dopo l'eventuale .org, così "etichetta: .org N"
		// punta a N; per le altre righe punta alla posizione corrente.
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
		pc += len(st.Data)
		// Le costanti .equ entrano nella symbol table come le label (nome → valore),
		// quindi sono usabili anche prima della loro definizione (risolte in passata 2).
		if st.Equ != nil {
			if err := syms.Define(st.Equ.Name, st.Equ.Value); err != nil {
				return nil, fmt.Errorf("riga %d: %w", st.Line, err)
			}
		}
	}

	// Passata 2: codifica, risolvendo le label con la symbol table.
	code := make([]byte, 0, pc)
	pc = 0
	for _, st := range stmts {
		if st.Org != nil {
			for pc < *st.Org { // riempi il vuoto fino all'indirizzo con NOP (0x00)
				code = append(code, 0x00)
				pc++
			}
			continue
		}
		if len(st.Data) > 0 { // .byte: emette i byte letterali
			code = append(code, st.Data...)
			pc += len(st.Data)
			continue
		}
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
