# 06 - Sottrazione BCD A Cifra Singola

## Obiettivo
Calcolare `7 - 5` in BCD usando il modello del prestito del 4004. Il flag Carry
assume il significato inverso rispetto alla sottrazione scolastica: `1` indica
nessun prestito.

## Lettura Del Codice
`STC` imposta il Carry in ingresso. `TCS` prepara nell'accumulatore la costante
di complemento decimale, `SUB R1` sottrae il sottraendo e `ADD R2` aggiunge il
minuendo. `DAA` completa la correzione BCD.

## Stato Atteso
La RAM riceve `2`. Il programma e' volutamente a cifra singola per rendere
leggibile la sequenza `TCS`, `SUB`, `ADD`, `DAA`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/06-sottrazione-bcd.asm -o sub-bcd.rom
go run ../go-4004/cmd/retronet-4004 -dump-ram sub-bcd.rom
```
