# 07 - Sottrazione BCD Multicifra

## Obiettivo
Eseguire `52 - 27 = 25` propagando il prestito tra cifre BCD. L'esempio e'
parallelo alla somma multicifra, ma con semantica di Carry opposta.

## Lettura Del Codice
Minuendo, sottraendo e risultato occupano registri RAM separati. Il ciclo visita
unita' e decine. `TCS` traduce il Carry precedente nella base corretta,
`SBM` sottrae la cifra in memoria, `ADM` aggiunge la cifra del minuendo e `DAA`
normalizza il risultato.

## Stato Atteso
La RAM del risultato contiene `5, 2`. Il Carry finale resta `1`, indicando che la
sottrazione complessiva non ha richiesto un prestito oltre la cifra piu'
significativa.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/07-sottrazione-multicifra.asm -o sub-multi.rom
go run ../go-4004/cmd/retronet-4004 -dump-ram sub-multi.rom
```
