# retronet-asm

Assembler modulare e multi-architettura dell'ecosistema **RetroNet**. Traduce un
sorgente testuale `.asm` in una ROM binaria `.rom`, risolvendo le label
automaticamente — niente più programmi scritti come array di byte contati a mano.

Architetture supportate: **Intel 4004** (`i4004`), **Intel 8008** (`i8008`),
**Intel 8080** (`i8080`) e **Intel 8086/8088** (`i8086`).

Le ROM prodotte sono eseguibili dagli emulatori
[retronet-4004](https://github.com/retronet-labs/retronet-4004),
[retronet-8008](https://github.com/retronet-labs/retronet-8008),
[retronet-8080](https://github.com/retronet-labs/retronet-8080) e, per i
programmi `.COM` didattici, da
[retronet-cpm](https://github.com/retronet-labs/retronet-cpm). Il backend
`i8086` genera anche **boot sector** avviabili da
[retronet-pc](https://github.com/retronet-labs/retronet-pc) (vedi
`examples/i8086-bootok.asm` ed `examples/i8086-echo.asm` e
[docs/arch-i8086.md](docs/arch-i8086.md)).

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
- Direttive: `.include "file.asm"` inserisce sorgenti locali, `.org <addr>`
  posiziona il codice, `.orgbase <addr>` cambia il PC logico senza padding,
  `.com` e' alias di `.orgbase 0x0100`, `.byte v1, v2, ...` emette dati in
  ROM (anche stringhe: `.byte "ciao", 0`, con escape `\n \r \t \0 \\ \"`).
- Arresto: `halt: JUN halt` (i4004) · `halt: JMP halt`/`HLT` (i8008).

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
| `i8080-hello.asm` (i8080) | stampa `HI` sulla porta terminale convenzionale `1` |
| `i8008-demo.asm` (i8008) | istruzioni 8008 a 1 byte senza operandi |
| `i8008-loop.asm` (i8008) | loop 8008: somma 5+4+3+2+1 = 15 (`LBI`/`ADB`/`DCB`/`JFZ`) |
| `i8008-sub.asm` (i8008) | subroutine 8008 `CAL`/`RET`: raddoppia 9 → 18 |
| `i8008-calc.asm` (i8008) | calcolatrice binaria a una cifra, 4 operatori, I/O terminale (`6*7=`→`42`) |
| `i8008-calc-multi.asm` (i8008) | calcolatrice binaria multi-cifra 0–255 (`12*12=`→`144`) |
| `i8086-bootok.asm` (i8086) | boot sector: messaggio via `INT 10h` su retronet-pc |
| `i8086-echo.asm` (i8086) | boot sector: eco dei tasti via `INT 16h`/`INT 10h` |
| `i8086-memdemo.asm` (i8086) | boot sector: stampa leggendo con `[msg+bx]` (operandi in memoria) |

I `*-bcd`/`multicifra`, assemblati, producono **gli stessi byte** delle ROM di
esempio di retronet-4004 (`testdata/`); gli `i8008-*` girano su retronet-8008 e
il suo `-disasm` ne ri-stampa gli identici mnemonici: è la validazione incrociata
assembler↔emulatore.

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
| Backend Intel 8008 | `arch/i8008`     | [docs/arch-i8008.md](docs/arch-i8008.md) · guida d'uso: [docs/guida-i8008.md](docs/guida-i8008.md) |
| Backend Intel 8086/8088 | `arch/i8086` | [docs/arch-i8086.md](docs/arch-i8086.md) |

Lexer, parser ed emitter sono **indipendenti dall'architettura**: ogni backend
vive in un pacchetto `arch/` che implementa `arch.Arch`.

---

## Roadmap

- [x] Backend `i4004`: tabella istruzioni, dimensionamento, codifica
- [x] Lexer, parser, symbol table, emitter a due passate
- [x] CLI `build`
- [x] Esempi + validazione contro le ROM golden di retronet-4004
- [x] Backend `i8008` (set completo, validato su retronet-8008)
- [x] Direttive `.org` (page alignment) e `.byte` (dati in ROM)
- [x] Direttiva `.equ` (costanti simboliche)
- [x] Backend `i8080`
- [x] Esempio end-to-end `i8080` -> `.COM` -> `retronet-cpm`
- [x] Direttiva `.com`/`.orgbase` per programmi CP/M `.COM` senza padding
- [x] Direttiva `.include` per librerie assembly locali
- [x] Backend `i8086` (registri/immediati, ALU, controllo, stringa) + boot sector
- [x] `i8086`: operandi in **memoria** (ModR/M 16 bit, `[bx+si]`/`[msg+bx]`/…) e `LEA`
- [x] Stringhe in `.byte "..."` (con escape) per i messaggi dei boot sector
- [ ] `i8086`: override di segmento (`[es:...]`)

---

## Licenza

MIT.
