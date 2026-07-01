# 08 - Memoria Video Diretta

## Obiettivo
Mostrare accesso diretto a memoria video senza chiamare il BIOS. L'esempio usa
`ES:DI` e `STOSW` per scrivere celle testuali MDA.

## Lettura Del Codice
`ES` viene impostato a `0xB000`, segmento tipico della memoria testuale MDA.
`DI` parte da zero. Ogni `STOSW` scrive `AX` in `ES:DI` e avanza `DI` di due
byte. Il byte basso della word e' il carattere, il byte alto l'attributo.

## Stato Atteso
La prima riga video contiene `OK` con attributo `07h`. Il programma introduce
l'idea che in real mode molte periferiche siano mappate in memoria.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/08-video-diretto.asm -o video86.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy video86.rom -steps 12000000
```
