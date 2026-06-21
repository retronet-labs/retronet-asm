// Package i8008 implementa l'architettura Intel 8008 per retronet-asm.
//
// I mnemonici seguono la convenzione del disassembler dell'emulatore
// retronet-8008 (es. LAB, ADM, INB, RFC, HLT): l'assembler e l'emulatore sono
// cosi' speculari. Le codifiche replicano i pattern di bit dell'ISA 8008.
//
// Copre l'intero set base: move tra registri/M, ALU registro e immediato,
// INr/DCr, rotate, controllo di flusso (JMP/CAL e condizionati, RET e
// condizionati, RST), I/O (INP/OUT) e HLT.
package i8008

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/retronet-labs/retronet-asm/arch"
)

// kind classifica come un'istruzione viene codificata.
type kind int

const (
	simple kind = iota // 1 byte, nessun operando
	imm                // 2 byte, 1 immediato 0-255
	addr               // 3 byte, 1 indirizzo/label 14 bit
	rst                // 1 byte, 1 vettore 0-7
	inp                // 1 byte, 1 porta input 0-7
	out                // 1 byte, 1 porta output 8-31
)

func (k kind) operands() int {
	if k == simple {
		return 0
	}
	return 1
}

func (k kind) size() int {
	switch k {
	case imm:
		return 2
	case addr:
		return 3
	default:
		return 1
	}
}

// instr descrive un'istruzione: opcode base + tipo di codifica.
type instr struct {
	op   byte
	kind kind
}

// regNames indicizza i codici registro 8008: A=0 .. L=6, M=7 (pseudo-registro).
var regNames = [8]string{"A", "B", "C", "D", "E", "H", "L", "M"}

// condNames indicizza i codici condizione: C=0, Z=1, S=2, P=3.
var condNames = [4]string{"C", "Z", "S", "P"}

// set mappa ogni mnemonico (MAIUSCOLO) alla sua descrizione.
var set = buildSet()

// buildSet genera la tabella delle istruzioni replicando i pattern di bit
// dell'8008, cosi' i mnemonici coincidono con quelli del disassembler emulatore.
func buildSet() map[string]instr {
	m := map[string]instr{
		"HLT": {0x00, simple}, // 0xFF e 0x01 sono alias di HLT, ma qui emettiamo 0x00
		"RET": {0x07, simple},
		"RLC": {0x02, simple}, "RRC": {0x0A, simple}, "RAL": {0x12, simple}, "RAR": {0x1A, simple},
		"NOP": {0xC0, simple}, // L A,A: trasferimento nullo
	}

	// Move Lr1r2 = 11 DDD SSS, con dst != src (dst == src e' NOP/HLT).
	for dst := 0; dst < 8; dst++ {
		for src := 0; src < 8; src++ {
			if dst == src {
				continue
			}
			op := byte(0xC0 | (dst << 3) | src)
			m["L"+regNames[dst]+regNames[src]] = instr{op, simple}
		}
	}

	// ALU registro = 10 GGG SSS: AD,AC,SU,SB,ND,XR,OR,CP applicati a r o M.
	aluPrefix := [8]string{"AD", "AC", "SU", "SB", "ND", "XR", "OR", "CP"}
	for g := 0; g < 8; g++ {
		for src := 0; src < 8; src++ {
			op := byte(0x80 | (g << 3) | src)
			m[aluPrefix[g]+regNames[src]] = instr{op, simple}
		}
	}

	// Increment/decrement = 00 RRR 000 / 00 RRR 001, solo B..L (codici 1..6).
	for r := 1; r <= 6; r++ {
		m["IN"+regNames[r]] = instr{byte(r << 3), simple}
		m["DC"+regNames[r]] = instr{byte((r << 3) | 0x01), simple}
	}

	// Ritorni condizionati = 00 0CC 011 (flag falso) / 00 1CC 011 (flag vero).
	for cc := 0; cc < 4; cc++ {
		m["RF"+condNames[cc]] = instr{byte(0x03 | (cc << 3)), simple}
		m["RT"+condNames[cc]] = instr{byte(0x23 | (cc << 3)), simple}
	}

	// Immediati (2 byte): load immediato LrI = 00 DDD 110, ALU immediato = 00 GGG 100.
	for dst := 0; dst < 8; dst++ {
		m["L"+regNames[dst]+"I"] = instr{byte(0x06 | (dst << 3)), imm}
	}
	immALU := [8]string{"ADI", "ACI", "SUI", "SBI", "NDI", "XRI", "ORI", "CPI"}
	for g := 0; g < 8; g++ {
		m[immALU[g]] = instr{byte(0x04 | (g << 3)), imm}
	}

	// Indirizzi (3 byte): JMP/CAL incondizionati e condizionati su C/Z/S/P.
	m["JMP"] = instr{0x44, addr}
	m["CAL"] = instr{0x46, addr}
	for cc := 0; cc < 4; cc++ {
		m["JF"+condNames[cc]] = instr{byte(0x40 | (cc << 3)), addr}
		m["JT"+condNames[cc]] = instr{byte(0x60 | (cc << 3)), addr}
		m["CF"+condNames[cc]] = instr{byte(0x42 | (cc << 3)), addr}
		m["CT"+condNames[cc]] = instr{byte(0x62 | (cc << 3)), addr}
	}

	// RST (1 byte) e I/O (1 byte). Per INP/OUT il numero di porta sta nel campo
	// a 5 bit dell'opcode: input 0-7, output 8-31.
	m["RST"] = instr{0x05, rst}
	m["INP"] = instr{0x41, inp}
	m["OUT"] = instr{0x41, out}

	return m
}

// I8008 implementa arch.Arch per l'Intel 8008.
type I8008 struct{}

// New restituisce l'architettura 8008 come arch.Arch.
func New() arch.Arch { return I8008{} }

func (I8008) Name() string { return "i8008" }

// Size valida il mnemonico e l'arita' degli operandi e restituisce 1, 2 o 3 byte.
func (I8008) Size(in arch.Instruction) (int, error) {
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
func (I8008) Encode(in arch.Instruction, pc int, resolve arch.Resolver) ([]byte, error) {
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

	case imm:
		v, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if v < 0 || v > 0xFF {
			return nil, fmt.Errorf("riga %d: immediato %d fuori range 0-255", in.Line, v)
		}
		return []byte{ins.op, byte(v)}, nil

	case addr:
		a, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if a < 0 || a > 0x3FFF {
			return nil, fmt.Errorf("riga %d: indirizzo 0x%X fuori range 14 bit", in.Line, a)
		}
		// low byte, poi high byte mascherato a 6 bit (indirizzo a 14 bit).
		return []byte{ins.op, byte(a & 0xFF), byte((a >> 8) & 0x3F)}, nil

	case rst:
		n, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if n < 0 || n > 7 {
			return nil, fmt.Errorf("riga %d: vettore RST %d fuori range 0-7", in.Line, n)
		}
		return []byte{ins.op | (byte(n) << 3)}, nil

	case inp:
		p, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if p < 0 || p > 7 {
			return nil, fmt.Errorf("riga %d: porta input %d fuori range 0-7", in.Line, p)
		}
		return []byte{ins.op | (byte(p) << 1)}, nil

	case out:
		p, err := parseValue(in.Operands[0], resolve)
		if err != nil {
			return nil, wrap(in.Line, err)
		}
		if p < 8 || p > 31 {
			return nil, fmt.Errorf("riga %d: porta output %d fuori range 8-31", in.Line, p)
		}
		return []byte{ins.op | (byte(p) << 1)}, nil
	}
	return nil, fmt.Errorf("riga %d: tipo di codifica non gestito per %s", in.Line, in.Mnemonic)
}

func wrap(line int, err error) error { return fmt.Errorf("riga %d: %w", line, err) }

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
