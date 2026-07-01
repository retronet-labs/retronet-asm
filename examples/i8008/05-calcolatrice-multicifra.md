# 05 - Calcolatrice Binaria Multicifra

## Obiettivo
Estendere la calcolatrice 8008 a numeri decimali di piu' cifre, mantenendo
calcolo interno binario a 8 bit.

## Lettura Del Codice
La routine `readnum` accumula un numero leggendo cifre ASCII e calcolando
`valore = valore * 10 + cifra`. Gli operatori sono gestiti con lo stesso schema
di dispatch dell'esempio precedente. La routine di stampa converte il byte in
centinaia, decine e unita'.

## Stato Atteso
Con input `12*12=`, l'output e' `144`. L'esempio mostra la distinzione fra
rappresentazione esterna decimale e rappresentazione interna binaria.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8008/05-calcolatrice-multicifra.asm -o calc-multi.rom
go run ../retronet-8008/cmd/retronet-8008 -bin calc-multi.rom -terminal-input "12*12=" -steps 2000000
```
