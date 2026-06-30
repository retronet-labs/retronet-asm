# La sintassi `.asm`

Guida alla sintassi dei sorgenti assembly accettati da `retronet-asm` (arch
i4004). È volutamente minimale.

```
retronet-asm build programma.asm -o programma.rom
```

---

## Direttiva di architettura: `.arch`

La prima riga "di codice" (saltando righe vuote e commenti) può dichiarare per
quale CPU è scritto il programma:

```asm
.arch i4004
```

L'assembler usa quel nome per scegliere il backend. La direttiva **non produce
byte** (è solo metadato). Se assente, si assume `i4004`. Un nome non registrato
dà errore. Le architetture disponibili sono `i4004`, `i6502`, `i8008`, `i8080`
e `i8086`.

---

## Direttiva include: `.include`

`.include "file.asm"` inserisce un altro sorgente prima del parsing. Il path e'
relativo al file che contiene la direttiva; gli include annidati sono permessi,
ma devono restare dentro la directory del sorgente principale.

```asm
.arch i8080
.include "lib/cpm-bdos.asm"
.com
```

Non ci sono fetch di rete, macro o package manager: e' solo un modo esplicito
per riusare costanti e piccoli blocchi didattici versionati nel workspace.

---

## Direttiva di posizione: `.org`

`.org <indirizzo>` posiziona il codice che segue a un indirizzo preciso della ROM,
riempiendo il vuoto con `NOP` (`0x00`). Serve soprattutto a rispettare il vincolo
"stessa pagina" dei salti `JCN`/`ISZ`/`JIN` nei programmi oltre i 256 byte.

```asm
.org 0x100          ; il codice seguente parte a 0x100 (decimale: .org 256)
```

L'indirizzo deve essere ≥ alla posizione corrente (non si torna indietro) e dentro
lo spazio ROM (`0x000`–`0xFFF`). Dettagli ed esempi in [`org.md`](org.md).

---

## Direttiva dati: `.byte`

`.byte v1, v2, ...` emette byte letterali (0–255) direttamente nella ROM: serve
per tabelle e costanti di dati. Può avere una label, che punta all'indirizzo del
primo byte.

```asm
tabella: .byte 0x41, 0x42, 0x43, 0    ; quattro byte in ROM
```

I valori sono decimali o esadecimali, separati da virgole (o spazi). È
indipendente dall'architettura.

---

## Direttiva parole: `.word`

`.word v1, v2, ...` emette parole a 16 bit little-endian. A differenza di
`.byte`, i valori possono essere label o costanti risolte in seconda passata.
Serve soprattutto per vettori 6502 e tabelle di indirizzi.

```asm
.arch i6502
.orgbase $8000
start:  NOP
        .org $FFFC
        .word start
```

Sono accettati numeri decimali, `0x` esadecimali, `$` esadecimali e `%` binari.

---

## Direttiva costanti: `.equ`

`.equ NOME valore` dà un nome a un numero, usabile ovunque vada un valore: immediati
(`LDM NOME`, `LAI NOME`), porte, vettori `RST`, target di salto (`JMP NOME`). Il nome
entra nella symbol table come una label, quindi è usabile anche **prima** della sua
definizione.

```asm
.equ COUNT 5
.equ PORTA 1
        LBI COUNT       ; B = 5
        INP PORTA       ; legge la porta 1
```

Il valore è un numero (decimale o esadecimale). È indipendente dall'architettura.

---

## Regole generali

- **Una istruzione per riga.**
- **Commenti**: tutto ciò che segue `;` fino a fine riga è ignorato.
- **Righe vuote** (o di solo commento) sono ignorate.
- **Maiuscole/minuscole**: i mnemonici sono case-insensitive (`ldm` = `LDM`).
  Le **label sono case-sensitive** (`Loop` ≠ `loop`).

## Label

Una label si definisce con un nome seguito da `:` e marca l'indirizzo
dell'istruzione successiva.

```asm
loop:   ADD R1     ; label e istruzione sulla stessa riga
fine:              ; label da sola (punta all'istruzione che segue)
```

Si usa come operando nei salti (`JUN`, `JMS`, `JCN`, `ISZ`): l'assembler la
sostituisce con l'indirizzo, anche se è definita più avanti nel file.

## Operandi

- **Registri**: `R0`–`R15` (case-insensitive). Le coppie (`FIM`, `SRC`, `FIN`,
  `JIN`) si indicano col registro pari: `FIM R0, 0x35`.
- **Numeri**: decimali (`12`) o esadecimali (`0x0C`).
- **Separatore**: la **virgola tra operandi è opzionale** —
  `FIM R0, 0x35` e `FIM R0 0x35` sono equivalenti.
- **6502**: gli operandi MOS come `#$01`, `$0200,X`, `($44),Y`, `<label`,
  `#<label` e `#>label` sono preservati e interpretati dal backend `i6502`.

## Arresto (HALT)

Il 4004 non ha un'istruzione di stop. Per convenzione si termina con un salto su
se stessi, usando una label:

```asm
halt:   JUN halt
```

---

## Esempio completo

```asm
; moltiplicazione 3 x 4 tramite addizioni ripetute
        LDM 0
        DCL                 ; banco RAM 0
        FIM R0, 0x03        ; R1 = 3 (addendo)
        FIM R2, 0x00        ; indirizzo RAM 0
        SRC R2
        LDM 12              ; contatore = 16 - 4
        XCH R4
loop:   ADD R1              ; A += 3
        ISZ R4, loop        ; ripeti 4 volte
        WRM                 ; scrivi il risultato in RAM
halt:   JUN halt
```

Assemblato, produce una ROM eseguibile da `retronet-4004`:

```
retronet-asm build moltiplicazione.asm
retronet-4004 -trace -dump-ram moltiplicazione.rom
```

---

## Istruzioni e operandi attesi

Per la tabella completa (mnemonico → opcode → tipo di codifica) vedi
[`docs/arch-i4004.md`](arch-i4004.md). In sintesi, per numero di operandi:

| Operandi | Istruzioni |
|----------|-----------|
| nessuno  | NOP, CLB, CLC, IAC, CMC, CMA, RAL, RAR, TCC, DAC, TCS, STC, DAA, KBP, DCL, WRM, WMP, WRR, WPM, WR0–WR3, SBM, RDM, RDR, ADM, RD0–RD3 |
| 1 registro | INC, ADD, SUB, LD, XCH |
| 1 immediato (0–15) | LDM, BBL |
| 1 registro pari | SRC, FIN, JIN |
| 1 indirizzo/label | JUN, JMS |
| condizione + indirizzo/label | JCN |
| registro + indirizzo/label | ISZ |
| registro pari + dato (0–255) | FIM |
