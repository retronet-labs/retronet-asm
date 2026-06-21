# Guida all'uso dell'Intel 8008

Guida pratica per scrivere, assemblare ed eseguire programmi per l'**Intel 8008**
con `retronet-asm` (l'assembler) e
[`retronet-8008`](https://github.com/retronet-labs/retronet-8008) (l'emulatore).
Per la tabella completa dell'ISA e delle codifiche vedi [`arch-i8008.md`](arch-i8008.md);
per la sintassi `.asm` vedi [`sintassi-asm.md`](sintassi-asm.md).

---

## 1. Cosa serve

Due strumenti (Go 1.26), uno per repo:

```bash
# nel repo retronet-asm
go build -o retronet-asm ./cmd/retronet-asm
# nel repo retronet-8008
go build -o retronet-8008 ./cmd/retronet-8008
```

Il flusso è sempre: scrivi `.asm` → assembli in `.rom` → esegui la `.rom`.

---

## 2. L'8008 in breve

- CPU **a 8 bit** (l'erede del 4004 a 4 bit). Lavora con byte 0–255.
- **Registri**: `A` (accumulatore, dove avvengono i calcoli) e `B C D E H L`
  (generali). `H` e `L` formano l'indirizzo a 14 bit per accedere alla **memoria**
  (lo pseudo-registro `M` = byte puntato da `HL`).
- **Flag**: `C` (carry/borrow), `Z` (zero), `S` (segno), `P` (parità).
- **Stack** hardware a 8 livelli per `CAL`/`RET` (il PC corrente è in cima): 7
  livelli utili di annidamento.
- **Memoria** fino a 16 KB; **I/O** separato: 8 porte input, 24 output.
- A differenza dell'8080, l'8008 **non ha `DAA`**: l'aritmetica decimale (BCD) va
  fatta a mano, quindi in pratica si lavora in **binario**.

---

## 3. Scrivere un programma `.asm`

Prima riga: `.arch i8008`. I **mnemonici codificano i registri nel nome** (come
il disassembler dell'emulatore):

| Categoria | Esempi | Significato |
|-----------|--------|-------------|
| Move | `LAB`, `LMA`, `LAM` | A←B, M←A (scrive in memoria), A←M (legge) |
| Immediato | `LAI 0x2A`, `LBI 5` | carica un valore in A / B |
| ALU registro | `ADB`, `SUC`, `CPM` | A+=B, A-=C, confronta A con M |
| ALU immediato | `ADI 1`, `SUI 10`, `CPI 100` | A+=1, A-=10, confronta A con 100 |
| Inc/Dec | `INB`, `DCC` | B++, C-- (solo B–L) |
| Salti | `JMP et`, `JTZ et`, `JFC et` | salta / se Zero / se NON Carry |
| Subroutine | `CAL et`, `RET` | chiama / ritorna |
| I/O | `INP 0`, `OUT 8` | A←porta 0, porta 8←A |
| Stop | `HLT` | ferma la CPU |

Condizioni dei salti: `C Z S P`. Prefisso `T` = "se vero", `F` = "se falso":
`JTZ`/`JFZ`, `JTC`/`JFC`, ecc. Stessa cosa per call (`CTZ`…) e return (`RTZ`…).

**Direttive** (vedi `sintassi-asm.md`):

```asm
.arch i8008          ; sceglie l'architettura
.equ COUNT 5         ; costante simbolica (usabile come numero o indirizzo)
.org 0x100           ; posiziona il codice a un indirizzo
tabella: .byte 1, 2, 3   ; dati letterali in ROM
```

---

## 4. Assemblare ed eseguire

```bash
# assembla
retronet-asm build programma.asm -o programma.rom

# esegui (mostra registri/flag/PC finali)
retronet-8008 -bin programma.rom

# disassembla N istruzioni (verifica la codifica)
retronet-8008 -bin programma.rom -disasm 8

# traccia ogni istruzione
retronet-8008 -bin programma.rom -trace -steps 100
```

> L'estensione `.rom` (prodotta da retronet-asm) si carica con `-bin`. Il nome
> `-rom` nell'emulatore serve invece a caricare ROM di profilo (`-rom nome=file`).

---

## 5. Input/Output: il terminale

L'8008 separa I/O dalla memoria. L'emulatore collega un **terminale ASCII** alla
porta **0 (input)** e **8 (output)**:

- `INP 0` legge il prossimo carattere in `A` (consuma un byte della coda);
- `OUT 8` stampa il carattere in `A`.

Si accodano i tasti con `-terminal-input` (abilita da solo il terminale):

```bash
retronet-8008 -bin calc.rom -terminal-input "12*12=" -steps 2000000
# -> 144
```

I caratteri sono ASCII: `'0'`=0x30 … `'9'`=0x39, `'+'`=0x2B, `'='`=0x3D, ecc.
Per stampare la cifra `d` (0–9) si fa `ADI 0x30` (+`'0'`) e poi `OUT 8`.

---

## 6. Esempi commentati

Tutti in [`examples/`](../examples/), assemblabili ed eseguibili come sopra.

- **`i8008-loop.asm`** — somma `5+4+3+2+1` con un loop (`LBI`/`ADB`/`DCB`/`JFZ`) e
  la costante `.equ COUNT 5`. Esegui: `retronet-8008 -bin loop.rom` → `A=0x0F`.
- **`i8008-sub.asm`** — subroutine `CAL`/`RET` che raddoppia 9 → `A=0x12`.
- **`i8008-calc.asm`** — calcolatrice a **una cifra**, 4 operatori, I/O terminale:
  `-terminal-input "6*7="` → `42`.
- **`i8008-calc-multi.asm`** — calcolatrice **multi-cifra** (binaria, 0–255):
  `-terminal-input "12*12="` → `144`. `readnum` costruisce il numero (`×10 + cifra`),
  il display converte il binario in decimale a 3 cifre.

Esempio minimo completo:

```asm
.arch i8008
        LAI 3           ; A = 3
        ADI 4           ; A = 3 + 4 = 7
        ADI 0x30        ; A += '0'  ->  '7'  (valido solo se il risultato è 0..9)
        OUT 8           ; stampa "7"
        HLT
```

(Per risultati ≥ 10 serve la conversione binario→decimale, vedi `i8008-calc-multi.asm`.)

---

## 7. Dove approfondire

- [`arch-i8008.md`](arch-i8008.md) — tabella ISA, `kind`, codifiche.
- [`sintassi-asm.md`](sintassi-asm.md) — sintassi `.asm` e direttive.
- Repo emulatore `retronet-8008` — profili macchina, front panel, trace JSON,
  conformance: opzioni avanzate della CLI (`-profiles`, `-panel`, `-io-trace`…).
