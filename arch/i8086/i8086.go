// Package i8086 implementa l'architettura Intel 8086/8088 per retronet-asm.
//
// Copre le istruzioni in real mode con operandi registro, immediato, segmento,
// memoria ([base+indice+disp] e [disp] diretto) e label (salti/chiamate), oltre
// alle istruzioni stringa e di controllo dei flag. Gli operandi in memoria usano
// la codifica ModR/M a 16 bit; per le forme memoria-immediato ambigue serve uno
// specificatore di dimensione (byte / word).
package i8086

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
)

var reg8 = map[string]byte{"AL": 0, "CL": 1, "DL": 2, "BL": 3, "AH": 4, "CH": 5, "DH": 6, "BH": 7}
var reg16 = map[string]byte{"AX": 0, "CX": 1, "DX": 2, "BX": 3, "SP": 4, "BP": 5, "SI": 6, "DI": 7}
var sreg = map[string]byte{"ES": 0, "CS": 1, "SS": 2, "DS": 3}

var aluGroup = map[string]byte{"ADD": 0, "OR": 1, "ADC": 2, "SBB": 3, "AND": 4, "SUB": 5, "XOR": 6, "CMP": 7}
var shiftGroup = map[string]byte{"ROL": 0, "ROR": 1, "RCL": 2, "RCR": 3, "SHL": 4, "SAL": 4, "SHR": 5, "SAR": 7}
var unaryGroup = map[string]byte{"NOT": 2, "NEG": 3, "MUL": 4, "IMUL": 5, "DIV": 6, "IDIV": 7}

var jcc = map[string]byte{
	"JO": 0, "JNO": 1, "JB": 2, "JC": 2, "JNAE": 2, "JAE": 3, "JNB": 3, "JNC": 3,
	"JE": 4, "JZ": 4, "JNE": 5, "JNZ": 5, "JBE": 6, "JNA": 6, "JA": 7, "JNBE": 7,
	"JS": 8, "JNS": 9, "JP": 10, "JPE": 10, "JNP": 11, "JPO": 11,
	"JL": 12, "JNGE": 12, "JGE": 13, "JNL": 13, "JLE": 14, "JNG": 14, "JG": 15, "JNLE": 15,
}

var noOperand = map[string]byte{
	"NOP": 0x90, "HLT": 0xF4, "IRET": 0xCF, "INT3": 0xCC, "INTO": 0xCE,
	"CLC": 0xF8, "STC": 0xF9, "CLI": 0xFA, "STI": 0xFB, "CLD": 0xFC, "STD": 0xFD, "CMC": 0xF5,
	"CBW": 0x98, "CWD": 0x99, "SAHF": 0x9E, "LAHF": 0x9F, "PUSHF": 0x9C, "POPF": 0x9D,
	"XLAT": 0xD7, "WAIT": 0x9B, "DAA": 0x27, "DAS": 0x2F, "AAA": 0x37, "AAS": 0x3F,
	"MOVSB": 0xA4, "MOVSW": 0xA5, "CMPSB": 0xA6, "CMPSW": 0xA7,
	"STOSB": 0xAA, "STOSW": 0xAB, "LODSB": 0xAC, "LODSW": 0xAD, "SCASB": 0xAE, "SCASW": 0xAF,
	"REP": 0xF3, "REPE": 0xF3, "REPZ": 0xF3, "REPNE": 0xF2, "REPNZ": 0xF2, "LOCK": 0xF0,
}

var loopOps = map[string]byte{"LOOPNZ": 0xE0, "LOOPNE": 0xE0, "LOOPZ": 0xE1, "LOOPE": 0xE1, "LOOP": 0xE2, "JCXZ": 0xE3}

type I8086 struct{}

// New crea il backend Intel 8086/8088.
func New() arch.Arch { return I8086{} }

func (I8086) Name() string { return "i8086" }

type opClass int

const (
	kImm opClass = iota
	kReg8
	kReg16
	kSreg
	kMem
)

func classify(s string) (opClass, byte) {
	u := strings.ToUpper(strings.TrimSpace(s))
	if c, ok := reg8[u]; ok {
		return kReg8, c
	}
	if c, ok := reg16[u]; ok {
		return kReg16, c
	}
	if c, ok := sreg[u]; ok {
		return kSreg, c
	}
	return kImm, 0
}

func isReg(o operand) bool { return o.kind == kReg8 || o.kind == kReg16 }
func isRM(o operand) bool  { return isReg(o) || o.kind == kMem }
func bw(word bool) byte {
	if word {
		return 1
	}
	return 0
}

// Size delega a Encode con resolver nullo: i simboli valgono 0 ma la lunghezza
// (dimensione degli immediati e degli spiazzamenti) e' deterministica.
func (a I8086) Size(in arch.Instruction) (int, error) {
	b, err := a.Encode(in, 0, nil)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// Encode produce i byte dell'istruzione all'indirizzo pc.
func (a I8086) Encode(in arch.Instruction, pc int, resolve arch.Resolver) ([]byte, error) {
	m := in.Mnemonic
	ops := in.Operands

	if op, ok := noOperand[m]; ok {
		return []byte{op}, nil
	}
	if cc, ok := jcc[m]; ok {
		if err := mustArity(in, 1); err != nil {
			return nil, err
		}
		return encodeRel8(0x70|cc, ops[0], pc, 2, resolve, in.Line)
	}
	if op, ok := loopOps[m]; ok {
		if err := mustArity(in, 1); err != nil {
			return nil, err
		}
		return encodeRel8(op, ops[0], pc, 2, resolve, in.Line)
	}

	switch m {
	case "INT":
		if err := mustArity(in, 1); err != nil {
			return nil, err
		}
		v, err := parseByteValue(ops[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{0xCD, v}, nil
	case "JMP":
		if len(ops) == 2 && strings.EqualFold(ops[0], "SHORT") {
			return encodeRel8(0xEB, ops[1], pc, 2, resolve, in.Line)
		}
		if err := mustArity(in, 1); err != nil {
			return nil, err
		}
		return encodeRel16(0xE9, ops[0], pc, resolve, in.Line)
	case "CALL":
		if err := mustArity(in, 1); err != nil {
			return nil, err
		}
		return encodeRel16(0xE8, ops[0], pc, resolve, in.Line)
	case "RET", "RETF":
		op := byte(0xC3)
		if m == "RETF" {
			op = 0xCB
		}
		if len(ops) == 0 {
			return []byte{op}, nil
		}
		v, err := parseWordValue(ops[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{op - 1, byte(v), byte(v >> 8)}, nil
	case "AAM", "AAD":
		base := byte(0x0A)
		if len(ops) == 1 {
			v, err := parseByteValue(ops[0], resolve)
			if err != nil {
				return nil, wrap(in.Line, err)
			}
			base = v
		}
		if m == "AAM" {
			return []byte{0xD4, base}, nil
		}
		return []byte{0xD5, base}, nil
	case "MOV":
		return encodeMOV(in, resolve)
	case "LEA":
		return encodeLEA(in, resolve)
	case "XCHG":
		return encodeXCHG(in, resolve)
	case "TEST":
		return encodeALU(in, true, resolve)
	case "INC", "DEC":
		return encodeIncDec(in, resolve)
	case "PUSH", "POP":
		return encodePushPop(in, resolve)
	case "IN", "OUT":
		return encodeInOut(in, resolve)
	}
	if _, ok := aluGroup[m]; ok {
		return encodeALU(in, false, resolve)
	}
	if g, ok := shiftGroup[m]; ok {
		return encodeShift(in, g, resolve)
	}
	if g, ok := unaryGroup[m]; ok {
		return encodeUnary(in, g, resolve)
	}
	return nil, fmt.Errorf("riga %d: mnemonico sconosciuto %q", in.Line, m)
}

// stripSize estrae un eventuale specificatore "byte"/"word" iniziale, restituendo
// la larghezza forzata (0 nessuna, 8, 16) e gli operandi rimanenti.
func stripSize(ops []string) (int, []string) {
	if len(ops) > 0 {
		switch strings.ToLower(ops[0]) {
		case "byte":
			return 8, ops[1:]
		case "word":
			return 16, ops[1:]
		}
	}
	return 0, ops
}

func encodeMOV(in arch.Instruction, resolve arch.Resolver) ([]byte, error) {
	forced, ops := stripSize(in.Operands)
	if len(ops) != 2 {
		return nil, arity2(in)
	}
	a, err := parseOperand(ops[0], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	b, err := parseOperand(ops[1], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}

	switch {
	case a.kind == kSreg: // MOV sreg, r/m16
		return append([]byte{0x8E}, modrmBytes(a.reg, b)...), nil
	case b.kind == kSreg: // MOV r/m16, sreg
		return append([]byte{0x8C}, modrmBytes(b.reg, a)...), nil
	case isReg(a) && b.kind == kImm: // MOV reg, imm
		if a.kind == kReg8 {
			v, err := parseByteValue(b.text, resolve)
			if err != nil {
				return nil, wrap(in.Line, err)
			}
			return []byte{0xB0 | a.reg, v}, nil
		}
		v, err := parseWordValue(b.text, resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{0xB8 | a.reg, byte(v), byte(v >> 8)}, nil
	case a.kind == kMem && b.kind == kImm: // MOV mem, imm (serve byte/word)
		w, err := memWidth(in, forced)
		if err != nil {
			return nil, err
		}
		return appendImm(append([]byte{0xC6 | bw(w)}, modrmBytes(0, a)...), b.text, w, resolve, in.Line)
	case isReg(b) && isRM(a): // MOV r/m, r
		return append([]byte{0x88 | bw(b.kind == kReg16)}, modrmBytes(b.reg, a)...), nil
	case isReg(a) && b.kind == kMem: // MOV r, m
		return append([]byte{0x8A | bw(a.kind == kReg16)}, modrmBytes(a.reg, b)...), nil
	}
	return nil, fmt.Errorf("riga %d: forma di MOV non supportata", in.Line)
}

func encodeLEA(in arch.Instruction, resolve arch.Resolver) ([]byte, error) {
	if err := mustArity(in, 2); err != nil {
		return nil, err
	}
	a, _ := parseOperand(in.Operands[0], resolve)
	b, err := parseOperand(in.Operands[1], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	if a.kind != kReg16 || b.kind != kMem {
		return nil, fmt.Errorf("riga %d: LEA vuole un registro a 16 bit e un operando in memoria", in.Line)
	}
	return append([]byte{0x8D}, modrmBytes(a.reg, b)...), nil
}

func encodeXCHG(in arch.Instruction, resolve arch.Resolver) ([]byte, error) {
	if err := mustArity(in, 2); err != nil {
		return nil, err
	}
	a, _ := parseOperand(in.Operands[0], resolve)
	b, _ := parseOperand(in.Operands[1], resolve)
	if isReg(b) && isRM(a) {
		return append([]byte{0x86 | bw(b.kind == kReg16)}, modrmBytes(b.reg, a)...), nil
	}
	if isReg(a) && b.kind == kMem {
		return append([]byte{0x86 | bw(a.kind == kReg16)}, modrmBytes(a.reg, b)...), nil
	}
	return nil, fmt.Errorf("riga %d: XCHG richiede due registri o registro/memoria", in.Line)
}

func encodeALU(in arch.Instruction, isTest bool, resolve arch.Resolver) ([]byte, error) {
	forced, ops := stripSize(in.Operands)
	if len(ops) != 2 {
		return nil, arity2(in)
	}
	a, err := parseOperand(ops[0], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	b, err := parseOperand(ops[1], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	g := aluGroup[in.Mnemonic]

	switch {
	case isReg(b) && isRM(a): // r/m, r
		w := b.kind == kReg16
		if isTest {
			return append([]byte{0x84 | bw(w)}, modrmBytes(b.reg, a)...), nil
		}
		return append([]byte{g<<3 | bw(w)}, modrmBytes(b.reg, a)...), nil
	case isReg(a) && b.kind == kMem: // r, r/m
		w := a.kind == kReg16
		if isTest {
			return append([]byte{0x84 | bw(w)}, modrmBytes(a.reg, b)...), nil
		}
		return append([]byte{g<<3 | 0x02 | bw(w)}, modrmBytes(a.reg, b)...), nil
	case b.kind == kImm: // r/m, imm
		// forma corta per l'accumulatore (AL/AX)
		if isReg(a) && a.reg == 0 {
			w := a.kind == kReg16
			op := byte(0xA8) // TEST AL/AX, imm
			if !isTest {
				op = g<<3 | 0x04
			}
			return appendImm([]byte{op | bw(w)}, b.text, w, resolve, in.Line)
		}
		w := a.kind == kReg16
		if a.kind == kMem {
			var err error
			if w, err = memWidth(in, forced); err != nil {
				return nil, err
			}
		}
		base, ext := byte(0x80), g // ALU: 80/81 /g
		if isTest {
			base, ext = 0xF6, 0 // TEST: F6/F7 /0
		}
		return appendImm(append([]byte{base | bw(w)}, modrmBytes(ext, a)...), b.text, w, resolve, in.Line)
	}
	return nil, fmt.Errorf("riga %d: forma di %s non supportata", in.Line, in.Mnemonic)
}

func encodeIncDec(in arch.Instruction, resolve arch.Resolver) ([]byte, error) {
	forced, ops := stripSize(in.Operands)
	if len(ops) != 1 {
		return nil, arity1(in)
	}
	a, err := parseOperand(ops[0], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	ext := byte(0) // INC
	if in.Mnemonic == "DEC" {
		ext = 1
	}
	if a.kind == kReg16 {
		base := byte(0x40)
		if in.Mnemonic == "DEC" {
			base = 0x48
		}
		return []byte{base | a.reg}, nil
	}
	w := a.kind == kReg16
	if a.kind == kMem {
		if w, err = memWidth(in, forced); err != nil {
			return nil, err
		}
	}
	return append([]byte{0xFE | bw(w)}, modrmBytes(ext, a)...), nil
}

func encodeUnary(in arch.Instruction, g byte, resolve arch.Resolver) ([]byte, error) {
	forced, ops := stripSize(in.Operands)
	if len(ops) != 1 {
		return nil, arity1(in)
	}
	a, err := parseOperand(ops[0], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	w := a.kind == kReg16
	if a.kind == kMem {
		if w, err = memWidth(in, forced); err != nil {
			return nil, err
		}
	} else if !isReg(a) {
		return nil, fmt.Errorf("riga %d: %s richiede registro o memoria", in.Line, in.Mnemonic)
	}
	return append([]byte{0xF6 | bw(w)}, modrmBytes(g, a)...), nil
}

func encodeShift(in arch.Instruction, g byte, resolve arch.Resolver) ([]byte, error) {
	forced, ops := stripSize(in.Operands)
	if len(ops) != 2 {
		return nil, arity2(in)
	}
	a, err := parseOperand(ops[0], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	count := strings.ToUpper(strings.TrimSpace(ops[1]))
	var op byte
	switch count {
	case "1":
		op = 0xD0
	case "CL":
		op = 0xD2
	default:
		return nil, fmt.Errorf("riga %d: %s accetta solo 1 oppure CL", in.Line, in.Mnemonic)
	}
	w := a.kind == kReg16
	if a.kind == kMem {
		if w, err = memWidth(in, forced); err != nil {
			return nil, err
		}
	} else if !isReg(a) {
		return nil, fmt.Errorf("riga %d: %s richiede registro o memoria", in.Line, in.Mnemonic)
	}
	return append([]byte{op | bw(w)}, modrmBytes(g, a)...), nil
}

func encodePushPop(in arch.Instruction, resolve arch.Resolver) ([]byte, error) {
	_, ops := stripSize(in.Operands) // PUSH/POP sono sempre a 16 bit: "word" e' opzionale
	if len(ops) != 1 {
		return nil, arity1(in)
	}
	a, err := parseOperand(ops[0], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	switch {
	case a.kind == kReg16:
		base := byte(0x50)
		if in.Mnemonic == "POP" {
			base = 0x58
		}
		return []byte{base | a.reg}, nil
	case a.kind == kSreg:
		if in.Mnemonic == "PUSH" {
			return []byte{0x06 | a.reg<<3}, nil
		}
		return []byte{0x07 | a.reg<<3}, nil
	case a.kind == kMem:
		if in.Mnemonic == "PUSH" {
			return append([]byte{0xFF}, modrmBytes(6, a)...), nil
		}
		return append([]byte{0x8F}, modrmBytes(0, a)...), nil
	}
	return nil, fmt.Errorf("riga %d: %s richiede registro a 16 bit, segmento o memoria", in.Line, in.Mnemonic)
}

func encodeInOut(in arch.Instruction, resolve arch.Resolver) ([]byte, error) {
	if err := mustArity(in, 2); err != nil {
		return nil, err
	}
	if in.Mnemonic == "IN" {
		acc, _ := classify(in.Operands[0])
		w := acc == kReg16
		if strings.EqualFold(strings.TrimSpace(in.Operands[1]), "DX") {
			return []byte{0xEC | bw(w)}, nil
		}
		v, err := parseByteValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{0xE4 | bw(w), v}, nil
	}
	acc, _ := classify(in.Operands[1])
	w := acc == kReg16
	if strings.EqualFold(strings.TrimSpace(in.Operands[0]), "DX") {
		return []byte{0xEE | bw(w)}, nil
	}
	v, err := parseByteValue(in.Operands[0], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	return []byte{0xE6 | bw(w), v}, nil
}

// memWidth determina la larghezza per un operando in memoria con immediato:
// richiede lo specificatore byte/word.
func memWidth(in arch.Instruction, forced int) (bool, error) {
	switch forced {
	case 8:
		return false, nil
	case 16:
		return true, nil
	}
	return false, fmt.Errorf("riga %d: %s su memoria richiede 'byte' o 'word'", in.Line, in.Mnemonic)
}

// appendImm aggiunge l'immediato (1 o 2 byte) a prefix.
func appendImm(prefix []byte, text string, word bool, resolve arch.Resolver, line int) ([]byte, error) {
	if word {
		v, err := parseWordValue(text, resolve)
		if err != nil {
			return nil, wrap(line, err)
		}
		return append(prefix, byte(v), byte(v>>8)), nil
	}
	v, err := parseByteValue(text, resolve)
	if err != nil {
		return nil, wrap(line, err)
	}
	return append(prefix, v), nil
}

func encodeRel8(op byte, target string, pc, size int, resolve arch.Resolver, line int) ([]byte, error) {
	dest, err := parseValue(target, resolve)
	if err != nil {
		return nil, wrap(line, err)
	}
	rel := dest - (pc + size)
	if resolve != nil && (rel < -128 || rel > 127) {
		return nil, fmt.Errorf("riga %d: salto fuori portata per rel8 (%d)", line, rel)
	}
	return []byte{op, byte(rel)}, nil
}

func encodeRel16(op byte, target string, pc int, resolve arch.Resolver, line int) ([]byte, error) {
	dest, err := parseValue(target, resolve)
	if err != nil {
		return nil, wrap(line, err)
	}
	rel := dest - (pc + 3)
	return []byte{op, byte(rel), byte(rel >> 8)}, nil
}

func arity1(in arch.Instruction) error { return mustArity(in, 1) }
func arity2(in arch.Instruction) error { return mustArity(in, 2) }

func mustArity(in arch.Instruction, want int) error {
	if len(in.Operands) != want {
		return fmt.Errorf("riga %d: %s vuole %d operandi, trovati %d", in.Line, in.Mnemonic, want, len(in.Operands))
	}
	return nil
}

func parseByteValue(s string, resolve arch.Resolver) (byte, error) {
	v, err := parseValue(s, resolve)
	if err != nil {
		return 0, err
	}
	if v < -128 || v > 0xFF {
		return 0, fmt.Errorf("valore %d fuori range byte", v)
	}
	return byte(v), nil
}

func parseWordValue(s string, resolve arch.Resolver) (uint16, error) {
	v, err := parseValue(s, resolve)
	if err != nil {
		return 0, err
	}
	if v < -32768 || v > 0xFFFF {
		return 0, fmt.Errorf("valore 0x%X fuori range 16 bit", v)
	}
	return uint16(v), nil
}

// parseValue interpreta un numero oppure risolve un simbolo. Con resolve=nil
// (fase di Size) i simboli valgono 0: la lunghezza non dipende dal valore.
func parseValue(s string, resolve arch.Resolver) (int, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return 0, fmt.Errorf("operando vuoto")
	}
	if t[0] == '-' || (t[0] >= '0' && t[0] <= '9') {
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

func parseNum(s string) (int, error) {
	t := strings.TrimSpace(s)
	neg := false
	if strings.HasPrefix(t, "-") {
		neg, t = true, t[1:]
	} else if strings.HasPrefix(t, "+") {
		t = t[1:]
	}
	var n int64
	var err error
	if strings.HasPrefix(strings.ToLower(t), "0x") {
		n, err = strconv.ParseInt(t[2:], 16, 32)
	} else {
		n, err = strconv.ParseInt(t, 10, 32)
	}
	if err != nil {
		return 0, fmt.Errorf("numero non valido %q", s)
	}
	if neg {
		n = -n
	}
	return int(n), nil
}

func wrap(line int, err error) error { return fmt.Errorf("riga %d: %w", line, err) }
