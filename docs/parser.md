# Il parser (token → statement)

Il **parser** è il secondo stadio: prende i token del lexer e li raggruppa in
**statement**, uno per riga del sorgente. È ancora **indipendente
dall'architettura**: riconosce la *forma* di una riga, non valida l'ISA.

```
token ──▶ [PARSER] ──▶ []Stmt ──▶ emitter
```

File: [`internal/parser/parser.go`](../internal/parser/parser.go)

---

## La grammatica di una riga

```
riga = [ label ':' ] [ mnemonico [ operando { ','? operando } ] ]
```

Ogni riga può avere:

- **solo una label**: `loop:`
- **solo un'istruzione**: `ADD R1`
- **entrambe**: `loop: ADD R1`
- **niente** (riga vuota o di solo commento) → saltata

Le righe sono delimitate dai token `Newline`/`EOF` prodotti dal lexer.

---

## Lo statement

```go
type Stmt struct {
	Label string            // label definita qui (vuota se assente)
	Instr *arch.Instruction // istruzione (nil se la riga ha solo una label)
	Line  int               // riga sorgente (1-based)
}
```

Una riga `loop: ADD R1` produce **un** `Stmt` con sia `Label` (`"loop"`) sia
`Instr` (`ADD [R1]`): l'emitter, nella passata 1, registrerà la label
all'indirizzo corrente e poi dimensionerà l'istruzione che la segue.

---

## Tre regole importanti

1. **La label si riconosce dal `:`** — un `Ident` è una label solo se è
   *seguito da* `:`. Così in `ADD R1` la parola `ADD` non viene scambiata per
   una label: diventa il mnemonico.

2. **Il mnemonico è normalizzato in MAIUSCOLO** — `ldm 5` e `LDM 5` sono
   equivalenti (la tabella ISA usa chiavi maiuscole).

3. **Gli operandi restano col testo originale** — i registri li normalizza
   l'encoder, ma le **label sono case-sensitive**: `loop` definita e `loop`
   riferita devono coincidere, quindi il parser non tocca gli operandi.

Le **virgole tra operandi sono separatori opzionali**: `FIM R0, 0x35` e
`FIM R0 0x35` producono gli stessi due operandi `["R0", "0x35"]`. La virgola
è solo zucchero sintattico per leggibilità.

---

## Cosa NON fa il parser

Non conosce l'ISA: non sa se `ADD` voglia un operando, se `R99` sia un registro
valido, o se la label `loop` esista. Produce solo la **struttura**. Il controllo
di arità e la codifica spettano ad `arch` (vedi `docs/arch-i4004.md`), la
risoluzione delle label all'emitter (vedi `docs/due-passate.md`).

Questa separazione mantiene il parser riusabile per tutte le architetture.

---

## Esempio

Sorgente:

```asm
        LDM 0
loop:   ADD R1
        ISZ R4, loop
halt:   JUN halt
```

Statement prodotti:

| # | Label | Mnemonico | Operandi | Line |
|---|-------|-----------|----------|------|
| 0 | —     | `LDM`     | `[0]`        | 1 |
| 1 | `loop`| `ADD`     | `[R1]`       | 2 |
| 2 | —     | `ISZ`     | `[R4, loop]` | 3 |
| 3 | `halt`| `JUN`     | `[halt]`     | 4 |

---

## API

```go
func Parse(toks []lexer.Token) ([]Stmt, error)
```

Restituisce la lista di statement, oppure un errore (con numero di riga) su una
riga malformata — ad esempio una riga che inizia con un numero o una virgola.

I test sono in [`internal/parser/parser_test.go`](../internal/parser/parser_test.go).
