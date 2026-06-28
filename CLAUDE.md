# CLAUDE.md — retronet-asm

Assembler **modulare e multi-architettura** dell'ecosistema RetroNet: da un
sorgente `.asm` produce una ROM binaria risolvendo le label. Backend supportati:
**i4004**, **i8008**, **i8080**, **i8086/8088**. Panoramica utente:
[README.md](README.md); sintassi: [docs/sintassi-asm.md](docs/sintassi-asm.md).

## Setup su una macchina nuova (handoff)

Repo Go autonomo, **nessuna dipendenza esterna** (la pipeline è tutta locale):
un clone pulito compila e testa subito.

```sh
go test ./...                                   # tutti i pacchetti verdi
go build -o retronet-asm ./cmd/retronet-asm     # binario CLI
go run ./cmd/retronet-asm build <file.asm> -o <out.rom>
```

La validazione incrociata (le ROM prodotte girano sugli emulatori) usa i repo
sibling `retronet-4004`/`8008`/`8080`/`pc`, ma **non** servono per build/test di
questo repo. Esempi i8086 → boot sector per `retronet-pc`.

## Comandi

- Test: `go test ./...` ; Formattazione: `gofmt -w .` ; Analisi: `go vet ./...`
- Build di un esempio: `go run ./cmd/retronet-asm build examples/i8086-bootok.asm -o bootok.rom`
- Senza `-o`, l'output prende il nome dell'input (`.rom`).

## Architettura

Pipeline a stadi **indipendenti dall'architettura**; l'ISA sta dietro
`arch.Arch` (`Size` senza resolver / `Encode` con resolver — servono **due
passate** perché un salto può puntare a una label più avanti):

```
.asm → internal/lexer → internal/parser → internal/emitter (2 passate) → .rom
                                              └── interroga arch/<cpu>
```

- `arch/` — un pacchetto per CPU (`i4004`, `i8008`, `i8080`, `i8086`), registrato
  nella mappa `arches` della CLI; si sceglie con `.arch <nome>` (default `i4004`).
- `internal/lexer` — token, incl. **String** (`.byte "..."` con escape
  `\n \r \t \0 \\ \"`) e **Mem** (operandi `[...]`, spazi rimossi).
- `internal/parser` — token → `[]Stmt`; `internal/symbols` — tabella label;
  `internal/emitter` — driver a due passate.
- Direttive: `.arch`, `.org` (page align), `.orgbase`/`.com` (PC logico senza
  padding), `.byte` (dati/stringhe), `.equ` (costanti), `.include` (sorgenti).

## Backend i8086 (il più recente)

Doc: [docs/arch-i8086.md](docs/arch-i8086.md). File: `arch/i8086/i8086.go`
(handler per famiglia) + `arch/i8086/mem.go` (operandi in memoria).

- Registri 8/16 bit + segmento; MOV, blocco ALU+TEST, INC/DEC, NEG/NOT/MUL/DIV,
  shift/rotate, PUSH/POP, XCHG, LEA, controllo (`JMP`/`Jcc`/`LOOP`/`CALL`/`RET`),
  INT/IRET, stringa + prefissi, flag misc, IN/OUT, BCD.
- **Operandi in memoria** (ModR/M a 16 bit): `[bx+si]`, `[bp]`, `[bx+0x10]`,
  `[msg+bx]` (disp simbolico), `[0x1234]`/`[msg]` (diretto). Forme
  memoria-immediato richiedono lo specificatore `byte`/`word`.
- **Regola chiave** (`Size`==`Encode`): la dimensione del disp si sceglie dalla
  **sintassi**, non dal valore risolto — un letterale ≤8 bit → disp8, un
  **simbolo → sempre disp16**. Così `Size` (resolver nil) ed `Encode` concordano
  e gli indirizzi non slittano tra le due passate. `Size` delega a `Encode(…,nil)`.
- **Boot sector**: `.orgbase 0x7C00` + codice + `.org 0x7DFE` + `.byte 0x55,0xAA`
  → immagine da 512 byte avviabile in `retronet-pc` come `-floppy`. Esempi:
  `i8086-bootok.asm`, `i8086-echo.asm`, `i8086-memdemo.asm` — validati e2e.
- Mancante: **override di segmento** (`[es:...]`).

## Convenzioni

- Commit atomici: `feat(asm|i8086|lexer|parser):`, `test(...)`, `docs(...)`.
- Ogni unità = codice + test + (se serve) esempio, con `go test ./...` verde
  prima del commit. Ogni parte essenziale ha un doc in `docs/`.

## Stato

Pipeline completa, testata e documentata; quattro backend (`i4004`/`i8008`/
`i8080`/`i8086`) pubblicati. Validazione incrociata attiva: ROM 4004 == golden di
retronet-4004; programmi 8008/8080 girano sui rispettivi emulatori e ne ri-stampa
gli mnemonici il `-disasm`; i boot sector i8086 bootano in retronet-pc.

Note utili da ricordare (4004): sottrazione BCD robusta = complemento a 10
(`M + comp9(S) + 1`); `×`/`÷` per addizioni/sottrazioni ripetute (O(valore));
attenzione al vincolo "stessa pagina" su `JCN`/`ISZ` oltre `0x100`.

Prossimi passi: override di segmento nell'i8086; eventuali nuovi backend con lo
stesso schema (pacchetto `arch/` + riga nella mappa `arches`).
