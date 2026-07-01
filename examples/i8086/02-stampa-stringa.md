# 02 - Stampa Stringa BIOS

## Obiettivo
Stampare una stringa da boot sector usando il servizio teletype del BIOS:
`INT 10h` con `AH = 0Eh`.

## Lettura Del Codice
`xor ax, ax` e `mov ds, ax` fissano `DS=0`, coerente con l'indirizzo fisico del
boot sector. `SI` punta alla stringa terminata da zero. `LODSB` carica il byte
successivo in `AL`; `CMP AL, 0` e `JE halt` riconoscono il terminatore.

## Stato Atteso
Il video mostra `RETRONET-PC: BOOT OK!`. Il programma dimostra la relazione tra
indirizzi logici generati dall'assembler e servizi BIOS disponibili in real mode.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/02-stampa-stringa.asm -o bootok.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy bootok.rom -steps 12000000
```
