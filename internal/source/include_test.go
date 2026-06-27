package source

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandIncludes(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "lib.asm"), ".equ BDOS 0x0005\n")
	main := filepath.Join(dir, "main.asm")
	mustWrite(t, main, ".arch i8080\n.include \"lib.asm\"\n.com\n")

	got, err := ExpandIncludes(main)
	if err != nil {
		t.Fatal(err)
	}
	want := ".arch i8080\n.equ BDOS 0x0005\n.com\n"
	if got != want {
		t.Fatalf("expanded=%q want=%q", got, want)
	}
}

func TestExpandIncludesNested(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "lib"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(dir, "lib", "defs.asm"), ".equ OUT 1\n")
	mustWrite(t, filepath.Join(dir, "lib", "all.asm"), ".include \"defs.asm\"\n.equ IN 0\n")
	main := filepath.Join(dir, "main.asm")
	mustWrite(t, main, ".include \"lib/all.asm\"\nHLT\n")

	got, err := ExpandIncludes(main)
	if err != nil {
		t.Fatal(err)
	}
	if want := ".equ OUT 1\n.equ IN 0\nHLT\n"; got != want {
		t.Fatalf("expanded=%q want=%q", got, want)
	}
}

func TestExpandIncludesRejectsEscapingRoot(t *testing.T) {
	dir := t.TempDir()
	parent := filepath.Dir(dir)
	mustWrite(t, filepath.Join(parent, "outside.asm"), ".equ BAD 1\n")
	main := filepath.Join(dir, "main.asm")
	mustWrite(t, main, ".include \"../outside.asm\"\n")

	if _, err := ExpandIncludes(main); err == nil || !strings.Contains(err.Error(), "fuori") {
		t.Fatalf("err=%v, want fuori-root", err)
	}
}

func TestExpandIncludesRejectsCycles(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.asm")
	b := filepath.Join(dir, "b.asm")
	mustWrite(t, a, ".include \"b.asm\"\n")
	mustWrite(t, b, ".include \"a.asm\"\n")

	if _, err := ExpandIncludes(a); err == nil || !strings.Contains(err.Error(), "ciclico") {
		t.Fatalf("err=%v, want ciclo", err)
	}
}

func mustWrite(t *testing.T, path string, data string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}
