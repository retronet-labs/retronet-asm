# 02 - Aritmetica Su Registri 8080

## Obiettivo
Calcolare `5 + 7` usando registri generali e accumulatore. L'esempio introduce
la forma classica dell'ALU 8080: sorgente esplicita, destinazione implicita `A`.

## Lettura Del Codice
`MVI` carica costanti in `B` e `C`. `MOV A, B` prepara l'accumulatore, poi
`ADD C` aggiorna `A` e i flag. `STA 0x2000` copia il risultato in memoria con
indirizzamento assoluto.

## Stato Atteso
`A` e la cella `0x2000` contengono `12`. I flag riflettono la somma, ma non
vengono ancora usati per saltare.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/02-aritmetica-registri.asm -o arit8080.rom
go run ../retronet-8080/cmd/retronet-8080 -bin arit8080.rom -trace -steps 20
```
