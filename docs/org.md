# La direttiva `.org`

> **Stato:** specifica della direttiva, **non ancora implementata** (vedi roadmap
> in `CLAUDE.md`, item 13). Questo documento descrive il comportamento previsto e
> serve da guida all'implementazione.

`.org <indirizzo>` ("origin") dice all'assembler **a quale indirizzo della ROM va
posizionato il codice che segue**. È una *direttiva*, non un'istruzione: non
genera un opcode, sposta il **contatore di posizione** (l'indirizzo a cui verrà
emesso il prossimo byte). Lo spazio lasciato vuoto viene riempito di `NOP`
(`0x00`).

```asm
.arch i4004
        LDM 0
        JUN start
.org 0x100          ; il codice seguente parte esatto a 0x100
start:  LDM 5
        ...
```

---

## Perché serve: il vincolo "stessa pagina"

Il 4004 divide la ROM in **pagine da 256 byte** (`0x000–0x0FF`, `0x100–0x1FF`, …).
Le istruzioni di salto **a 2 byte ma con indirizzo a 8 bit** — `JCN`, `ISZ`, `JIN`
— codificano solo il byte basso del target: la **pagina** (i 4 bit alti) è quella
del PC corrente. Quindi questi salti **non possono uscire dalla pagina** in cui si
trovano.

`JUN` e `JMS` invece hanno l'indirizzo completo a 12 bit e attraversano le pagine
senza problemi.

In un programma piccolo (< 256 byte) tutto sta in pagina 0 e non te ne accorgi
mai. Ma quando un programma cresce oltre i 256 byte, un `JCN`/`ISZ` può trovarsi
**a cavallo di un confine di pagina** rispetto al suo target, e saltare
nell'indirizzo sbagliato.

### Esempio reale

Nella `calcolatrice-completa.asm` (384 byte) un `JCN` a `0x0FC` aveva come target
una label finita a `0x100`. Avendo `JCN` solo il byte basso (`0x00`) e prendendo
la pagina dal PC corrente (pagina 0), il salto andava a `0x000` invece che a
`0x100` — riavviando il programma. Senza `.org` l'unico rimedio è **riorganizzare
il codice a mano** perché quel salto e il suo target cadano nella stessa pagina.

Con `.org` lo si risolve in modo esplicito: si allinea il blocco a un confine di
pagina noto.

```asm
.org 0x100          ; le subroutine partono a inizio pagina 1:
pdig:   ...          ; nessun salto interno attraversa più 0x100
clr2:   ...
```

---

## A cosa serve, in pratica

- **Allineare un loop o un blocco a inizio pagina**, così nessun salto interno
  (`JCN`/`ISZ`) attraversa un confine.
- **Tenere insieme un salto e il suo target** nella stessa pagina, in modo
  prevedibile invece che "sperando" nella disposizione automatica.
- **Mettere tabelle, vettori o handler a indirizzi fissi e noti** (utile più
  avanti, p.es. con backend che usano salti tabellari).

---

## Sintassi

```asm
.org 0x100          ; indirizzo esadecimale
.org 256            ; equivalente in decimale
```

- Un solo operando: l'indirizzo di destinazione (decimale o esadecimale, come gli
  altri numeri — vedi [`sintassi-asm.md`](sintassi-asm.md)).
- La direttiva **non produce byte propri**: produce solo il padding necessario per
  arrivare all'indirizzo richiesto.
- Si può usare più volte nello stesso file.

### Regole

- L'indirizzo deve essere **maggiore o uguale** alla posizione corrente: non si
  può "tornare indietro" (sovrascriverebbe codice già emesso). → **errore**.
- L'indirizzo deve stare nello spazio ROM del 4004 (`0x000`–`0xFFF`, 12 bit). Un
  valore fuori range → **errore**.
- Il riempimento usa `0x00`, che sul 4004 è `NOP`: se per errore l'esecuzione ci
  finisce dentro, scorre fino al codice successivo senza effetti.

---

## Esempi

### 1. Posizionare il codice a un indirizzo

```asm
.arch i4004
.org 0x010
start:  LDM 7
halt:   JUN halt
```

I byte `0x000`–`0x00F` sono padding (`NOP`), `start` è a `0x010`.

```
offset  byte   significato
0x000   00     NOP   ┐
…       00     NOP   │ padding generato da .org 0x010
0x00F   00     NOP   ┘
0x010   D7     LDM 7   <- start
0x011   40     JUN ┐
0x012   11     …   ┘ -> 0x011 (halt)
```

### 2. Allineare le subroutine a inizio pagina

```asm
.arch i4004
        ; --- programma principale (pagina 0) ---
        LDM 0
        JMS pdig
        JUN halt
halt:   JUN halt

.org 0x100          ; subroutine tutte in pagina 1
pdig:   XCH R7      ; i JCN/ISZ interni non rischiano di attraversare 0x100
        ...
        BBL 0
```

`JMS pdig` funziona perché `JMS` usa l'indirizzo completo a 12 bit; i salti
*interni* a `pdig` restano dentro pagina 1.

### 3. Più direttive `.org`

```asm
.org 0x100
tabella: ...

.org 0x200
codice:  ...
```

### 4. Errore: indirizzo all'indietro

```asm
        LDM 0       ; emesso a 0x000 (1 byte) -> posizione corrente 0x001
.org 0x000          ; ERRORE: 0x000 < 0x001, sovrapporrebbe codice
```

---

## Note di implementazione (per domani)

Si incastra nella pipeline esistente **senza toccare l'interfaccia `arch`**: la
direttiva riguarda gli *indirizzi*, non l'ISA (come `.arch`).

1. **Riconoscimento** (`internal/source` o parser): trattare `.org <numero>` come
   uno statement speciale, p.es. `Stmt` con un campo `Org *uint16` (analogo a come
   `.arch` è gestito come metadato).
2. **Passata 1 — indirizzi/Size** (`internal/emitter`): si tiene un contatore
   `pc`. A `.org N`:
   - se `N < pc` → errore ("indirizzo .org all'indietro");
   - se `N > 0xFFF` → errore ("fuori dallo spazio ROM");
   - altrimenti `pc = N`. Le label successive vengono registrate nella symbol
     table agli indirizzi basati su `N`.
3. **Passata 2 — Encode**: emettere `N - pc` byte di padding `0x00`, poi
   continuare a codificare normalmente. Il `pc` deve coincidere con quello della
   passata 1.

Vedi [`due-passate.md`](due-passate.md) per il funzionamento delle due passate e
della symbol table.

### Test da aggiungere

- `.org` in avanti: il padding ha la lunghezza giusta e le label successive hanno
  l'indirizzo atteso.
- una label dopo `.org` usata come target di `JUN` punta all'indirizzo corretto.
- `.org` all'indietro → errore.
- `.org` fuori range (`> 0xFFF`) → errore.
- `.org` multipli nello stesso file.
