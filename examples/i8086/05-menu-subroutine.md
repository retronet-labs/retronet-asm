# 05 - Menu Con Subroutine

## Obiettivo
Combinare input BIOS, confronti, salti condizionali e subroutine di stampa in un
boot sector piu' strutturato.

## Lettura Del Codice
Il programma stampa un menu, legge un tasto con `INT 16h` e confronta `AL` con
i codici ASCII di `1` e `2`. La subroutine `print` riceve in `SI` l'indirizzo di
una stringa terminata da zero e usa `LODSB` piu' `INT 10h` per stamparla. `CALL`
e `RET` usano lo stack reale dell'8086.

## Stato Atteso
Premendo `1` viene stampato `Hai scelto SOMMA`; premendo `2` viene stampato
`Fine`; altri tasti ristampano un messaggio di errore e ripetono la lettura.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/05-menu-subroutine.asm -o menu.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy menu.rom -keys "1" -steps 14000000
```
