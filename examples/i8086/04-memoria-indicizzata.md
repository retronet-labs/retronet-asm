# 04 - Memoria Indicizzata

## Obiettivo
Mostrare gli operandi in memoria del backend i8086, in particolare la forma
`[msg+bx]` con spiazzamento simbolico e registro base.

## Lettura Del Codice
`BX` parte da zero e viene usato come indice nella stringa. `mov al, [msg+bx]`
legge un byte da memoria senza usare `LODSB`; `inc bx` avanza al carattere
successivo. Il terminatore zero viene riconosciuto con `cmp al, 0`.

## Stato Atteso
Il boot sector stampa `Indexed memory OK!`. L'esempio e' utile per verificare
che l'assembler emetta ModR/M e displacement coerenti con le label.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/04-memoria-indicizzata.asm -o memdemo.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy memdemo.rom -steps 12000000
```
