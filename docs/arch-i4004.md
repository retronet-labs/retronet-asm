# L'architettura: interfaccia `Arch` e backend `i4004`

L'assembler è **multi-architettura**: lexer, parser ed emitter non sanno nulla
del 4004. Tutta la conoscenza specifica di una CPU (quali istruzioni esistono,
quanti byte occupano, come si codificano) vive dietro un'interfaccia, `arch.Arch`,
e in un pacchetto per ogni CPU (`arch/i4004`, in futuro `arch/i8008`, `arch/i8080`).

```
parser ──▶ []Stmt ──▶ [EMITTER] ──interroga──▶ arch.Arch (es. i4004) ──▶ byte
```

File: [`arch/arch.go`](../arch/arch.go) · [`arch/i4004/i4004.go`](../arch/i4004/i4004.go)

---

## Il contratto: `arch.Arch`

```go
type Instruction struct {
	Mnemonic string   // "LDM", "JUN", ... (MAIUSCOLO)
	Operands []string // operandi grezzi: ["R4", "loop"]
	Line     int      // riga sorgente, per gli errori
}

type Resolver func(name string) (addr int, ok bool)

type Arch interface {
	Name() string
	Size(in Instruction) (int, error)
	Encode(in Instruction, pc int, resolve Resolver) ([]byte, error)
}
```

`Instruction` è l'AST minimo che il parser produce: un mnemonico e i suoi
operandi **ancora come stringhe**. È l'architettura a interpretarli.

### Perché `Size` ed `Encode` sono separati

È il cuore del modello a **due passate** (vedi `docs/due-passate.md`):

- **Passata 1 — indirizzi.** Per sapere dove cade ogni istruzione (e quindi a
  quale indirizzo corrisponde una `label:`), serve la **dimensione** di ogni
  istruzione. La dimensione **non dipende mai dalle label** → `Size` non riceve
  il `Resolver`.
- **Passata 2 — byte.** Ora le label hanno un indirizzo noto, quindi si può
  **codificare** → `Encode` riceve il `Resolver` per tradurre `JUN loop` nel suo
  indirizzo.

`Resolver` è un tipo-funzione: l'architettura non sa *come* sono memorizzate le
label, chiede solo "qual è l'indirizzo di `loop`?". L'emitter gli passa una
funzione che cerca nella symbol table.

---

## Il backend `i4004`

`I4004` implementa `arch.Arch`. `New()` lo restituisce come `arch.Arch`.

### Classificazione: `kind`

Ogni istruzione del 4004 ricade in uno di **8 tipi di codifica**, che
determinano quanti operandi servono e quanti byte produce:

| `kind`     | Byte | Operandi | Esempi |
|------------|------|----------|--------|
| `simple`   | 1    | 0        | NOP, DAA, WRM |
| `reg`      | 1    | 1 (registro) | ADD, INC, LD, XCH, SUB |
| `imm`      | 1    | 1 (0–15) | LDM, BBL |
| `regPair`  | 1    | 1 (registro pari) | SRC, FIN, JIN |
| `addr12`   | 2    | 1 (addr/label) | JUN, JMS |
| `condAddr` | 2    | 2 (cond, addr/label) | JCN |
| `regAddr`  | 2    | 2 (registro, addr/label) | ISZ |
| `regImm`   | 2    | 2 (registro pari, dato 8 bit) | FIM |

`kind.operands()` e `kind.size()` derivano arità e dimensione dal tipo.

### La tabella `set`

Un'unica mappa traduce ogni mnemonico nei due dati che servono: **opcode base**
e **kind**.

```go
var set = map[string]instr{
	"NOP": {0x00, simple},
	"ADD": {0x80, reg},
	"LDM": {0xD0, imm},
	"SRC": {0x21, regPair},
	"JUN": {0x40, addr12},
	"JCN": {0x10, condAddr},
	"ISZ": {0x70, regAddr},
	"FIM": {0x20, regImm},
	// ... tutte le 46 istruzioni
}
```

È l'unico punto del progetto che "conosce" gli opcode del 4004.

### `Size`

Fa solo: lookup del mnemonico + controllo dell'arità → restituisce 1 o 2.
Non tocca i valori né le label, quindi può girare nella passata 1.

### `Encode`

Un ramo per ogni `kind`. Regole di codifica:

| `kind`     | byte 0 | byte 1 |
|------------|--------|--------|
| `simple`   | `op` | — |
| `reg`      | `op \| reg` | — |
| `imm`      | `op \| (v & 0x0F)` (v in 0–15) | — |
| `regPair`  | `op \| (reg &^ 1)` (forza pari) | — |
| `addr12`   | `op \| (addr>>8 & 0x0F)` | `addr & 0xFF` |
| `condAddr` | `op \| (cond & 0x0F)` | `addr & 0xFF` |
| `regAddr`  | `op \| reg` | `addr & 0xFF` |
| `regImm`   | `op \| (reg &^ 1)` | `dato & 0xFF` |

Note:

- **`&^ 1`** azzera il bit 0 del numero di registro: le coppie (SRC/FIN/JIN/FIM)
  devono indirizzare un registro pari (`R0`, `R2`, ...).
- **Indirizzi a 8 bit** (`condAddr`, `regAddr`): si emette solo il byte basso.
  Sul 4004 questi salti restano nella *pagina* corrente (i 4 bit alti vengono
  dal PC a runtime), quindi l'assembler scrive solo gli 8 bit bassi della label.
- I controlli di range (immediato 0–15, dato 0–255, indirizzo 0–0xFFF) producono
  errori con il numero di riga.

### Dove entrano le label

Solo dentro `Encode`, attraverso `parseAddr`:

```
operando inizia con una cifra?  → è un numero (parseNum: decimale o 0x esadecimale)
altrimenti                       → è una label → resolve(nome)
```

Così `JUN 0x100` e `JUN loop` passano per lo stesso ramo `addr12`: il primo
risolve un numero, il secondo chiede l'indirizzo al `Resolver`.

### Parsing degli operandi

- `parseReg("R5")` → `5` (case-insensitive, valida 0–15).
- `parseNum("12")` / `parseNum("0x0C")` → intero (decimale o esadecimale).
- `parseAddr` → numero **oppure** label risolta (vedi sopra).

---

## Aggiungere una nuova architettura (in futuro)

Per `i8008`/`i8080` basta un nuovo pacchetto `arch/iXXXX` che implementa
`arch.Arch` (la sua tabella + `Size`/`Encode`). Lexer, parser ed emitter restano
identici: ricevono l'arch da usare e la interrogano. È il senso della struttura
`arch/`.

I test sono in [`arch/i4004/i4004_test.go`](../arch/i4004/i4004_test.go).
