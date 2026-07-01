# 10 - Input Bufferizzato

## Obiettivo
Leggere caratteri da tastiera, salvarli in memoria e ristamparli. L'esempio
riunisce BIOS keyboard, BIOS video, `STOSB`, buffer e terminatore zero.

## Lettura Del Codice
`INT 16h` legge un carattere in `AL`. Se `AL` e' `0x0D` (`Enter`), la riga e'
terminata. Altrimenti `STOSB` scrive il carattere in `ES:DI` e avanza `DI`;
subito dopo il carattere viene visualizzato con `INT 10h`. Al termine viene
aggiunto uno zero e il buffer viene stampato come stringa.

## Stato Atteso
Il boot sector fa eco ai caratteri durante la digitazione e, dopo `Enter`,
ristampa la riga salvata. Per restare didattico non include controllo di
overflow: il buffer e' volutamente piccolo.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/10-input-buffer.asm -o input86.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy input86.rom -keys $'CIAO\r' -steps 16000000
```
