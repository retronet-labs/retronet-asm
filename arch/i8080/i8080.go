// Package i8080 implementa l'architettura Intel 8080 per retronet-asm.
package i8080

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
)

type kind int

const (
	simple kind = iota
	reg
	regReg
	regImm
	pair
	pairImm
	addr
	imm8
	port
	rst
)

type instr struct {
	op   byte
	kind kind
}

var regs = map[string]byte{"B": 0, "C": 1, "D": 2, "E": 3, "H": 4, "L": 5, "M": 6, "A": 7}
var pairs = map[string]byte{"B": 0, "D": 1, "H": 2, "SP": 3}
var stackPairs = map[string]byte{"B": 0, "D": 1, "H": 2, "PSW": 3}

var set = buildSet()

func buildSet() map[string]instr {
	m := map[string]instr{
		"NOP": {0x00, simple}, "HLT": {0x76, simple},
		"RLC": {0x07, simple}, "RRC": {0x0F, simple}, "RAL": {0x17, simple}, "RAR": {0x1F, simple},
		"DAA": {0x27, simple}, "CMA": {0x2F, simple}, "STC": {0x37, simple}, "CMC": {0x3F, simple},
		"XCHG": {0xEB, simple}, "XTHL": {0xE3, simple}, "SPHL": {0xF9, simple}, "PCHL": {0xE9, simple},
		"RET": {0xC9, simple}, "EI": {0xFB, simple}, "DI": {0xF3, simple},

		"MOV": {0x40, regReg}, "MVI": {0x06, regImm},
		"LXI": {0x01, pairImm}, "INX": {0x03, pair}, "DCX": {0x0B, pair}, "DAD": {0x09, pair},
		"INR": {0x04, reg}, "DCR": {0x05, reg},
		"LDAX": {0x0A, pair}, "STAX": {0x02, pair},
		"LDA": {0x3A, addr}, "STA": {0x32, addr}, "LHLD": {0x2A, addr}, "SHLD": {0x22, addr},
		"JMP": {0xC3, addr}, "CALL": {0xCD, addr},
		"PUSH": {0xC5, pair}, "POP": {0xC1, pair},
		"IN": {0xDB, port}, "OUT": {0xD3, port}, "RST": {0xC7, rst},
	}
	for g, name := range []string{"ADD", "ADC", "SUB", "SBB", "ANA", "XRA", "ORA", "CMP"} {
		m[name] = instr{byte(0x80 | (g << 3)), reg}
	}
	for g, name := range []string{"ADI", "ACI", "SUI", "SBI", "ANI", "XRI", "ORI", "CPI"} {
		m[name] = instr{byte(0xC6 | (g << 3)), imm8}
	}
	conds := []string{"NZ", "Z", "NC", "C", "PO", "PE", "P", "M"}
	for c, name := range conds {
		m["J"+name] = instr{byte(0xC2 | (c << 3)), addr}
		m["C"+name] = instr{byte(0xC4 | (c << 3)), addr}
		m["R"+name] = instr{byte(0xC0 | (c << 3)), simple}
	}
	return m
}

type I8080 struct{}

func New() arch.Arch { return I8080{} }

func (I8080) Name() string { return "i8080" }

func (I8080) Size(in arch.Instruction) (int, error) {
	ins, ok := set[in.Mnemonic]
	if !ok {
		return 0, fmt.Errorf("riga %d: mnemonico sconosciuto %q", in.Line, in.Mnemonic)
	}
	if err := validateArity(in, ins); err != nil {
		return 0, err
	}
	switch ins.kind {
	case regImm, imm8, port:
		return 2, nil
	case pairImm, addr:
		return 3, nil
	default:
		return 1, nil
	}
}

func (I8080) Encode(in arch.Instruction, pc int, resolve arch.Resolver) ([]byte, error) {
	ins, ok := set[in.Mnemonic]
	if !ok {
		return nil, fmt.Errorf("riga %d: mnemonico sconosciuto %q", in.Line, in.Mnemonic)
	}
	if err := validateArity(in, ins); err != nil {
		return nil, err
	}
	switch ins.kind {
	case simple:
		return []byte{ins.op}, nil
	case reg:
		r, err := parseReg(in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if ins.op&0xC0 == 0x80 {
			return []byte{ins.op | r}, nil
		}
		return []byte{ins.op | r<<3}, nil
	case regReg:
		dst, err := parseReg(in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		src, err := parseReg(in.Operands[1])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if dst == 6 && src == 6 {
			return nil, fmt.Errorf("riga %d: MOV M,M codifica HLT sull'8080", in.Line)
		}
		return []byte{ins.op | dst<<3 | src}, nil
	case regImm:
		r, err := parseReg(in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		v, err := parseByteValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{ins.op | r<<3, v}, nil
	case pair:
		p, err := parsePair(in.Mnemonic, in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if (in.Mnemonic == "LDAX" || in.Mnemonic == "STAX") && p > 1 {
			return nil, fmt.Errorf("riga %d: %s accetta solo B o D", in.Line, in.Mnemonic)
		}
		return []byte{ins.op | p<<4}, nil
	case pairImm:
		p, err := parsePair(in.Mnemonic, in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		v, err := parseWordValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{ins.op | p<<4, byte(v), byte(v >> 8)}, nil
	case addr:
		v, err := parseWordValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{ins.op, byte(v), byte(v >> 8)}, nil
	case imm8, port:
		v, err := parseByteValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{ins.op, v}, nil
	case rst:
		v, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if v < 0 || v > 7 {
			return nil, fmt.Errorf("riga %d: vettore RST %d fuori range 0-7", in.Line, v)
		}
		return []byte{ins.op | byte(v)<<3}, nil
	}
	return nil, fmt.Errorf("riga %d: tipo di codifica non gestito per %s", in.Line, in.Mnemonic)
}

func validateArity(in arch.Instruction, ins instr) error {
	want := 0
	switch ins.kind {
	case reg, pair, addr, imm8, port, rst:
		want = 1
	case regReg, regImm, pairImm:
		want = 2
	}
	if len(in.Operands) != want {
		return fmt.Errorf("riga %d: %s vuole %d operandi, trovati %d", in.Line, in.Mnemonic, want, len(in.Operands))
	}
	return nil
}

func parseReg(s string) (byte, error) {
	r, ok := regs[strings.ToUpper(strings.TrimSpace(s))]
	if !ok {
		return 0, fmt.Errorf("registro non valido %q", s)
	}
	return r, nil
}

func parsePair(mnemonic string, s string) (byte, error) {
	name := strings.ToUpper(strings.TrimSpace(s))
	if mnemonic == "PUSH" || mnemonic == "POP" {
		p, ok := stackPairs[name]
		if !ok {
			return 0, fmt.Errorf("coppia stack non valida %q", s)
		}
		return p, nil
	}
	p, ok := pairs[name]
	if !ok {
		return 0, fmt.Errorf("coppia registro non valida %q", s)
	}
	return p, nil
}

func parseByteValue(s string, resolve arch.Resolver) (byte, error) {
	v, err := parseValue(s, resolve)
	if err != nil {
		return 0, err
	}
	if v < 0 || v > 0xFF {
		return 0, fmt.Errorf("valore %d fuori range 0-255", v)
	}
	return byte(v), nil
}

func parseWordValue(s string, resolve arch.Resolver) (uint16, error) {
	v, err := parseValue(s, resolve)
	if err != nil {
		return 0, err
	}
	if v < 0 || v > 0xFFFF {
		return 0, fmt.Errorf("indirizzo 0x%X fuori range 16 bit", v)
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
	return int(n), nil
}

func wrap(line int, err error) error { return fmt.Errorf("riga %d: %w", line, err) }
