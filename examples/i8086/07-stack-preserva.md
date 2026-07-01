# 07 - Stack E Preservazione Registri

## Obiettivo
Usare `PUSH` e `POP` in una subroutine per preservare registri temporanei. La
disciplina di salvataggio e ripristino e' essenziale nei programmi real mode.

## Lettura Del Codice
Il chiamante passa in `SI` l'indirizzo della stringa. `print_preserve` salva
`AX` e `SI`, usa `LODSB` e `INT 10h` per stampare, poi ripristina i registri in
ordine inverso. `RET` rientra al chiamante usando l'indirizzo salvato da `CALL`.

## Stato Atteso
Il boot sector stampa due righe. Il valore iniziale di `SI` del chiamante viene
preservato dalla subroutine, anche se internamente `LODSB` lo incrementa.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8086/07-stack-preserva.asm -o stack86.rom
go run ../retronet-pc/cmd/retronet-pc -bios ../retronet-pc/GLABIOS_0.4.2_8X.ROM -floppy stack86.rom -steps 12000000
```
