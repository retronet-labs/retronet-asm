// Package i4004 implementa l'architettura Intel 4004 per retronet-asm:
// la tabella delle istruzioni e la logica di dimensionamento e codifica.
package i4004

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
)

// kind classifica come un'istruzione viene codificata (quanti byte, quali operandi).
type kind int

const (
	simple   kind = iota // 1 byte, nessun operando
	reg                  // 1 byte + registro
	imm                  // 1 byte + immediato 0-15
	regPair              // 1 byte + registro pari (coppia Rr/Rr+1)
	addr12               // 2 byte + indirizzo 12 bit / label
	condAddr             // 2 byte + condizione + indirizzo 8 bit / label
	regAddr              // 2 byte + registro + indirizzo 8 bit / label
	regImm               // 2 byte + registro pari + dato 8 bit
)

// operands è il numero di operandi attesi dal tipo.
func (k kind) operands() int {
	switch k {
	case simple:
		return 0
	case condAddr, regAddr, regImm:
		return 2
	default:
		return 1
	}
}

// size è la lunghezza in byte dell'istruzione (1 o 2).
func (k kind) size() int {
	switch k {
	case addr12, condAddr, regAddr, regImm:
		return 2
	default:
		return 1
	}
}

// instr descrive un'istruzione: opcode base + tipo di codifica.
type instr struct {
	op   byte
	kind kind
}

// set mappa ogni mnemonico (MAIUSCOLO) alla sua descrizione.
var set = map[string]instr{
	// 1 byte, nessun operando
	"NOP": {0x00, simple},
	"WRM": {0xE0, simple}, "WMP": {0xE1, simple}, "WRR": {0xE2, simple}, "WPM": {0xE3, simple},
	"WR0": {0xE4, simple}, "WR1": {0xE5, simple}, "WR2": {0xE6, simple}, "WR3": {0xE7, simple},
	"SBM": {0xE8, simple}, "RDM": {0xE9, simple}, "RDR": {0xEA, simple}, "ADM": {0xEB, simple},
	"RD0": {0xEC, simple}, "RD1": {0xED, simple}, "RD2": {0xEE, simple}, "RD3": {0xEF, simple},
	"CLB": {0xF0, simple}, "CLC": {0xF1, simple}, "IAC": {0xF2, simple}, "CMC": {0xF3, simple},
	"CMA": {0xF4, simple}, "RAL": {0xF5, simple}, "RAR": {0xF6, simple}, "TCC": {0xF7, simple},
	"DAC": {0xF8, simple}, "TCS": {0xF9, simple}, "STC": {0xFA, simple}, "DAA": {0xFB, simple},
	"KBP": {0xFC, simple}, "DCL": {0xFD, simple},

	// 1 byte + registro
	"INC": {0x60, reg}, "ADD": {0x80, reg}, "SUB": {0x90, reg}, "LD": {0xA0, reg}, "XCH": {0xB0, reg},

	// 1 byte + immediato 0-15
	"BBL": {0xC0, imm}, "LDM": {0xD0, imm},

	// 1 byte + coppia di registri pari
	"SRC": {0x21, regPair}, "FIN": {0x30, regPair}, "JIN": {0x31, regPair},

	// 2 byte
	"JUN": {0x40, addr12}, "JMS": {0x50, addr12},
	"JCN": {0x10, condAddr},
	"ISZ": {0x70, regAddr},
	"FIM": {0x20, regImm},
}

// I4004 implementa arch.Arch per l'Intel 4004.
type I4004 struct{}

// New restituisce l'architettura 4004 come arch.Arch.
func New() arch.Arch { return I4004{} }

func (I4004) Name() string { return "i4004" }

// Size valida il mnemonico e l'arità degli operandi e restituisce 1 o 2 byte.
func (I4004) Size(in arch.Instruction) (int, error) {
	ins, ok := set[in.Mnemonic]
	if !ok {
		return 0, fmt.Errorf("riga %d: mnemonico sconosciuto %q", in.Line, in.Mnemonic)
	}
	if len(in.Operands) != ins.kind.operands() {
		return 0, fmt.Errorf("riga %d: %s vuole %d operandi, trovati %d",
			in.Line, in.Mnemonic, ins.kind.operands(), len(in.Operands))
	}
	return ins.kind.size(), nil
}

// Encode produce i byte dell'istruzione, risolvendo le label tramite resolve.
func (I4004) Encode(in arch.Instruction, pc int, resolve arch.Resolver) ([]byte, error) {
	ins, ok := set[in.Mnemonic]
	if !ok {
		return nil, fmt.Errorf("riga %d: mnemonico sconosciuto %q", in.Line, in.Mnemonic)
	}
	if len(in.Operands) != ins.kind.operands() {
		return nil, fmt.Errorf("riga %d: %s vuole %d operandi, trovati %d",
			in.Line, in.Mnemonic, ins.kind.operands(), len(in.Operands))
	}

	switch ins.kind {
	case simple:
		return []byte{ins.op}, nil

	case reg:
		r, err := parseReg(in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{ins.op | r}, nil

	case imm:
		v, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if v < 0 || v > 15 {
			return nil, fmt.Errorf("riga %d: immediato %d fuori range 0-15", in.Line, v)
		}
		return []byte{ins.op | byte(v)}, nil

	case regPair:
		r, err := parseReg(in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{ins.op | (r &^ 1)}, nil // forza il registro pari

	case addr12:
		addr, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if addr < 0 || addr > 0x0FFF {
			return nil, fmt.Errorf("riga %d: indirizzo 0x%X fuori range 12 bit", in.Line, addr)
		}
		return []byte{ins.op | byte(addr>>8&0x0F), byte(addr & 0xFF)}, nil

	case condAddr:
		cond, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		addr, err := parseValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{ins.op | byte(cond&0x0F), byte(addr & 0xFF)}, nil

	case regAddr:
		r, err := parseReg(in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		addr, err := parseValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		return []byte{ins.op | r, byte(addr & 0xFF)}, nil

	case regImm:
		r, err := parseReg(in.Operands[0])
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		data, err := parseValue(in.Operands[1], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if data < 0 || data > 0xFF {
			return nil, fmt.Errorf("riga %d: dato 0x%X fuori range 8 bit", in.Line, data)
		}
		return []byte{ins.op | (r &^ 1), byte(data)}, nil
	}
	return nil, fmt.Errorf("riga %d: tipo di codifica non gestito per %s", in.Line, in.Mnemonic)
}

func wrap(line int, err error) error { return fmt.Errorf("riga %d: %w", line, err) }

// parseReg interpreta "R0".."R15" (case-insensitive) come numero di registro 0-15.
func parseReg(s string) (byte, error) {
	t := strings.ToUpper(strings.TrimSpace(s))
	if len(t) < 2 || t[0] != 'R' {
		return 0, fmt.Errorf("registro non valido %q (atteso R0-R15)", s)
	}
	n, err := strconv.Atoi(t[1:])
	if err != nil || n < 0 || n > 15 {
		return 0, fmt.Errorf("registro non valido %q (atteso R0-R15)", s)
	}
	return byte(n), nil
}

// parseNum interpreta un numero decimale ("12") o esadecimale ("0x0C").
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

// parseValue interpreta un operando numerico: un numero (decimale, esadecimale o
// negativo) oppure un simbolo (label o costante .equ) risolto con resolve.
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
