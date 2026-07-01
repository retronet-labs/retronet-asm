# Backend i8086

Il backend `i8086` assembla l'Intel **8086/8088** in real mode. È pensato per
programmi a registri e per i **boot sector** dell'IBM PC/XT emulato da
[retronet-pc](https://github.com/retronet-labs/retronet-pc).

## Cosa supporta

- **Registri**: a 8 bit (`AL CL DL BL AH CH DH BH`), a 16 bit
  (`AX CX DX BX SP BP SI DI`) e di segmento (`ES CS SS DS`).
- **MOV** tra registri, da/verso segmento, e immediato → registro.
- **Aritmetico-logiche** `ADD OR ADC SBB AND SUB XOR CMP` e `TEST`, nelle forme
  registro-registro, accumulatore-immediato (forma corta) e registro-immediato.
- **INC/DEC**, **PUSH/POP** (registri e segmenti), **XCHG**.
- **NEG NOT MUL IMUL DIV IDIV** su registro.
- **Shift/rotate** `ROL ROR RCL RCR SHL SAL SHR SAR` per `1` o per `CL`.
- **Controllo di flusso**: `JMP` (vicino, o `JMP SHORT`), `CALL`, `RET`/`RETF`
  (con o senza immediato), tutti i salti condizionati `Jcc`, `LOOP*`, `JCXZ`.
- **Interrupt**: `INT n`, `INT3`, `INTO`, `IRET`.
- **Stringa**: `MOVSB/W STOSB/W LODSB/W SCASB/W CMPSB/W` e i prefissi
  `REP REPE REPNE LOCK` (su riga propria).
- **Flag e varie**: `CLC STC CLI STI CLD STD CMC`, `CBW CWD`, `SAHF LAHF`,
  `PUSHF POPF`, `IN/OUT` (porta immediata o `DX`), `NOP HLT WAIT XLAT`,
  aggiustamenti BCD `DAA DAS AAA AAS AAM AAD`.

### Operandi in memoria

Sono supportati gli operandi in **memoria** con l'indirizzamento a 16 bit
dell'8086, in MOV, nel blocco aritmetico-logico, TEST, INC/DEC, NEG/NOT/MUL/DIV,
shift, PUSH/POP e LEA:

- combinazioni base+indice: `[bx+si]`, `[bx+di]`, `[bp+si]`, `[bp+di]`;
- registro singolo: `[bx]`, `[bp]`, `[si]`, `[di]`;
- con spiazzamento, numerico o simbolico: `[bx+0x10]`, `[bp-2]`, `[msg+bx]`;
- indirizzo diretto: `[0x1234]`, `[msg]`.

La dimensione dello spiazzamento (disp8/disp16) si sceglie dalla sintassi: un
letterale che entra in 8 bit usa disp8, un simbolo usa sempre disp16 (così Size ed
Encode concordano senza risolvere le label). Per le forme **memoria-immediato**,
ambigue sulla larghezza, serve uno specificatore `byte`/`word`:

```asm
        mov   al, [bx+si]      ; r <- m
        mov   [di], cx         ; m <- r
        add   ax, [msg+bx]     ; indicizzato con spiazzamento simbolico
        mov   byte [bx], 0     ; memoria-immediato: serve byte/word
        inc   word [count]
        lea   si, [bx+di]      ; carica l'indirizzo effettivo
```

Esempio completo: [`examples/i8086/04-memoria-indicizzata.asm`](../examples/i8086/04-memoria-indicizzata.asm)
(stampa un messaggio leggendolo con `[msg+bx]`). Non sono ancora supportati gli
**override di segmento** (`[es:...]`).

## Boot sector

Un boot sector è un settore da 512 byte caricato dal BIOS a `0000:7C00`, con la
firma `0x55 0xAA` agli offset 510-511. Si scrive così:

```asm
.arch i8086
.orgbase 0x7C00          ; le label risolvono all'indirizzo di caricamento
        mov si, msg
        ...
msg:    .byte "Ciao", 0
        .org 0x7DFE       ; riempi (con 0) fino all'offset 510
        .byte 0x55, 0xAA  ; firma di boot
```

- `.orgbase 0x7C00` imposta la base logica degli indirizzi **senza** aggiungere
  padding all'inizio del file: l'output parte dal primo byte di codice, ma
  `mov si, msg` riceve l'indirizzo assoluto `0x7C00 + offset`.
- `.org 0x7DFE` riempie di zeri fino all'offset 510.
- L'immagine prodotta è di 512 byte; `retronet-pc` la accetta come floppy
  riempiendola al formato standard.

Il backend `i8086` non e' limitato ai boot sector: puo' assemblare anche normali
frammenti real mode o immagini caricate da un altro loader. Gli esempi in
`examples/i8086/` sono boot sector perche' sono eseguibili direttamente da
`retronet-pc` senza introdurre un formato DOS `.COM` o un loader separato.

Esempio completo: [`examples/i8086/02-stampa-stringa.asm`](../examples/i8086/02-stampa-stringa.asm)
(stampa un messaggio via `INT 10h`) ed [`examples/i8086/03-echo-tastiera.asm`](../examples/i8086/03-echo-tastiera.asm)
(eco dei tasti via `INT 16h`/`INT 10h`).

```bash
go run ./cmd/retronet-asm build examples/i8086/02-stampa-stringa.asm -o bootok.rom
# poi, in retronet-pc:
go run ./cmd/retronet-pc -bios <BIOS.ROM> -floppy bootok.rom
```

## Note di codifica

- I salti condizionati e `LOOP/JCXZ` usano sempre uno spiazzamento a 8 bit
  (`rel8`): il bersaglio dev'essere a portata (−128..+127). `JMP` senza `SHORT`
  usa `rel16` (3 byte) per raggiungere qualsiasi distanza.
- Per l'accumulatore (`AL`/`AX`) le ALU-immediato usano la forma corta (es.
  `CMP AL, 0` → `3C 00`); per gli altri registri la forma di gruppo `80/81`.
