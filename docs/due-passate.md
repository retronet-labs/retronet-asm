# Le due passate: symbols + emitter

L'ultimo stadio dell'assembler trasforma gli **statement** del parser in **byte**
di ROM. Lo fa in **due passate**, con l'aiuto di una **symbol table**.

```
[]Stmt ──▶ [EMITTER]  ──passata 1──▶  symbol table (label → indirizzo)
                      ──passata 2──▶  byte ROM
```

File: [`internal/emitter/emitter.go`](../internal/emitter/emitter.go) ·
[`internal/symbols/symbols.go`](../internal/symbols/symbols.go)

---

## Perché due passate

Un salto può riferirsi a una label **definita più avanti**:

```asm
        JUN end      ; 'end' non è ancora stata vista
        NOP
end:    JUN end      ; ...viene definita qui
```

Per codificare `JUN end` serve l'indirizzo di `end`. Ma per conoscere
quell'indirizzo bisogna prima sapere quanti byte occupano le istruzioni che
stanno **prima** di `end`. Quindi:

- **Passata 1 — indirizzi.** Si scorre il programma tenendo un *program counter*
  `pc`. Per ogni istruzione si chiede all'architettura la sua **dimensione**
  (`Size`) e si avanza `pc`. Ogni `label:` viene registrata nella symbol table
  **all'indirizzo `pc` corrente** (cioè punta all'istruzione che la segue).
  Qui non si emette ancora nessun byte.

- **Passata 2 — byte.** Si riscorre il programma e si **codifica** ogni
  istruzione (`Encode`), passando la funzione di lookup della symbol table come
  `Resolver`. Ora `JUN end` trova l'indirizzo di `end` e produce i byte giusti.

La separazione `Size`/`Encode` dell'interfaccia `arch.Arch` esiste proprio per
questo: `Size` non ha bisogno delle label, `Encode` sì.

---

## La symbol table

Una mappa `nome → indirizzo`, con due operazioni:

```go
func (t *Table) Define(name string, addr int) error // errore se label duplicata
func (t *Table) Lookup(name string) (int, bool)      // firma di arch.Resolver
```

`Lookup` ha **esattamente** la firma di `arch.Resolver`, quindi l'emitter la
passa direttamente a `Encode` senza adattatori:

```go
b, err := a.Encode(*st.Instr, pc, syms.Lookup)
```

`Define` rifiuta una label già presente: due `loop:` nello stesso file sono un
errore (con il numero di riga).

---

## L'emitter, passo per passo

```go
func Assemble(stmts []parser.Stmt, a arch.Arch) ([]byte, error) {
	syms := symbols.New()

	// Passata 1: indirizzi + label
	pc := 0
	for _, st := range stmts {
		if st.Label != "" {
			syms.Define(st.Label, pc) // (errore gestito)
		}
		if st.Instr != nil {
			sz, _ := a.Size(*st.Instr) // (errore gestito)
			pc += sz
		}
	}

	// Passata 2: byte
	pc = 0
	for _, st := range stmts {
		if st.Instr == nil {
			continue // riga di sola label: 0 byte
		}
		b, _ := a.Encode(*st.Instr, pc, syms.Lookup) // (errore gestito)
		code = append(code, b...)
		pc += len(b)
	}
	return code, nil
}
```

Le due passate **devono produrre gli stessi indirizzi**: per questo `pc` viene
fatto avanzare con la stessa logica in entrambe (in passata 1 con `Size`, in
passata 2 con `len(byte)` — che per costruzione coincidono).

---

## La validazione che chiude il cerchio

I byte prodotti dall'emitter sono confrontati, nei test, con le **ROM golden** di
[`retronet-4004`](https://github.com/retronet-labs/retronet-4004) (`testdata/`):
gli stessi byte che l'emulatore esegue correttamente. Se l'assembler di un
sorgente `.asm` produce byte identici alla ROM costruita a mano nell'emulatore,
i due moduli sono coerenti tra loro. È il test più importante del progetto.

I test sono in [`internal/emitter/emitter_test.go`](../internal/emitter/emitter_test.go)
e [`internal/symbols/symbols_test.go`](../internal/symbols/symbols_test.go).
