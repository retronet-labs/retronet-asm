# Backend i6502

Il backend `i6502` assembla il set **documentato** del MOS/NMOS 6502. Usa la
sintassi MOS standard:

```asm
.arch i6502
.orgbase $8000

start:
        LDA #$01
        STA $0200
        LDA ($44),Y
        BNE start

        .org $FFFC
        .word start
```

## Numeri E Direttive

Sono accettati numeri decimali, `0x` esadecimali, `$` esadecimali e `%` binari:

```asm
LDA #10
LDA #0x0A
LDA #$0A
LDA #%00001010
```

`.word` emette parole little-endian e risolve label/costanti in seconda passata.
E' utile per vettori 6502:

```asm
.org $FFFC
.word start
```

## Addressing Mode

| Sintassi | Modo |
|----------|------|
| `LDA #$01` | immediate |
| `ASL A` | accumulator |
| `LDA $20` | zero page |
| `LDA $20,X` | zero page,X |
| `LDA $2000` | absolute |
| `LDA $2000,Y` | absolute,Y |
| `LDA ($44,X)` | indexed indirect |
| `LDA ($44),Y` | indirect indexed |
| `JMP ($12FF)` | indirect |
| `BNE label` | relative |

## Zero Page Vs Absolute

La dimensione deve essere stabile tra prima e seconda passata:

- letterale `0..$FF` usa zero page se la forma esiste;
- letterale `>$FF` usa absolute;
- simbolo usa absolute per default;
- `<label` forza zero page;
- `#<label` e `#>label` emettono byte basso/alto immediato.

Esempio:

```asm
ptr:    .byte 0
        LDA ptr      ; absolute, per non cambiare dimensione tra passate
        LDA <ptr     ; zero page forzata
        LDA #>start  ; byte alto dell'indirizzo di start
```

## Scope

Il backend rifiuta gli opcode illegali/non documentati e le forme non valide
(`STX $20,X`, `JMP #$12`, ecc.). L'obiettivo e' restare speculare al core
`retronet-6502`, che implementa lo stesso perimetro NMOS documentato.
