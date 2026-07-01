# 02 - Loop E Flag Zero

## Obiettivo
Calcolare `5 + 4 + 3 + 2 + 1 = 15` usando un contatore e un salto condizionale.

## Lettura Del Codice
`LBI COUNT` inizializza `B` come contatore. `LAI 0` azzera l'accumulatore. Nel
ciclo, `ADB` aggiunge il contatore ad `A`, `DCB` decrementa `B` e aggiorna il
flag Zero. `JFZ loop` salta se Zero e' falso, quindi il ciclo termina quando
`B` diventa zero.

## Stato Atteso
Alla fine `A = 0x0F` e la CPU e' in `HLT`. Il punto teorico e' che il controllo
di flusso non confronta direttamente due valori: osserva un flag prodotto
dall'istruzione precedente.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8008/02-loop.asm -o loop.rom
go run ../retronet-8008/cmd/retronet-8008 -bin loop.rom -steps 100
```
