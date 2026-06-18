# Il lexer (tokenizer)

Il **lexer** è il primo stadio dell'assembler: trasforma il **testo** del
sorgente `.asm` in una sequenza di **token**, le unità minime con cui lavorerà
il parser. È l'unico stadio che "legge i caratteri"; tutto il resto della
pipeline ragiona su token.

```
sorgente .asm ──▶ [LEXER] ──▶ token ──▶ parser ──▶ emitter ──▶ .rom
```

File: [`internal/lexer/lexer.go`](../internal/lexer/lexer.go)

---

## Cosa fa (e cosa NON fa)

**Fa:**

- spezza il testo in token;
- salta spazi, tabulazioni e commenti;
- emette un token di fine riga (`Newline`) come separatore di istruzioni;
- tiene il **numero di riga** di ogni token, per messaggi d'errore precisi.

**Non fa** (di proposito):

- non sa cosa sia un *mnemonico*, un *registro* o una *label*. Per il lexer
  `LDM`, `R4` e `loop` sono tutti la stessa cosa: un **identificatore** (`Ident`).
  La distinzione la fanno il parser e l'architettura (`arch/i4004`).
- non valida i numeri (controlla solo che inizino con una cifra), non risolve
  niente, non conosce l'ISA.

Questa separazione netta è ciò che tiene il lexer **indipendente
dall'architettura**: lo stesso lexer servirà per 4004, 8008 e 8080.

---

## I tipi di token

| `Type`    | Significato | Esempio |
|-----------|-------------|---------|
| `Ident`   | identificatore: inizia con lettera o `_` | `LDM`, `R4`, `loop` |
| `Number`  | numero: inizia con una cifra | `12`, `0x0C`, `255` |
| `Colon`   | due punti (definizione di label) | `:` |
| `Comma`   | virgola (separatore di operandi) | `,` |
| `Newline` | fine riga (separa gli statement) | a-capo |
| `EOF`     | fine input | — |

Ogni token è:

```go
type Token struct {
	Type Type
	Text string // testo originale (es. "0x0C", "loop")
	Line int    // riga 1-based, per gli errori
}
```

---

## Le regole di scansione

Il lexer scorre la stringa carattere per carattere e decide in base al primo
carattere del prossimo token:

| Carattere iniziale | Azione |
|--------------------|--------|
| spazio, tab, `\r`  | ignorato |
| `;`                | commento: salta fino a fine riga |
| `\n`               | emette `Newline`, incrementa il contatore di riga |
| `:`                | emette `Colon` |
| `,`                | emette `Comma` |
| cifra `0-9`        | legge un `Number` (prende tutti gli alfanumerici di seguito) |
| lettera o `_`      | legge un `Ident` (lettere, cifre, `_`) |
| qualsiasi altro    | **errore**: "carattere inatteso", con la riga |

### Perché il numero legge "tutti gli alfanumerici"

Per gestire l'esadecimale `0x0C`: la scansione parte dalla cifra `0` e continua
finché trova caratteri alfanumerici, catturando `0x0C` come un unico `Number`.
La validazione vera del valore (decimale vs esadecimale, range) avviene più
avanti, in `arch/i4004` quando l'operando viene codificato.

### Perché `Ident` e `Number` si distinguono dal primo carattere

`R4` inizia con una lettera → `Ident`. `12` inizia con una cifra → `Number`.
Così un registro non viene mai scambiato per un numero, e il parser riceve già
i due casi separati.

---

## Esempio

Sorgente:

```asm
loop: ADD R1, 0x0C   ; somma l'addendo
```

Token prodotti:

```
Ident("loop")  Colon(":")  Ident("ADD")  Ident("R1")
Comma(",")     Number("0x0C")  Newline   EOF
```

Il commento `; somma l'addendo` sparisce; restano solo i token "utili".

---

## `Newline` è un token, non viene scartato

L'assembler ha **una istruzione per riga** (con l'eccezione `label: istruzione`).
Il parser ha quindi bisogno di sapere dove finisce uno statement: il token
`Newline` è quel confine. Una riga di solo commento produce comunque il suo
`Newline` (la riga "esiste", anche se vuota di codice).

---

## I numeri di riga viaggiano fino in fondo

Ogni token porta la sua `Line`. Parser ed emitter la propagano negli errori,
così un problema segnalato dall'encoder (es. *"registro non valido"*) può dire
**a quale riga** del sorgente si riferisce. È il motivo per cui il lexer traccia
le righe fin da subito.

---

## API

```go
func Tokenize(src string) ([]Token, error)
```

Restituisce l'intera lista di token (terminata da `EOF`) oppure un errore al
primo carattere inatteso. Non ha stato: una stringa entra, una slice di token
esce.

I test sono in [`internal/lexer/lexer_test.go`](../internal/lexer/lexer_test.go).
