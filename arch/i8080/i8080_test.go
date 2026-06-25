package i8080

import (
	"reflect"
	"testing"

	"github.com/retronet-labs/retronet-asm/arch"
)

func TestEncodeSmallProgram(t *testing.T) {
	a := New()
	src := []arch.Instruction{
		{Mnemonic: "LXI", Operands: []string{"H", "0x1234"}, Line: 1},
		{Mnemonic: "MVI", Operands: []string{"A", "0x2A"}, Line: 2},
		{Mnemonic: "MOV", Operands: []string{"M", "A"}, Line: 3},
		{Mnemonic: "HLT", Line: 4},
	}
	var out []byte
	pc := 0
	for _, in := range src {
		size, err := a.Size(in)
		if err != nil {
			t.Fatal(err)
		}
		code, err := a.Encode(in, pc, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(code) != size {
			t.Fatalf("%s len=%d size=%d", in.Mnemonic, len(code), size)
		}
		out = append(out, code...)
		pc += size
	}
	want := []byte{0x21, 0x34, 0x12, 0x3E, 0x2A, 0x77, 0x76}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("code=% X want=% X", out, want)
	}
}

func TestLabelsAndConditionals(t *testing.T) {
	a := New()
	in := arch.Instruction{Mnemonic: "JNZ", Operands: []string{"loop"}, Line: 10}
	code, err := a.Encode(in, 0, func(name string) (int, bool) {
		return 0x0100, name == "loop"
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0xC2, 0x00, 0x01}
	if !reflect.DeepEqual(code, want) {
		t.Fatalf("code=% X want=% X", code, want)
	}
}

func TestRejectsMOVMM(t *testing.T) {
	_, err := New().Encode(arch.Instruction{Mnemonic: "MOV", Operands: []string{"M", "M"}, Line: 1}, 0, nil)
	if err == nil {
		t.Fatal("MOV M,M accepted, want error")
	}
}
