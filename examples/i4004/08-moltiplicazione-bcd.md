# 08 - Moltiplicazione BCD

## Obiettivo
Calcolare `25 x 5 = 125` in BCD con addizioni ripetute. Si passa da operazioni
su singole cifre a un algoritmo con ciclo esterno e ciclo interno.

## Lettura Del Codice
Il moltiplicando `25` e' memorizzato in RAM come cifre BCD. Il contatore esterno
ripete l'addizione cinque volte; il ciclo interno somma tre cifre, propagando il
Carry con `ADM` e `DAA`. I puntatori RAM vengono incrementati dopo ogni cifra.

## Stato Atteso
Il prodotto e' salvato come `5, 2, 1`, cioe' `125`. Il caso mostra il costo
computazionale di una moltiplicazione su una CPU senza istruzione dedicata.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/08-moltiplicazione-bcd.asm -o mul-bcd.rom
go run ../go-4004/cmd/retronet-4004 -dump-ram mul-bcd.rom
```
