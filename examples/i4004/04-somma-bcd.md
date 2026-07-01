# 04 - Somma BCD A Cifra Singola

## Obiettivo
Introdurre l'aritmetica decimale codificata in BCD. La somma `7 + 5` produce una
cifra unita' e un riporto decimale.

## Lettura Del Codice
Gli operandi vengono scritti in RAM con `WRM`, poi riletti con `RDM` e sommati
con `ADM`. `CLC` azzera il riporto in ingresso. Dopo la somma binaria, `DAA`
corregge il risultato in BCD: l'accumulatore diventa `2` e il Carry rappresenta
la decina.

## Stato Atteso
La cella risultato contiene la cifra delle unita' `2`; la porta riceve lo stesso
valore. Il Carry vale `1`, cioe' la somma completa e' `12`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/04-somma-bcd.asm -o somma-bcd.rom
go run ../go-4004/cmd/retronet-4004 -dump-ram somma-bcd.rom
```
