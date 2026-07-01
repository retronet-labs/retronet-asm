# 01 - Reset Vector 6502

## Obiettivo
Costruire una ROM 6502 minima con vettore di reset. Il programma inizializza lo
stack e scrive un byte osservabile in RAM.

## Lettura Del Codice
`.orgbase $8000` dice all'assembler che il codice sara' visto dalla CPU a
partire da `$8000`. `SEI` e `CLD` fissano uno stato iniziale semplice. `LDX #$FF`
e `TXS` impostano lo stack pointer. Alla fine `.org $FFFC` posiziona il vettore
RESET, emesso da `.word reset` in little-endian.

## Stato Atteso
Dopo il reset, la cella `$0200` contiene `$01` e la CPU resta nel loop `halt`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i6502/01-reset-vector.asm -o reset6502.rom
go run ../retronet-6502/cmd/retronet-6502 -bin reset6502.rom -load 0x8000 -trace -steps 20
```
