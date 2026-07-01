# 09 - Divisione BCD

## Obiettivo
Calcolare `7 / 2 = 3` con resto `1` usando sottrazioni ripetute. L'esempio
introduce un ramo condizionale basato sul Carry.

## Lettura Del Codice
Il registro del resto parte dal dividendo. A ogni giro viene tentata una
sottrazione del divisore. Se il Carry indica assenza di prestito, `JCN` porta a
`commit`, dove il resto viene aggiornato e il quoziente incrementato. Se invece
la sottrazione richiede prestito, il ciclo termina.

## Stato Atteso
La RAM contiene quoziente `3` e resto `1`. Didatticamente e' importante notare
che la sottrazione non valida viene scartata.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/09-divisione-bcd.asm -o div-bcd.rom
go run ../go-4004/cmd/retronet-4004 -dump-ram div-bcd.rom
```
