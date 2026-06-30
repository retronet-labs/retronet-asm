// Package emitter assembla gli statement in byte ROM con due passate:
// prima calcola gli indirizzi (e raccoglie le label), poi codifica.
package emitter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
	"github.com/retronet-labs/retronet-asm/internal/parser"
	"github.com/retronet-labs/retronet-asm/internal/symbols"
)

// Assemble traduce gli statement in byte usando l'architettura a.
func Assemble(stmts []parser.Stmt, a arch.Arch) ([]byte, error) {
	syms := symbols.New()
	maxAddr := maxAddress(a)

	// Passata 1: assegna gli indirizzi e registra le label.
	base := 0
	pc := 0
	for _, st := range stmts {
		if st.OrgBase != nil {
			if pc != base {
				return nil, fmt.Errorf("riga %d: .orgbase deve precedere codice e dati", st.Line)
			}
			if *st.OrgBase < 0 || *st.OrgBase > maxAddr {
				return nil, fmt.Errorf("riga %d: .orgbase 0x%X fuori range (max 0x%X)", st.Line, *st.OrgBase, maxAddr)
			}
			base = *st.OrgBase
			pc = base
		}
		if st.Org != nil {
			if *st.Org < pc {
				return nil, fmt.Errorf("riga %d: .org 0x%03X precede la posizione corrente 0x%03X", st.Line, *st.Org, pc)
			}
			if *st.Org > maxAddr {
				return nil, fmt.Errorf("riga %d: .org 0x%X fuori dallo spazio ROM (max 0x%X)", st.Line, *st.Org, maxAddr)
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
		pc += 2 * len(st.Words)
		// Le costanti .equ entrano nella symbol table come le label (nome → valore),
		// quindi sono usabili anche prima della loro definizione (risolte in passata 2).
		if st.Equ != nil {
			if err := syms.Define(st.Equ.Name, st.Equ.Value); err != nil {
				return nil, fmt.Errorf("riga %d: %w", st.Line, err)
			}
		}
	}

	// Passata 2: codifica, risolvendo le label con la symbol table.
	code := make([]byte, 0, pc-base)
	base = 0
	pc = 0
	for _, st := range stmts {
		if st.OrgBase != nil {
			if len(code) > 0 || pc != base {
				return nil, fmt.Errorf("riga %d: .orgbase deve precedere codice e dati", st.Line)
			}
			base = *st.OrgBase
			pc = base
			continue
		}
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
		if len(st.Words) > 0 {
			for _, w := range st.Words {
				v, err := parseWord(w, syms.Lookup)
				if err != nil {
					return nil, fmt.Errorf("riga %d: %w", st.Line, err)
				}
				code = append(code, byte(v), byte(v>>8))
				pc += 2
			}
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

func parseWord(s string, resolve arch.Resolver) (uint16, error) {
	v, err := parseValue(s, resolve)
	if err != nil {
		return 0, err
	}
	if v < 0 || v > 0xFFFF {
		return 0, fmt.Errorf(".word 0x%X fuori range 16 bit", v)
	}
	return uint16(v), nil
}

func parseValue(s string, resolve arch.Resolver) (int, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return 0, fmt.Errorf("operando vuoto")
	}
	if isNumber(t) {
		return parseNum(t)
	}
	if resolve == nil {
		return 0, fmt.Errorf("simbolo %q non risolvibile", t)
	}
	v, ok := resolve(t)
	if !ok {
		return 0, fmt.Errorf("simbolo non definito: %q", t)
	}
	return v, nil
}

func parseNum(s string) (int, error) {
	t := strings.TrimSpace(s)
	var n int64
	var err error
	switch {
	case strings.HasPrefix(strings.ToLower(t), "0x"):
		n, err = strconv.ParseInt(t[2:], 16, 32)
	case strings.HasPrefix(t, "$"):
		n, err = strconv.ParseInt(t[1:], 16, 32)
	case strings.HasPrefix(t, "%"):
		n, err = strconv.ParseInt(t[1:], 2, 32)
	default:
		n, err = strconv.ParseInt(t, 10, 32)
	}
	if err != nil {
		return 0, fmt.Errorf("numero non valido %q", s)
	}
	return int(n), nil
}

func isNumber(s string) bool {
	if s == "" {
		return false
	}
	if s[0] == '$' || s[0] == '%' {
		return true
	}
	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
	}
	return s != "" && s[0] >= '0' && s[0] <= '9'
}

func maxAddress(a arch.Arch) int {
	if a.Name() == "i4004" {
		return 0x0FFF
	}
	return 0xFFFF
}
