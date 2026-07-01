# retronet-asm

Assembler modulare e multi-architettura dell'ecosistema **RetroNet**. Traduce un
sorgente testuale `.asm` in una ROM binaria `.rom`, risolvendo le label
automaticamente — niente più programmi scritti come array di byte contati a mano.

Architetture supportate: **Intel 4004** (`i4004`), **MOS/NMOS 6502** (`i6502`),
**Intel 8008** (`i8008`), **Intel 8080** (`i8080`) e **Intel 8086/8088**
(`i8086`).

Le ROM prodotte sono eseguibili dagli emulatori
[retronet-4004](https://github.com/retronet-labs/retronet-4004),
[retronet-6502](https://github.com/retronet-labs/retronet-6502),
[retronet-8008](https://github.com/retronet-labs/retronet-8008),
[retronet-8080](https://github.com/retronet-labs/retronet-8080) e, per i
programmi `.COM` didattici, da
[retronet-cpm](https://github.com/retronet-labs/retronet-cpm). Il backend
`i8086` genera anche **boot sector** avviabili da
[retronet-pc](https://github.com/retronet-labs/retronet-pc) (vedi
`examples/i8086/02-stampa-stringa.asm` ed `examples/i8086/03-echo-tastiera.asm` e
[docs/arch-i8086.md](docs/arch-i8086.md)).

---

## Quick start

```bash
# assembla un sorgente .asm in una .rom
go run ./cmd/retronet-asm build examples/i4004/05-somma-multicifra.asm -o out.rom

# (senza -o, l'output prende il nome dell'input: examples/i4004/05-somma-multicifra.rom)
go run ./cmd/retronet-asm build examples/i4004/05-somma-multicifra.asm

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
  ROM (anche stringhe: `.byte "ciao", 0`, con escape `\n \r \t \0 \\ \"`),
  `.word v1, v2, ...` emette parole little-endian risolte in seconda passata.
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

In [`examples/`](examples/) gli esempi sono divisi per architettura. Ogni
programma ha un `.asm` assemblabile e un `.md` didattico con spiegazione,
comandi e stato atteso.

| Cartella | Contenuto |
|----------|-----------|
| [`examples/i4004`](examples/i4004) | serie completa: output base, aritmetica, BCD e calcolatrici |
| [`examples/i8008`](examples/i8008) | cinque esempi: istruzioni base, loop, subroutine e calcolatrici binarie |
| [`examples/i8080`](examples/i8080) | dieci esempi: I/O, registri, loop, memoria, stack, 16 bit e CP/M `.COM` |
| [`examples/i6502`](examples/i6502) | vettore reset, branch, decimal mode, addressing e stack |
| [`examples/i8086`](examples/i8086) | dieci boot sector PC con BIOS, tastiera, stack, video, memoria e stringhe |

I programmi BCD 4004 migrati sotto `examples/i4004/` continuano a produrre gli
stessi byte golden delle ROM storiche. Gli esempi 8008 restano allineati al
disassembler di `retronet-8008`; i 6502 includono vettori hardware; gli 8086
sono boot sector da 512 byte con firma `55 AA`.

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
| Backend MOS/NMOS 6502 | `arch/i6502` | [docs/arch-i6502.md](docs/arch-i6502.md) |
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
- [x] Direttiva `.word` e backend `i6502` con sintassi MOS standard
- [ ] `i8086`: override di segmento (`[es:...]`)

---

## Licenza

MIT.
