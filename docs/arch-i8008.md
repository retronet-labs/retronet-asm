# Il backend `i8008`

`arch/i8008` implementa l'interfaccia `arch.Arch` per l'**Intel 8008**, allo
stesso modo di `arch/i4004`: lexer, parser ed emitter non cambiano, ricevono
l'architettura e la interrogano in due passate (`Size`, poi `Encode`). Si
seleziona con `.arch i8008` sulla prima riga del sorgente.

File: [`arch/i8008/i8008.go`](../arch/i8008/i8008.go) · interfaccia: [`arch-i4004.md`](arch-i4004.md)

---

## Convenzione dei mnemonici = il disassembler dell'emulatore

I mnemonici 8008 **codificano i registri nel nome** (non come operando): `LAB`
= carica A da B, `ADC` = A += C, `INB` = incrementa B. Solo immediati, indirizzi,
`RST` e porte hanno un operando.

La scelta chiave: i mnemonici coincidono **esattamente** con quelli prodotti dal
disassembler di [retronet-8008](https://github.com/retronet-labs/retronet-8008).
Così assembler ed emulatore sono speculari — un `.asm` assemblato qui, eseguito
con `retronet-8008 -disasm`, ri-stampa gli stessi mnemonici. È anche la
**validazione**: assemblo → eseguo sull'emulatore → verifico.

Per non scrivere a mano ~150 voci, la tabella `set` è **generata
programmaticamente** dai pattern di bit dell'ISA (gli stessi del decoder
dell'emulatore), quindi le due parti restano allineate per costruzione.

---

## Classificazione: `kind`

| `kind`   | Byte | Operandi | Esempi |
|----------|------|----------|--------|
| `simple` | 1    | 0        | `LAB`, `ADM`, `INB`, `RLC`, `RET`, `RFC`, `HLT` |
| `imm`    | 2    | 1 (0–255)        | `LAI`…`LMI`, `ADI`…`CPI` |
| `addr`   | 3    | 1 (addr/label)   | `JMP`, `CAL`, `JFZ`/`JTC`…, `CFc`/`CTc` |
| `rst`    | 1    | 1 (vettore 0–7)  | `RST` |
| `inp`    | 1    | 1 (porta 0–7)    | `INP` |
| `out`    | 1    | 1 (porta 8–31)   | `OUT` |

`kind.operands()` e `kind.size()` derivano arità e dimensione dal tipo.

---

## Le famiglie e le loro codifiche

I registri usano i codici 8008: `A=0, B=1, C=2, D=3, E=4, H=5, L=6, M=7`
(`M` = byte di memoria puntato da `HL`). Le condizioni: `C=0, Z=1, S=2, P=3`.

| Famiglia | Pattern di bit | Opcode base | Note |
|----------|----------------|-------------|------|
| Move `Lr1r2` | `11 DDD SSS` | — | `dst != src`; `dst == src` è `NOP`/`HLT` |
| Load imm. `LrI` | `00 DDD 110` | `0x06` | + byte immediato |
| ALU registro | `10 GGG SSS` | `0x80` | `AD,AC,SU,SB,ND,XR,OR,CP` × reg/`M` |
| ALU immediato | `00 GGG 100` | `0x04` | + byte immediato |
| `INr`/`DCr` | `00 RRR 000/001` | — | solo B–L (non A, non `M`) |
| Rotate | — | `02/0A/12/1A` | solo il flag Carry |
| `JMP`/`CAL` | `01 xxx 100/110` | `0x44`/`0x46` | + indirizzo 14 bit |
| Jump cond. | `01 0CC 000` (F) / `01 1CC 000` (T) | `0x40`/`0x60` | + indirizzo |
| Call cond. | `01 0CC 010` / `01 1CC 010` | `0x42`/`0x62` | + indirizzo |
| `RET` / cond. | `00 xxx 111` / `00 0CC 011`,`00 1CC 011` | `0x07` / `0x03`,`0x23` | — |
| `RST n` | `00 NNN 101` | `0x05` | vettore a pagina 0 (`n*8`) |
| `INP`/`OUT` | `01 MMMMM 1` | `0x41` | campo 5 bit: in 0–7, out 8–31 |
| `HLT` | `00 000 00X` (+ alias `0xFF`) | `0x00` | emesso come `0x00` |

### Indirizzi a 14 bit

L'8008 ha un bus indirizzi a **14 bit** (16 KB), diverso dai 12 bit + pagina del
4004. `Encode` emette **low byte, poi high byte mascherato a 6 bit**:

```
JMP/CAL/Jcc/Ccc  →  [ opcode , addr & 0xFF , (addr >> 8) & 0x3F ]
```

Le label passano per `parseAddr` (numero se inizia con una cifra, altrimenti
`resolve(nome)`), come nell'i4004.

### Cosa l'assembler non emette (di proposito)

Sono **alias ridondanti**, non istruzioni mancanti: le altre due codifiche di
`HLT` (`0x01`, `0xFF`), i move "stessa coppia" (`LAA`… = no-op, c'è solo `NOP`),
e gli alias di `JMP`/`CAL`. L'emulatore li *decodifica* tutti; l'assembler
genera la forma canonica.

---

## Validazione

Esempi in [`examples/i8008/`](../examples/i8008/), tutti verificati assemblando e poi
eseguendo la ROM su `retronet-8008`:

| File | Cosa mostra | Risultato |
|------|-------------|-----------|
| `01-demo.asm` | istruzioni a 1 byte senza operandi | round-trip `-disasm` identico |
| `02-loop.asm` | loop con `LBI`/`ADB`/`DCB`/`JFZ` (somma 5+4+3+2+1) | `A = 0x0F` |
| `03-subroutine.asm`  | subroutine `CAL`/`RET` (raddoppia 9) | `A = 0x12` |

Esempio di catena completa:

```console
$ retronet-asm build examples/i8008/02-loop.asm -o loop.rom
$ retronet-8008 -bin loop.rom -disasm 6
0000: 0E 05    LBI #0x05
0002: 06 00    LAI #0x00
0004: 81       ADB
0005: 09       DCB
0006: 48 04 00 JFZ 0x0004
0009: 00       HLT
$ retronet-8008 -bin loop.rom        # A=0x0F  Halted=true
```

I test del backend sono in [`arch/i8008/i8008_test.go`](../arch/i8008/i8008_test.go).
