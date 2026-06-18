// Package symbols è la tabella delle label: nome → indirizzo nella ROM.
package symbols

import "fmt"

// Table associa ogni label al suo indirizzo.
type Table struct {
	m map[string]int
}

// New crea una tabella vuota.
func New() *Table {
	return &Table{m: make(map[string]int)}
}

// Define registra una label a un indirizzo. Errore se già definita.
func (t *Table) Define(name string, addr int) error {
	if _, ok := t.m[name]; ok {
		return fmt.Errorf("label duplicata: %q", name)
	}
	t.m[name] = addr
	return nil
}

// Lookup restituisce l'indirizzo di una label e true se definita.
// Ha esattamente la firma di arch.Resolver, quindi si passa così com'è all'emitter.
func (t *Table) Lookup(name string) (int, bool) {
	addr, ok := t.m[name]
	return addr, ok
}
