# 08 - Stringa Da Memoria A Porta

## Obiettivo
Scorrere una stringa in memoria e inviarla a una porta di output. L'esempio
introduce il pattern `HL` come puntatore e `M` come byte puntato.

## Lettura Del Codice
`LXI H, msg` inizializza il puntatore. A ogni iterazione `MOV A, M` legge il
carattere corrente; `CPI 0` riconosce il terminatore. Se il byte non e' zero,
`OUT 1` lo invia alla porta e `INX H` passa al carattere successivo.

## Stato Atteso
Un terminale collegato alla porta `1` riceve `8080`. Il programma mostra come
l'8080 emuli un accesso indicizzato usando una coppia registro.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/08-stringa-porta.asm -o stringa8080.rom
go run ../retronet-8080/cmd/retronet-8080 -bin stringa8080.rom -terminal -steps 80
```
