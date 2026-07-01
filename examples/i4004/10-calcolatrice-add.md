# 10 - Calcolatrice: Addizione

## Obiettivo
Leggere due cifre dalla periferica virtuale e visualizzare la somma BCD. Questo
e' il primo esempio interattivo della serie 4004.

## Lettura Del Codice
`RDR` legge i due tasti e li deposita in `R0` e `R1`. La somma usa `ADD` e
`DAA`; la cifra delle unita' viene salvata in `R2`, mentre `TCC` converte il
Carry nella cifra delle decine. `WMP` invia in sequenza decine e unita'.

## Stato Atteso
Con input `75`, il display riceve `12`. L'esempio separa chiaramente ingresso,
calcolo e presentazione del risultato.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/10-calcolatrice-add.asm -o calc.rom
printf '75' | go run ../go-4004/cmd/retronet-4004 -io calc.rom
```
