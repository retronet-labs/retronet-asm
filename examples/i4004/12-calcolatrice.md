# 12 - Calcolatrice A Una Cifra

## Obiettivo
Realizzare una calcolatrice a una cifra con quattro operatori: addizione,
sottrazione, moltiplicazione e divisione.

## Lettura Del Codice
La prima parte legge operandi e operatore. Il dispatch confronta il codice
dell'operatore con i simboli convenzionali. Somma e sottrazione usano le routine
BCD gia' viste; moltiplicazione e divisione sono costruite con cicli di
addizioni o sottrazioni ripetute.

## Stato Atteso
Con input `6*7`, il display mostra `42`. L'esempio e' il primo programma della
serie in cui controllo di flusso, I/O e aritmetica cooperano in modo organico.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/12-calcolatrice.asm -o calc.rom
printf '6*7' | go run ../go-4004/cmd/retronet-4004 -io calc.rom
```
