# 06 - Contatore Con LOOP

## Obiettivo
Mostrare l'istruzione `LOOP`, che decrementa `CX` e salta finche' `CX` non
diventa zero. Il boot sector stampa le cifre ASCII da `0` a `9`.

## Lettura Del Codice
`CX` contiene il numero di iterazioni. `BL` contiene il codice ASCII corrente.
A ogni giro il BIOS teletype (`INT 10h`, `AH=0Eh`) stampa `AL`; poi `INC BL`
passa alla cifra successiva. `LOOP print_digit` aggiorna `CX` implicitamente.

## Stato Atteso
Sul video compare `0123456789`. L'esempio distingue `LOOP`, basato su `CX`, dai
salti condizionali basati sui flag.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/06-contatore-loop.asm -o loop86.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy loop86.rom -steps 12000000
```
