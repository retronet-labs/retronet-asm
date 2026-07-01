# 14 - Calcolatrice Con Due Operandi Multicifra

## Obiettivo
Combinare input multicifra, deposito in RAM e aritmetica BCD su due operandi.
Il programma gestisce addizione e sottrazione.

## Lettura Del Codice
Il primo operando viene letto in `R0..R3`, poi copiato in RAM in ordine
little-endian. Dopo l'operatore viene letto il secondo operando. Le routine di
calcolo riusano il modello a puntatori paralleli gia' introdotto negli esempi
di somma e sottrazione multicifra.

## Stato Atteso
Input `47+58=` produce `105`; input `52-27=` produce `25`. La RAM diventa la
struttura dati principale, mentre i registri servono da puntatori e contatori.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/14-calcolatrice-due.asm -o calc-due.rom
printf '47+58=' | go run ../go-4004/cmd/retronet-4004 -io calc-due.rom
```
