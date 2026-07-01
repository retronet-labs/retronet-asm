# 10 - CP/M Echo

## Obiettivo
Completare la serie 8080 con un programma `.COM` interattivo: prompt, lettura
di un carattere tramite BDOS e ristampa.

## Lettura Del Codice
`.com` imposta l'origine logica a `0x0100`. La funzione BDOS `9` stampa la
stringa terminata da `$`; la funzione `1` legge un carattere da console e lo
ritorna in `A`; la funzione `2` stampa il carattere contenuto in `E`. La funzione
`0` termina il programma.

## Stato Atteso
In `retronet-cpm`, il programma mostra `Premi un tasto: `, legge un carattere e
lo ristampa. BDOS puo' gia' ecoare l'input, ma qui l'output esplicito rende
visibile il contratto della funzione `2`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/10-cpm-echo.asm -o ECHO.COM
go run ../retronet-cpm/cmd/retronet-cpm -run ECHO.COM -input "Z"
```
