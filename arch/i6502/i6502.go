// Package i6502 implementa il backend MOS/NMOS 6502 per retronet-asm.
package i6502

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
)

type mode int

const (
	imp mode = iota
	acc
	imm
	zp
	zpx
	zpy
	abs
	absx
	absy
	ind
	indx
	indy
	rel
)

type I6502 struct{}

// New restituisce il backend 6502.
func New() arch.Arch { return I6502{} }

func (I6502) Name() string { return "i6502" }

var set = buildSet()

func buildSet() map[string]map[mode]byte {
	m := map[string]map[mode]byte{}
	add := func(name string, md mode, op byte) {
		if m[name] == nil {
			m[name] = map[mode]byte{}
		}
		m[name][md] = op
	}
	add(impName("BRK"), imp, 0x00)
	add(impName("NOP"), imp, 0xEA)
	add(impName("RTI"), imp, 0x40)
	add(impName("RTS"), imp, 0x60)
	for _, x := range []struct {
		n  string
		op byte
	}{
		{"CLC", 0x18}, {"SEC", 0x38}, {"CLI", 0x58}, {"SEI", 0x78},
		{"CLV", 0xB8}, {"CLD", 0xD8}, {"SED", 0xF8},
		{"PHA", 0x48}, {"PHP", 0x08}, {"PLA", 0x68}, {"PLP", 0x28},
		{"TAX", 0xAA}, {"TAY", 0xA8}, {"TSX", 0xBA}, {"TXA", 0x8A},
		{"TXS", 0x9A}, {"TYA", 0x98}, {"INX", 0xE8}, {"INY", 0xC8},
		{"DEX", 0xCA}, {"DEY", 0x88},
	} {
		add(x.n, imp, x.op)
	}
	for _, x := range []struct {
		n  string
		op byte
	}{
		{"BPL", 0x10}, {"BMI", 0x30}, {"BVC", 0x50}, {"BVS", 0x70},
		{"BCC", 0x90}, {"BCS", 0xB0}, {"BNE", 0xD0}, {"BEQ", 0xF0},
	} {
		add(x.n, rel, x.op)
	}
	addALU := func(n string, ops ...byte) {
		add(n, imm, ops[0])
		add(n, zp, ops[1])
		add(n, zpx, ops[2])
		add(n, abs, ops[3])
		add(n, absx, ops[4])
		add(n, absy, ops[5])
		add(n, indx, ops[6])
		add(n, indy, ops[7])
	}
	addALU("ORA", 0x09, 0x05, 0x15, 0x0D, 0x1D, 0x19, 0x01, 0x11)
	addALU("AND", 0x29, 0x25, 0x35, 0x2D, 0x3D, 0x39, 0x21, 0x31)
	addALU("EOR", 0x49, 0x45, 0x55, 0x4D, 0x5D, 0x59, 0x41, 0x51)
	addALU("ADC", 0x69, 0x65, 0x75, 0x6D, 0x7D, 0x79, 0x61, 0x71)
	addALU("SBC", 0xE9, 0xE5, 0xF5, 0xED, 0xFD, 0xF9, 0xE1, 0xF1)
	addALU("CMP", 0xC9, 0xC5, 0xD5, 0xCD, 0xDD, 0xD9, 0xC1, 0xD1)
	addLoad := addALU
	addLoad("LDA", 0xA9, 0xA5, 0xB5, 0xAD, 0xBD, 0xB9, 0xA1, 0xB1)
	add("LDX", imm, 0xA2)
	add("LDX", zp, 0xA6)
	add("LDX", zpy, 0xB6)
	add("LDX", abs, 0xAE)
	add("LDX", absy, 0xBE)
	add("LDY", imm, 0xA0)
	add("LDY", zp, 0xA4)
	add("LDY", zpx, 0xB4)
	add("LDY", abs, 0xAC)
	add("LDY", absx, 0xBC)
	add("STA", zp, 0x85)
	add("STA", zpx, 0x95)
	add("STA", abs, 0x8D)
	add("STA", absx, 0x9D)
	add("STA", absy, 0x99)
	add("STA", indx, 0x81)
	add("STA", indy, 0x91)
	add("STX", zp, 0x86)
	add("STX", zpy, 0x96)
	add("STX", abs, 0x8E)
	add("STY", zp, 0x84)
	add("STY", zpx, 0x94)
	add("STY", abs, 0x8C)
	for _, x := range []struct {
		n            string
		acc, zp, zpx byte
		abs, absx    byte
	}{
		{"ASL", 0x0A, 0x06, 0x16, 0x0E, 0x1E},
		{"ROL", 0x2A, 0x26, 0x36, 0x2E, 0x3E},
		{"LSR", 0x4A, 0x46, 0x56, 0x4E, 0x5E},
		{"ROR", 0x6A, 0x66, 0x76, 0x6E, 0x7E},
	} {
		add(x.n, acc, x.acc)
		add(x.n, zp, x.zp)
		add(x.n, zpx, x.zpx)
		add(x.n, abs, x.abs)
		add(x.n, absx, x.absx)
	}
	add("INC", zp, 0xE6)
	add("INC", zpx, 0xF6)
	add("INC", abs, 0xEE)
	add("INC", absx, 0xFE)
	add("DEC", zp, 0xC6)
	add("DEC", zpx, 0xD6)
	add("DEC", abs, 0xCE)
	add("DEC", absx, 0xDE)
	add("CPX", imm, 0xE0)
	add("CPX", zp, 0xE4)
	add("CPX", abs, 0xEC)
	add("CPY", imm, 0xC0)
	add("CPY", zp, 0xC4)
	add("CPY", abs, 0xCC)
	add("BIT", zp, 0x24)
	add("BIT", abs, 0x2C)
	add("JMP", abs, 0x4C)
	add("JMP", ind, 0x6C)
	add("JSR", abs, 0x20)
	return m
}

func impName(s string) string { return s }

func (a I6502) Size(in arch.Instruction) (int, error) {
	md, _, err := classifyOperand(in, nil)
	if err != nil {
		return 0, err
	}
	if _, ok := opcode(in.Mnemonic, md); !ok {
		return 0, fmt.Errorf("riga %d: forma non valida: %s %s", in.Line, in.Mnemonic, strings.Join(in.Operands, ","))
	}
	return modeSize(md), nil
}

func (a I6502) Encode(in arch.Instruction, pc int, resolve arch.Resolver) ([]byte, error) {
	md, operand, err := classifyOperand(in, resolve)
	if err != nil {
		return nil, err
	}
	op, ok := opcode(in.Mnemonic, md)
	if !ok {
		return nil, fmt.Errorf("riga %d: forma non valida: %s %s", in.Line, in.Mnemonic, strings.Join(in.Operands, ","))
	}
	switch md {
	case imp, acc:
		return []byte{op}, nil
	case imm, zp, zpx, zpy, indx, indy:
		v, err := operand.byteValue(resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{op, v}, nil
	case rel:
		target, err := operand.wordValue(resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		off := int(target) - (pc + 2)
		if off < -128 || off > 127 {
			return nil, fmt.Errorf("riga %d: branch fuori range (%d)", in.Line, off)
		}
		return []byte{op, byte(int8(off))}, nil
	default:
		v, err := operand.wordValue(resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{op, byte(v), byte(v >> 8)}, nil
	}
}

func opcode(mnemonic string, md mode) (byte, bool) {
	forms, ok := set[mnemonic]
	if !ok {
		return 0, false
	}
	op, ok := forms[md]
	return op, ok
}

func modeSize(md mode) int {
	switch md {
	case imp, acc:
		return 1
	case imm, zp, zpx, zpy, indx, indy, rel:
		return 2
	default:
		return 3
	}
}

type operand struct {
	text  string
	force byte // 0 none, '<' low/zeropage, '>' high byte immediate
}

func classifyOperand(in arch.Instruction, resolve arch.Resolver) (mode, operand, error) {
	m := in.Mnemonic
	ops := normalizeOps(in.Operands)
	if _, isBranch := set[m][rel]; isBranch {
		if len(ops) != 1 {
			return 0, operand{}, fmt.Errorf("riga %d: %s vuole 1 operando", in.Line, m)
		}
		return rel, operand{text: ops[0]}, nil
	}
	if len(ops) == 0 {
		return imp, operand{}, nil
	}
	if len(ops) != 1 {
		return 0, operand{}, fmt.Errorf("riga %d: troppi operandi per %s", in.Line, m)
	}
	raw := strings.TrimSpace(ops[0])
	up := strings.ToUpper(raw)
	if up == "A" {
		return acc, operand{}, nil
	}
	if strings.HasPrefix(raw, "#") {
		o := parseOperandText(raw[1:])
		return imm, o, nil
	}
	if strings.HasPrefix(raw, "(") {
		return classifyParen(raw)
	}
	if strings.HasSuffix(up, ",X") {
		o := parseOperandText(raw[:len(raw)-2])
		if zeroPage(o, resolve) {
			return zpx, o, nil
		}
		return absx, o, nil
	}
	if strings.HasSuffix(up, ",Y") {
		o := parseOperandText(raw[:len(raw)-2])
		if zeroPage(o, resolve) {
			return zpy, o, nil
		}
		return absy, o, nil
	}
	o := parseOperandText(raw)
	if m == "JMP" || m == "JSR" {
		return abs, o, nil
	}
	if zeroPage(o, resolve) {
		return zp, o, nil
	}
	return abs, o, nil
}

func normalizeOps(ops []string) []string {
	if len(ops) == 2 {
		u := strings.ToUpper(strings.TrimSpace(ops[1]))
		if u == "X" || u == "Y" {
			return []string{strings.TrimSpace(ops[0]) + "," + u}
		}
	}
	return ops
}

func classifyParen(raw string) (mode, operand, error) {
	up := strings.ToUpper(strings.ReplaceAll(raw, " ", ""))
	switch {
	case strings.HasSuffix(up, ",X)"):
		return indx, parseOperandText(raw[1 : len(raw)-3]), nil
	case strings.HasSuffix(up, "),Y"):
		return indy, parseOperandText(raw[1 : len(raw)-3]), nil
	case strings.HasSuffix(up, ")"):
		return ind, parseOperandText(raw[1 : len(raw)-1]), nil
	default:
		return 0, operand{}, fmt.Errorf("operando indiretto non valido %q", raw)
	}
}

func parseOperandText(s string) operand {
	t := strings.TrimSpace(s)
	if strings.HasPrefix(t, "<") || strings.HasPrefix(t, ">") {
		return operand{text: strings.TrimSpace(t[1:]), force: t[0]}
	}
	return operand{text: t}
}

func zeroPage(o operand, resolve arch.Resolver) bool {
	if o.force == '<' {
		return true
	}
	if !isLiteral(o.text) {
		return false
	}
	v, err := parseNum(o.text)
	return err == nil && v >= 0 && v <= 0xFF
}

func (o operand) byteValue(resolve arch.Resolver) (byte, error) {
	v, err := o.value(resolve)
	if err != nil {
		return 0, err
	}
	if o.force == '>' {
		return byte(v >> 8), nil
	}
	if o.force == '<' {
		return byte(v), nil
	}
	if v < 0 || v > 0xFF {
		return 0, fmt.Errorf("valore 0x%X fuori range 8 bit", v)
	}
	return byte(v), nil
}

func (o operand) wordValue(resolve arch.Resolver) (uint16, error) {
	v, err := o.value(resolve)
	if err != nil {
		return 0, err
	}
	if v < 0 || v > 0xFFFF {
		return 0, fmt.Errorf("indirizzo 0x%X fuori range 16 bit", v)
	}
	return uint16(v), nil
}

func (o operand) value(resolve arch.Resolver) (int, error) {
	if isLiteral(o.text) {
		return parseNum(o.text)
	}
	if resolve == nil {
		return 0, nil
	}
	v, ok := resolve(o.text)
	if !ok {
		return 0, fmt.Errorf("simbolo non definito: %q", o.text)
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

func isLiteral(s string) bool {
	s = strings.TrimSpace(s)
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

func wrap(line int, err error) error { return fmt.Errorf("riga %d: %w", line, err) }
