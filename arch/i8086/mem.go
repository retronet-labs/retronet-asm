package i8086

import (
	"fmt"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
)

// operand e' un operando decodificato: registro, immediato o riferimento in
// memoria con il suo ModR/M (mod/rm) e gli eventuali byte di spiazzamento.
type operand struct {
	kind opClass
	reg  byte   // codice per kReg8/kReg16/kSreg
	mod  byte   // per kMem
	rm   byte   // per kMem
	disp []byte // per kMem: 0, 1 o 2 byte di spiazzamento
	text string // testo grezzo (per gli immediati)
}

func (o operand) isMem() bool { return o.kind == kMem }

// memForms mappa la combinazione base+indice all'indice rm dell'indirizzamento a
// 16 bit dell'8086.
var memForms = map[string]byte{
	"bx+si": 0, "si+bx": 0,
	"bx+di": 1, "di+bx": 1,
	"bp+si": 2, "si+bp": 2,
	"bp+di": 3, "di+bp": 3,
	"si": 4, "di": 5, "bp": 6, "bx": 7,
}

// parseOperand classifica un operando. Per la memoria calcola mod/rm e i byte di
// spiazzamento. La dimensione dello spiazzamento dipende dalla SINTASSI (un
// letterale numerico sceglie disp8/disp16 dal valore; un simbolo usa sempre
// disp16), cosi' Size ed Encode concordano anche senza risolvere le label.
func parseOperand(s string, resolve arch.Resolver) (operand, error) {
	t := strings.TrimSpace(s)
	if strings.HasPrefix(t, "[") {
		return parseMem(t, resolve)
	}
	k, code := classify(t)
	return operand{kind: k, reg: code, text: t}, nil
}

func parseMem(s string, resolve arch.Resolver) (operand, error) {
	inner := strings.TrimSuffix(strings.TrimPrefix(s, "["), "]")
	if inner == "" {
		return operand{}, fmt.Errorf("operando in memoria vuoto")
	}

	var regs []string
	dispText := ""
	hasDisp := false
	for _, t := range splitTerms(inner) {
		low := strings.ToLower(t.text)
		if low == "bx" || low == "bp" || low == "si" || low == "di" {
			if t.neg {
				return operand{}, fmt.Errorf("registro %q non puo' essere sottratto", t.text)
			}
			regs = append(regs, low)
			continue
		}
		if hasDisp {
			return operand{}, fmt.Errorf("troppi termini di spiazzamento in %q", s)
		}
		hasDisp = true
		if t.neg {
			dispText = "-" + t.text
		} else {
			dispText = t.text
		}
	}

	// Indirizzo diretto [disp]: mod=00, rm=110, sempre disp16.
	if len(regs) == 0 {
		v, err := evalDisp(dispText, resolve)
		if err != nil {
			return operand{}, err
		}
		return operand{kind: kMem, mod: 0, rm: 6, disp: word(v)}, nil
	}

	var rm byte
	var ok bool
	if len(regs) == 1 {
		rm, ok = memForms[regs[0]]
	} else if len(regs) == 2 {
		rm, ok = memForms[regs[0]+"+"+regs[1]]
	}
	if !ok {
		return operand{}, fmt.Errorf("indirizzamento non valido %q", inner)
	}

	if !hasDisp {
		if rm == 6 { // [bp] da solo richiede disp8=0 (mod=00 rm=110 e' l'indiretto)
			return operand{kind: kMem, mod: 1, rm: 6, disp: []byte{0}}, nil
		}
		return operand{kind: kMem, mod: 0, rm: rm}, nil
	}

	v, err := evalDisp(dispText, resolve)
	if err != nil {
		return operand{}, err
	}
	// Letterale numerico che entra in 8 bit -> disp8; altrimenti disp16.
	if isLiteral(dispText) && v >= -128 && v <= 127 {
		return operand{kind: kMem, mod: 1, rm: rm, disp: []byte{byte(v)}}, nil
	}
	return operand{kind: kMem, mod: 2, rm: rm, disp: word(v)}, nil
}

type term struct {
	text string
	neg  bool
}

func splitTerms(s string) []term {
	var terms []term
	start, neg := 0, false
	flush := func(end int) {
		if end > start {
			terms = append(terms, term{text: s[start:end], neg: neg})
		}
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '+' || s[i] == '-' {
			flush(i)
			neg = s[i] == '-'
			start = i + 1
		}
	}
	flush(len(s))
	return terms
}

// isLiteral indica se il testo e' un numero letterale (non un simbolo).
func isLiteral(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if s[0] == '-' || s[0] == '+' {
		s = s[1:]
	}
	return s != "" && s[0] >= '0' && s[0] <= '9'
}

// evalDisp valuta uno spiazzamento (numero letterale o simbolo). Con resolve=nil
// (fase di Size) un simbolo vale 0: la dimensione e' comunque disp16 perche'
// isLiteral e' falso.
func evalDisp(s string, resolve arch.Resolver) (int, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return 0, nil
	}
	if isLiteral(t) {
		return parseNum(t)
	}
	if resolve == nil {
		return 0, nil
	}
	v, ok := resolve(t)
	if !ok {
		return 0, fmt.Errorf("simbolo non definito: %q", t)
	}
	return v, nil
}

func word(v int) []byte { return []byte{byte(v), byte(v >> 8)} }

// modrmBytes compone il ModR/M (campo reg + r/m dell'operando) e i byte di
// spiazzamento, sia per un operando in memoria sia per un registro (mod=11).
func modrmBytes(regField byte, rm operand) []byte {
	if rm.kind != kMem {
		return []byte{0xC0 | regField<<3 | rm.reg}
	}
	return append([]byte{rm.mod<<6 | regField<<3 | rm.rm}, rm.disp...)
}
