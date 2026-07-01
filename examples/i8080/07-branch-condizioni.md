# 07 - Branch Condizionali

## Obiettivo
Usare i flag prodotti da `CPI` per scegliere un ramo. Il programma classifica
un valore rispetto a una soglia.

## Lettura Del Codice
`CPI 5` sottrae logicamente `5` da `A` senza modificare `A`, aggiornando pero'
i flag. Se il valore e' minore della soglia, il Carry viene impostato e `JC`
salta a `minore`. In caso contrario viene salvato `1`.

## Stato Atteso
Con `A = 7`, la cella `0x2103` contiene `1`. Cambiando il valore iniziale sotto
`5`, la stessa cella conterrebbe `0`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/07-branch-condizioni.asm -o branch8080.rom
go run ../retronet-8080/cmd/retronet-8080 -bin branch8080.rom -trace -steps 20
```
