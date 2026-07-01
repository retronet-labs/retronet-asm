# 09 - Somma A 16 Bit

## Obiettivo
Mostrare la famiglia di istruzioni su coppie registro. `DAD` somma una coppia a
`HL`, rendendo possibile aritmetica a 16 bit su una CPU a 8 bit.

## Lettura Del Codice
`LXI H, 0x1234` e `LXI D, 0x0102` caricano due parole. `DAD D` calcola
`HL = HL + DE`. `SHLD result` salva `L` e poi `H`, cioe' il risultato in formato
little-endian.

## Stato Atteso
`result` contiene i byte `0x36, 0x13`, corrispondenti alla parola `0x1336`.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/09-somma-16bit.asm -o somma16.rom
go run ../retronet-8080/cmd/retronet-8080 -bin somma16.rom -trace -steps 20
```
