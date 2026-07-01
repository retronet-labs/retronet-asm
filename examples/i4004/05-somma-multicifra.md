# 05 - Somma BCD Multicifra

## Obiettivo
Generalizzare la somma BCD a piu' cifre: `47 + 58 = 105`. L'esempio mostra la
propagazione del riporto tra cifre decimali.

## Lettura Del Codice
Le cifre sono in RAM in ordine little-endian: prima unita', poi decine. Tre
puntatori paralleli scorrono addendo A, addendo B e risultato. A ogni iterazione
`ADM` include il Carry precedente, `DAA` normalizza la cifra BCD e `WRM` salva il
risultato parziale. `TCC` trasforma l'ultimo Carry nella cifra delle centinaia.

## Stato Atteso
Il risultato in RAM e' `5, 0, 1`, quindi il numero decimale `105`. Il punto
centrale e' che il flag Carry diventa parte dello stato dell'algoritmo.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/05-somma-multicifra.asm -o somma-multi.rom
go run ../go-4004/cmd/retronet-4004 -dump-ram somma-multi.rom
```
