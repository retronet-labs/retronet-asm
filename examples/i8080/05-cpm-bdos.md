# 05 - Programma CP/M Con BDOS

## Obiettivo
Generare un programma `.COM` CP/M-like. L'esempio usa il punto di ingresso BDOS
standard `CALL 0005h` per stampare una stringa terminata da `$`.

## Lettura Del Codice
`.com` imposta l'origine logica a `0x0100`, come richiesto dai programmi CP/M.
`DE` punta alla stringa, `C` contiene il numero della funzione BDOS. La funzione
`9` stampa fino a `$`; la funzione `0` termina il programma.

## Stato Atteso
Eseguito in `retronet-cpm`, il programma stampa `CIAO DA CP/M` e ritorna al
sistema. L'esempio separa chiaramente il formato binario `.COM` dalla ROM nuda.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/05-cpm-bdos.asm -o CIAO.COM
go run ../retronet-cpm/cmd/retronet-cpm -run CIAO.COM
```
