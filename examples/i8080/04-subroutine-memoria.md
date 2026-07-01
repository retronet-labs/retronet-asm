# 04 - Subroutine E Memoria

## Obiettivo
Somma di una piccola tabella in memoria tramite subroutine. Si introducono
`CALL`, `RET`, la coppia `HL` e l'operando memoria `M`.

## Lettura Del Codice
`LXI H, numeri` mette in `HL` l'indirizzo della tabella. La subroutine `sum3`
usa `M` come alias di `memoria[HL]`; dopo ogni lettura, `INX H` passa
all'elemento successivo. `B` conta gli elementi rimasti. `RET` restituisce il
controllo al chiamante.

## Stato Atteso
La label `totale` riceve `12`. L'esempio mostra come l'8080 rappresenti molti
accessi memoria attraverso la coppia `HL`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/04-subroutine-memoria.asm -o mem8080.rom
go run ../retronet-8080/cmd/retronet-8080 -bin mem8080.rom -trace -steps 80
```
