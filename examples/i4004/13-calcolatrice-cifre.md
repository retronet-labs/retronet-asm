# 13 - Calcolatrice: Input Multicifra

## Obiettivo
Leggere una sequenza di cifre e mostrarla con soppressione degli zeri iniziali.
Questo esempio prepara l'input dei programmi multicifra successivi.

## Lettura Del Codice
I registri `R0..R3` funzionano come registro a scorrimento: a ogni nuova cifra,
le cifre precedenti si spostano verso la posizione piu' significativa. Quando
arriva `=`, la routine di stampa visita le cifre da sinistra a destra e usa un
flag per evitare zeri iniziali.

## Stato Atteso
Con input `308=`, il display mostra `308`. Con input equivalente a zero, il
programma stampa comunque una singola cifra `0`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/13-calcolatrice-cifre.asm -o calc-cifre.rom
printf '308=' | go run ../go-4004/cmd/retronet-4004 -io calc-cifre.rom
```
