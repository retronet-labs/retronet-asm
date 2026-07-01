# 01 - Hello Su Porta 8080

## Obiettivo
Mostrare l'I/O separato dell'Intel 8080 con un programma minimo che scrive
`HI` sulla porta convenzionale `1`.

## Lettura Del Codice
`MVI A, 0x48` carica il codice ASCII di `H` nell'accumulatore. `OUT 1` invia
quel byte alla porta. La sequenza viene ripetuta per `I`, poi `HLT` arresta la
CPU. Non c'e' memoria dati: tutto passa dall'accumulatore.

## Stato Atteso
Un emulatore che collega la porta `1` a un terminale visualizza `HI`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/01-hello-port.asm -o hello8080.rom
go run ../retronet-8080/cmd/retronet-8080 -bin hello8080.rom -terminal -steps 20
```
