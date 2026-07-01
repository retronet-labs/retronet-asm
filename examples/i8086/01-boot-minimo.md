# 01 - Boot Sector Minimo 8086

## Obiettivo
Creare il boot sector 8086 piu' piccolo utile: codice eseguibile, loop stabile
e firma BIOS `55 AA` agli ultimi due byte.

## Lettura Del Codice
`.orgbase 0x7C00` fa risolvere le label come se il BIOS avesse caricato il
settore a `0000:7C00`, senza aggiungere padding all'inizio del file. `CLI`
disabilita interrupt, `HLT` ferma la CPU fino al prossimo evento e `JMP halt`
mantiene il controllo nel settore.

## Stato Atteso
Il file prodotto e' lungo 512 byte e termina con `0x55, 0xAA`. Non produce
output visibile: serve come base strutturale per i boot sector successivi.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/01-boot-minimo.asm -o boot-minimo.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy boot-minimo.rom -steps 12000000
```
