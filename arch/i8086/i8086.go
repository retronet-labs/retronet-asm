// Package i8086 implementa l'architettura Intel 8086/8088 per retronet-asm.
//
// Copre le istruzioni in real mode con operandi di tipo registro, immediato,
// segmento, label (salti/chiamate) e le istruzioni stringa e di controllo dei
// flag: abbastanza per scrivere boot sector e programmi a registri. Gli operandi
// in memoria con parentesi (es. [bx+si]) NON sono ancora supportati: userebbero
// la decodifica completa del ModR/M e una sintassi col lexer esteso.
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

// aluGroup mappa i mnemonici aritmetico-logici al loro indice di gruppo (0-7).
var aluGroup = map[string]byte{"ADD": 0, "OR": 1, "ADC": 2, "SBB": 3, "AND": 4, "SUB": 5, "XOR": 6, "CMP": 7}

// shiftGroup mappa shift/rotate all'estensione /n del gruppo D0-D3.
var shiftGroup = map[string]byte{"ROL": 0, "ROR": 1, "RCL": 2, "RCR": 3, "SHL": 4, "SAL": 4, "SHR": 5, "SAR": 7}

// unaryGroup mappa NOT/NEG/MUL/IMUL/DIV/IDIV all'estensione /n del gruppo F6/F7.
var unaryGroup = map[string]byte{"NOT": 2, "NEG": 3, "MUL": 4, "IMUL": 5, "DIV": 6, "IDIV": 7}

// jcc mappa i salti condizionati al codice di condizione (opcode 0x70|cc).
var jcc = map[string]byte{
	"JO": 0, "JNO": 1, "JB": 2, "JC": 2, "JNAE": 2, "JAE": 3, "JNB": 3, "JNC": 3,
	"JE": 4, "JZ": 4, "JNE": 5, "JNZ": 5, "JBE": 6, "JNA": 6, "JA": 7, "JNBE": 7,
	"JS": 8, "JNS": 9, "JP": 10, "JPE": 10, "JNP": 11, "JPO": 11,
	"JL": 12, "JNGE": 12, "JGE": 13, "JNL": 13, "JLE": 14, "JNG": 14, "JG": 15, "JNLE": 15,
}

// noOperand sono i mnemonici a 1 byte senza operandi.
var noOperand = map[string]byte{
	"NOP": 0x90, "HLT": 0xF4, "IRET": 0xCF, "INT3": 0xCC, "INTO": 0xCE,
	"CLC": 0xF8, "STC": 0xF9, "CLI": 0xFA, "STI": 0xFB, "CLD": 0xFC, "STD": 0xFD, "CMC": 0xF5,
	"CBW": 0x98, "CWD": 0x99, "SAHF": 0x9E, "LAHF": 0x9F, "PUSHF": 0x9C, "POPF": 0x9D,
	"XLAT": 0xD7, "WAIT": 0x9B, "DAA": 0x27, "DAS": 0x2F, "AAA": 0x37, "AAS": 0x3F,
	"MOVSB": 0xA4, "MOVSW": 0xA5, "CMPSB": 0xA6, "CMPSW": 0xA7,
	"STOSB": 0xAA, "STOSW": 0xAB, "LODSB": 0xAC, "LODSW": 0xAD, "SCASB": 0xAE, "SCASW": 0xAF,
	"REP": 0xF3, "REPE": 0xF3, "REPZ": 0xF3, "REPNE": 0xF2, "REPNZ": 0xF2, "LOCK": 0xF0,
}

// loopOps sono i salti basati su CX, tutti con rel8 (2 byte).
var loopOps = map[string]byte{"LOOPNZ": 0xE0, "LOOPNE": 0xE0, "LOOPZ": 0xE1, "LOOPE": 0xE1, "LOOP": 0xE2, "JCXZ": 0xE3}

type I8086 struct{}

// New crea il backend Intel 8086/8088.
func New() arch.Arch { return I8086{} }

func (I8086) Name() string { return "i8086" }

// --- classificazione degli operandi ---

type opClass int

const (
	kImm opClass = iota
	kReg8
	kReg16
	kSreg
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

func mod11(reg, rm byte) byte { return 0xC0 | reg<<3 | rm }

// --- Size: lunghezza in byte senza risolvere le label ---

func (a I8086) Size(in arch.Instruction) (int, error) {
	m := in.Mnemonic
	ops := in.Operands
	if _, ok := noOperand[m]; ok {
		return checkArity(in, 0, 1)
	}
	if _, ok := jcc[m]; ok {
		return checkArity(in, 1, 2)
	}
	if _, ok := loopOps[m]; ok {
		return checkArity(in, 1, 2)
	}
	switch m {
	case "INT":
		return checkArity(in, 1, 2)
	case "CALL":
		return checkArity(in, 1, 3)
	case "JMP":
		if len(ops) == 2 && strings.EqualFold(ops[0], "SHORT") {
			return 2, nil
		}
		return checkArity(in, 1, 3)
	case "RET", "RETF":
		if len(ops) == 0 {
			return 1, nil
		}
		return checkArity(in, 1, 3)
	case "AAM", "AAD":
		return 2, nil
	case "MOV":
		return sizeMOV(in)
	case "XCHG":
		return checkArity(in, 2, 2) // forma registro-registro
	case "TEST":
		return sizeTestALU(in, true)
	case "INC", "DEC":
		return sizeIncDec(in)
	case "PUSH", "POP":
		return 1, mustArity(in, 1)
	case "IN", "OUT":
		return sizeInOut(in)
	}
	if _, ok := aluGroup[m]; ok {
		return sizeTestALU(in, false)
	}
	if _, ok := shiftGroup[m]; ok {
		return checkArity(in, 2, 2)
	}
	if _, ok := unaryGroup[m]; ok {
		return checkArity(in, 1, 2)
	}
	return 0, fmt.Errorf("riga %d: mnemonico sconosciuto %q", in.Line, m)
}

func sizeMOV(in arch.Instruction) (int, error) {
	if err := mustArity(in, 2); err != nil {
		return 0, err
	}
	d, _ := classify(in.Operands[0])
	s, _ := classify(in.Operands[1])
	switch {
	case d == kSreg || s == kSreg:
		return 2, nil // 8E / 8C + ModRM
	case d == kReg8 && s == kImm:
		return 2, nil // B0+r ib
	case d == kReg16 && s == kImm:
		return 3, nil // B8+r iw
	default:
		return 2, nil // reg,reg: opcode + ModRM
	}
}

func sizeTestALU(in arch.Instruction, isTest bool) (int, error) {
	if err := mustArity(in, 2); err != nil {
		return 0, err
	}
	d, dc := classify(in.Operands[0])
	s, _ := classify(in.Operands[1])
	if s != kImm { // reg,reg
		return 2, nil
	}
	// destinazione, immediato
	switch {
	case d == kReg8 && dc == 0: // AL
		return 2, nil
	case d == kReg16 && dc == 0: // AX
		return 3, nil
	case isTest && d == kReg8:
		return 3, nil // F6 /0 modrm ib
	case isTest && d == kReg16:
		return 4, nil // F7 /0 modrm iw
	case d == kReg8:
		return 3, nil // 80 /g modrm ib
	default:
		return 4, nil // 81 /g modrm iw
	}
}

func sizeIncDec(in arch.Instruction) (int, error) {
	if err := mustArity(in, 1); err != nil {
		return 0, err
	}
	if c, _ := classify(in.Operands[0]); c == kReg16 {
		return 1, nil
	}
	return 2, nil // reg8: FE /n
}

func sizeInOut(in arch.Instruction) (int, error) {
	if err := mustArity(in, 2); err != nil {
		return 0, err
	}
	// la porta e' DX (1 byte) oppure un immediato (2 byte)
	port := in.Operands[0]
	if in.Mnemonic == "IN" {
		port = in.Operands[1]
	}
	if strings.EqualFold(strings.TrimSpace(port), "DX") {
		return 1, nil
	}
	return 2, nil
}

// --- Encode ---

func (a I8086) Encode(in arch.Instruction, pc int, resolve arch.Resolver) ([]byte, error) {
	m := in.Mnemonic
	ops := in.Operands

	if op, ok := noOperand[m]; ok {
		return []byte{op}, nil
	}
	if cc, ok := jcc[m]; ok {
		return encodeRel8(0x70|cc, ops[0], pc, 2, resolve, in.Line)
	}
	if op, ok := loopOps[m]; ok {
		return encodeRel8(op, ops[0], pc, 2, resolve, in.Line)
	}

	switch m {
	case "INT":
		v, err := parseByteValue(ops[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{0xCD, v}, nil
	case "JMP":
		if len(ops) == 2 && strings.EqualFold(ops[0], "SHORT") {
			return encodeRel8(0xEB, ops[1], pc, 2, resolve, in.Line)
		}
		return encodeRel16(0xE9, ops[0], pc, resolve, in.Line)
	case "CALL":
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
		return []byte{op - 1, byte(v), byte(v >> 8)}, nil // C2/CA imm16
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
	case "XCHG":
		return encodeXCHG(in)
	case "TEST":
		return encodeTestALU(in, true, resolve)
	case "INC", "DEC":
		return encodeIncDec(in)
	case "PUSH", "POP":
		return encodePushPop(in)
	case "IN", "OUT":
		return encodeInOut(in, resolve)
	}
	if _, ok := aluGroup[m]; ok {
		return encodeTestALU(in, false, resolve)
	}
	if g, ok := shiftGroup[m]; ok {
		return encodeShift(in, g)
	}
	if g, ok := unaryGroup[m]; ok {
		return encodeUnary(in, g)
	}
	return nil, fmt.Errorf("riga %d: mnemonico sconosciuto %q", in.Line, m)
}

func encodeMOV(in arch.Instruction, resolve arch.Resolver) ([]byte, error) {
	d, dc := classify(in.Operands[0])
	s, sc := classify(in.Operands[1])
	switch {
	case d == kSreg: // MOV sreg, r16
		if s != kReg16 {
			return nil, fmt.Errorf("riga %d: MOV su registro di segmento richiede un registro a 16 bit", in.Line)
		}
		return []byte{0x8E, mod11(dc, sc)}, nil
	case s == kSreg: // MOV r16, sreg
		if d != kReg16 {
			return nil, fmt.Errorf("riga %d: MOV da registro di segmento richiede un registro a 16 bit", in.Line)
		}
		return []byte{0x8C, mod11(sc, dc)}, nil
	case d == kReg8 && s == kReg8:
		return []byte{0x88, mod11(sc, dc)}, nil
	case d == kReg16 && s == kReg16:
		return []byte{0x89, mod11(sc, dc)}, nil
	case d == kReg8 && s == kImm:
		v, err := parseByteValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{0xB0 | dc, v}, nil
	case d == kReg16 && s == kImm:
		v, err := parseWordValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{0xB8 | dc, byte(v), byte(v >> 8)}, nil
	}
	return nil, fmt.Errorf("riga %d: forma di MOV non supportata", in.Line)
}

func encodeXCHG(in arch.Instruction) ([]byte, error) {
	d, dc := classify(in.Operands[0])
	s, sc := classify(in.Operands[1])
	if d == kReg16 && s == kReg16 {
		return []byte{0x87, mod11(sc, dc)}, nil
	}
	if d == kReg8 && s == kReg8 {
		return []byte{0x86, mod11(sc, dc)}, nil
	}
	return nil, fmt.Errorf("riga %d: XCHG richiede due registri della stessa larghezza", in.Line)
}

func encodeTestALU(in arch.Instruction, isTest bool, resolve arch.Resolver) ([]byte, error) {
	d, dc := classify(in.Operands[0])
	s, sc := classify(in.Operands[1])
	w := d == kReg16

	if s != kImm { // reg, reg
		if d != s {
			return nil, fmt.Errorf("riga %d: operandi di larghezza diversa", in.Line)
		}
		if isTest {
			op := byte(0x84)
			if w {
				op = 0x85
			}
			return []byte{op, mod11(sc, dc)}, nil
		}
		g := aluGroup[in.Mnemonic]
		op := g << 3 // forma r/m,r (00)
		if w {
			op |= 1
		}
		return []byte{op, mod11(sc, dc)}, nil
	}

	// destinazione, immediato
	if isTest {
		if dc == 0 { // TEST AL/AX, imm (A8/A9)
			if w {
				v, err := parseWordValue(in.Operands[1], resolve)
				if err != nil {
					return nil, wrap(in.Line, err)
				}
				return []byte{0xA9, byte(v), byte(v >> 8)}, nil
			}
			v, err := parseByteValue(in.Operands[1], resolve)
			if err != nil {
				return nil, wrap(in.Line, err)
			}
			return []byte{0xA8, v}, nil
		}
		return encodeGroupImm(in, 0xF6, 0, w, dc, resolve) // F6/F7 /0
	}

	g := aluGroup[in.Mnemonic]
	if dc == 0 { // ALU AL/AX, imm (forma accumulatore)
		if w {
			v, err := parseWordValue(in.Operands[1], resolve)
			if err != nil {
				return nil, wrap(in.Line, err)
			}
			return []byte{g<<3 | 0x05, byte(v), byte(v >> 8)}, nil
		}
		v, err := parseByteValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{g<<3 | 0x04, v}, nil
	}
	return encodeGroupImm(in, 0x80, g, w, dc, resolve) // 80/81 /g
}

// encodeGroupImm codifica la forma "gruppo + immediato" (80/81 per ALU, F6/F7 per
// TEST) con ModRM mod=11.
func encodeGroupImm(in arch.Instruction, baseOp, g byte, w bool, rm byte, resolve arch.Resolver) ([]byte, error) {
	op := baseOp
	if w {
		op |= 1
	}
	modrm := mod11(g, rm)
	if w {
		v, err := parseWordValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{op, modrm, byte(v), byte(v >> 8)}, nil
	}
	v, err := parseByteValue(in.Operands[1], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	return []byte{op, modrm, v}, nil
}

func encodeIncDec(in arch.Instruction) ([]byte, error) {
	c, code := classify(in.Operands[0])
	if c == kReg16 {
		base := byte(0x40) // INC
		if in.Mnemonic == "DEC" {
			base = 0x48
		}
		return []byte{base | code}, nil
	}
	if c == kReg8 {
		ext := byte(0) // INC
		if in.Mnemonic == "DEC" {
			ext = 1
		}
		return []byte{0xFE, mod11(ext, code)}, nil
	}
	return nil, fmt.Errorf("riga %d: %s richiede un registro", in.Line, in.Mnemonic)
}

func encodePushPop(in arch.Instruction) ([]byte, error) {
	c, code := classify(in.Operands[0])
	switch c {
	case kReg16:
		base := byte(0x50) // PUSH
		if in.Mnemonic == "POP" {
			base = 0x58
		}
		return []byte{base | code}, nil
	case kSreg:
		if in.Mnemonic == "PUSH" {
			return []byte{0x06 | code<<3}, nil
		}
		return []byte{0x07 | code<<3}, nil
	}
	return nil, fmt.Errorf("riga %d: %s richiede un registro a 16 bit o di segmento", in.Line, in.Mnemonic)
}

func encodeShift(in arch.Instruction, g byte) ([]byte, error) {
	c, code := classify(in.Operands[0])
	w := c == kReg16
	if c != kReg8 && c != kReg16 {
		return nil, fmt.Errorf("riga %d: %s richiede un registro", in.Line, in.Mnemonic)
	}
	count := strings.ToUpper(strings.TrimSpace(in.Operands[1]))
	var op byte
	switch count {
	case "1":
		op = 0xD0
	case "CL":
		op = 0xD2
	default:
		return nil, fmt.Errorf("riga %d: %s accetta solo 1 oppure CL come conteggio", in.Line, in.Mnemonic)
	}
	if w {
		op |= 1
	}
	return []byte{op, mod11(g, code)}, nil
}

func encodeUnary(in arch.Instruction, g byte) ([]byte, error) {
	c, code := classify(in.Operands[0])
	if c != kReg8 && c != kReg16 {
		return nil, fmt.Errorf("riga %d: %s richiede un registro", in.Line, in.Mnemonic)
	}
	op := byte(0xF6)
	if c == kReg16 {
		op = 0xF7
	}
	return []byte{op, mod11(g, code)}, nil
}

func encodeInOut(in arch.Instruction, resolve arch.Resolver) ([]byte, error) {
	if in.Mnemonic == "IN" {
		acc, _ := classify(in.Operands[0])
		w := acc == kReg16
		if strings.EqualFold(strings.TrimSpace(in.Operands[1]), "DX") {
			return []byte{inOutOp(0xEC, w)}, nil
		}
		v, err := parseByteValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{inOutOp(0xE4, w), v}, nil
	}
	// OUT
	acc, _ := classify(in.Operands[1])
	w := acc == kReg16
	if strings.EqualFold(strings.TrimSpace(in.Operands[0]), "DX") {
		return []byte{inOutOp(0xEE, w)}, nil
	}
	v, err := parseByteValue(in.Operands[0], resolve)
	if err != nil {
		return nil, wrap(in.Line, err)
	}
	return []byte{inOutOp(0xE6, w), v}, nil
}

func inOutOp(base byte, w bool) byte {
	if w {
		return base | 1
	}
	return base
}

// encodeRel8 codifica un salto a spiazzamento di 8 bit (opcode + rel8).
func encodeRel8(op byte, target string, pc, size int, resolve arch.Resolver, line int) ([]byte, error) {
	dest, err := parseValue(target, resolve)
	if err != nil {
		return nil, wrap(line, err)
	}
	rel := dest - (pc + size)
	if rel < -128 || rel > 127 {
		return nil, fmt.Errorf("riga %d: salto fuori portata per rel8 (%d)", line, rel)
	}
	return []byte{op, byte(rel)}, nil
}

// encodeRel16 codifica un salto/chiamata a spiazzamento di 16 bit.
func encodeRel16(op byte, target string, pc int, resolve arch.Resolver, line int) ([]byte, error) {
	dest, err := parseValue(target, resolve)
	if err != nil {
		return nil, wrap(line, err)
	}
	rel := dest - (pc + 3)
	return []byte{op, byte(rel), byte(rel >> 8)}, nil
}

// --- helper di arieta' e parsing ---

func checkArity(in arch.Instruction, want, size int) (int, error) {
	if err := mustArity(in, want); err != nil {
		return 0, err
	}
	return size, nil
}

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

func parseValue(s string, resolve arch.Resolver) (int, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return 0, fmt.Errorf("operando vuoto")
	}
	if t[0] == '-' || (t[0] >= '0' && t[0] <= '9') {
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
	neg := false
	if strings.HasPrefix(t, "-") {
		neg = true
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
