# 03 - Subroutine E Stack Hardware

## Obiettivo
Introdurre `CAL` e `RET` sull'8008. La subroutine `dbl` raddoppia il contenuto
dell'accumulatore.

## Lettura Del Codice
Il programma carica `9` in `A`, chiama `dbl` e poi si ferma. `CAL dbl` salva
l'indirizzo di ritorno nello stack hardware interno della CPU. Nella subroutine,
`ADA` somma `A` con se stesso; `RET` ripristina il PC salvato.

## Stato Atteso
Alla fine `A = 18` (`0x12`). L'esempio evidenzia che lo stack dell'8008 non e'
una zona di memoria generale: e' una struttura interna limitata per il controllo
di flusso.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8008/03-subroutine.asm -o sub.rom
go run ../retronet-8008/cmd/retronet-8008 -bin sub.rom -steps 100
```
