# 09 - REP MOVSB

## Obiettivo
Usare le istruzioni stringa dell'8086 per copiare un blocco di memoria. `REP`
ripete `MOVSB` per il numero di byte indicato in `CX`.

## Lettura Del Codice
`DS:SI` punta alla sorgente e `ES:DI` al buffer di destinazione. `CX=6` indica
la lunghezza di `COPIA!`. Il prefisso `REP` precede `MOVSB`; al termine, `DI`
punta al byte successivo, dove viene scritto il terminatore zero.

## Stato Atteso
Il buffer contiene la copia della stringa e viene stampato con la subroutine
`print`. L'esempio chiarisce il ruolo implicito di `SI`, `DI`, `CX`, `DS` ed
`ES` nelle istruzioni stringa.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/09-rep-movsb.asm -o repmovsb.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy repmovsb.rom -steps 12000000
```
