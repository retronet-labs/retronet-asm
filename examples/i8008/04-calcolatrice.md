# 04 - Calcolatrice Binaria A Una Cifra

## Obiettivo
Leggere dal terminale una espressione del tipo `6*7=` e stampare il risultato
decimale. L'8008 lavora in binario: non esiste `DAA`.

## Lettura Del Codice
`INP 0` legge caratteri ASCII dalla porta terminale; sottraendo `ZERO` si ottiene
la cifra numerica. Il dispatch confronta l'operatore con `CPI`. Addizione e
sottrazione sono dirette, mentre moltiplicazione e divisione usano cicli
elementari. L'output converte il valore binario in cifre ASCII.

## Stato Atteso
Con input terminale `6*7=`, il programma stampa `42`. Il limite didattico e'
esplicito: operandi a una cifra e risultato nell'intervallo 0..255.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8008/04-calcolatrice.asm -o calc.rom
go run ../retronet-8008/cmd/retronet-8008 -bin calc.rom -terminal-input "6*7=" -steps 100000
```
