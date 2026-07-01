# 02 - Loop E Branch Relativo

## Obiettivo
Calcolare `5 x 3` come somma ripetuta su 6502. L'esempio introduce `DEX` e
`BNE`, cioe' il ciclo basato sul flag Zero.

## Lettura Del Codice
`X` contiene il numero di iterazioni, `A` l'accumulatore. Ogni giro esegue
`ADC #$03`; `CLC` prima della somma evita che un Carry precedente alteri il
risultato. `DEX` decrementa `X` e imposta Zero quando il contatore arriva a
zero. `BNE loop` usa un offset relativo calcolato dall'assembler.

## Stato Atteso
La cella `$0200` contiene `$0F`. Il vettore a `$FFFC` rende il file una ROM
avviabile nel modello 6502.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i6502/02-loop.asm -o loop6502.rom
go run ../retronet-6502/cmd/retronet-6502 -bin loop6502.rom -load 0x8000 -trace -steps 80
```
