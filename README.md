# retronet-asm

Assembler modulare e multi-architettura dell'ecosistema **RetroNet**. Traduce un
sorgente testuale `.asm` in una ROM binaria `.rom`, risolvendo le label
automaticamente — niente più programmi scritti come array di byte contati a mano.

Prima architettura supportata: **Intel 4004** (`i4004`). In arrivo: `i8008`, `i8080`.

La ROM prodotta è eseguibile dall'emulatore
[retronet-4004](https://github.com/retronet-labs/retronet-4004).

---

## Quick start

```bash
# assembla un sorgente .asm in una .rom
go run ./cmd/retronet-asm build examples/somma-multicifra.asm -o out.rom

# (senza -o, l'output prende il nome dell'input: examples/somma-multicifra.rom)
go run ./cmd/retronet-asm build examples/somma-multicifra.asm

# tutti i test
go test ./...

# compila il binario
go build -o retronet-asm ./cmd/retronet-asm
```

Poi esegui la ROM con l'emulatore:

```bash
retronet-4004 -trace -dump-ram out.rom
```

---

## Sintassi in breve

- `.arch <nome>` sulla prima riga sceglie l'architettura (default `i4004`).
- Una istruzione per riga; `;` inizia un commento.
- `label:` definisce una label (anche `loop: ADD R1` sulla stessa riga).
- Registri `R0`–`R15`; numeri decimali (`12`) o esadecimali (`0x0C`).
- Virgola tra operandi opzionale; mnemonici case-insensitive, label case-sensitive.
- Arresto: `halt: JUN halt`.

```asm
        LDM 0
        DCL
        FIM R0, 0x03      ; R1 = 3
        LDM 12
        XCH R4
loop:   ADD R1
        ISZ R4, loop
        WRM
halt:   JUN halt
```

Riferimento completo: [`docs/sintassi-asm.md`](docs/sintassi-asm.md).

---

## Esempi

In [`examples/`](examples/):

| File | Cosa fa |
|------|---------|
| `hello4004.asm`     | scrive 5 sulla porta di output (il "ciao mondo") |
| `add.asm`           | addizione binaria 4 + 3 |
| `moltiplicazione.asm` | 3 × 4 con loop ISZ |
| `somma-bcd.asm`     | calcolatrice BCD a cifra singola (7 + 5) |
| `somma-multicifra.asm` | addizione BCD multi-cifra (47 + 58 = 105) |
| `sottrazione-bcd.asm` | sottrazione BCD a cifra singola (7 − 5, TCS) |
| `sottrazione-multicifra.asm` | sottrazione BCD multi-cifra (52 − 27 = 25) |
| `moltiplicazione-bcd.asm` | moltiplicazione per addizioni ripetute (25 × 5 = 125) |
| `divisione-bcd.asm` | divisione per sottrazioni ripetute, con JCN (7 / 2 = 3 r 1) |

Gli ultimi tre, assemblati, producono **gli stessi byte** delle ROM di esempio
di retronet-4004 (`testdata/`): è la validazione incrociata assembler↔emulatore.

---

## Architettura

Pipeline a stadi indipendenti, con la conoscenza della CPU isolata dietro
un'interfaccia (`arch.Arch`):

```
.asm → lexer → parser → emitter (2 passate) → .rom
                            └── interroga arch/i4004
```

| Stadio | Pacchetto | Doc |
|--------|-----------|-----|
| Tokenizer        | `internal/lexer`   | [docs/lexer.md](docs/lexer.md) |
| Parser           | `internal/parser`  | [docs/parser.md](docs/parser.md) |
| Symbol table + emitter | `internal/symbols`, `internal/emitter` | [docs/due-passate.md](docs/due-passate.md) |
| Backend Intel 4004 | `arch/i4004`     | [docs/arch-i4004.md](docs/arch-i4004.md) |

Lexer, parser ed emitter sono **indipendenti dall'architettura**: aggiungere
`i8008`/`i8080` significa scrivere un nuovo pacchetto `arch/` che implementa
`arch.Arch`, senza toccare il resto.

---

## Roadmap

- [x] Backend `i4004`: tabella istruzioni, dimensionamento, codifica
- [x] Lexer, parser, symbol table, emitter a due passate
- [x] CLI `build`
- [x] Esempi + validazione contro le ROM golden di retronet-4004
- [ ] Backend `i8008`
- [ ] Backend `i8080`

---

## Licenza

MIT.
