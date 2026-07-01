# 04 - Modi Di Indirizzamento

## Obiettivo
Confrontare indirizzamento immediato, assoluto indicizzato e indiretto
indicizzato. L'esempio usa una tabella e un puntatore in zero page.

## Lettura Del Codice
`#<tabella` e `#>tabella` estraggono byte basso e alto dell'indirizzo della
tabella. Questi byte formano un puntatore in `$00/$01`. `LDA tabella,X` legge
con base assoluta piu' indice; `LDA ($00),Y` legge tramite puntatore zero page
piu' indice `Y`.

## Stato Atteso
`$0200` riceve `$11`, il primo elemento della tabella. `$0201` riceve `$33`, il
terzo elemento. Il documento evidenzia perche' zero page e absolute hanno
dimensioni diverse nell'encoding.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i6502/04-addressing.asm -o addressing6502.rom
go run ../retronet-6502/cmd/retronet-6502 -bin addressing6502.rom -load 0x8000 -trace -steps 80
```
