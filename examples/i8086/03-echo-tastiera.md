# 03 - Echo Tastiera

## Obiettivo
Leggere tasti dal BIOS e ristamparli a video. L'esempio introduce `INT 16h` per
la tastiera e riusa `INT 10h` per l'output.

## Lettura Del Codice
`mov ah, 0x00` seleziona la funzione BIOS di attesa tasto. Dopo `INT 16h`, il
codice ASCII e' in `AL`. Impostando `AH=0x0E`, `INT 10h` stampa quello stesso
carattere. Il ciclo `jmp read` rende il programma interattivo.

## Stato Atteso
Ogni tasto premuto viene visualizzato. Non c'e' buffer applicativo: il BIOS
fornisce direttamente il carattere gia' letto.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/03-echo-tastiera.asm -o echo.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy echo.rom -keys "ABC" -steps 14000000
```
