# 03 - Loop Con Flag

## Obiettivo
Usare il flag Zero per controllare un ciclo. Il programma somma i numeri da 5 a
1 e salva `15`.

## Lettura Del Codice
`B` e' il contatore. Ogni iterazione aggiunge `B` ad `A`, poi `DCR B` decrementa
il contatore e aggiorna il flag Zero. `JNZ loop` torna all'inizio finche' il
flag Zero non e' impostato.

## Stato Atteso
La cella `0x2001` contiene `0x0F`. L'esempio chiarisce il contratto tra una
istruzione che produce flag e il salto condizionale che li consuma.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/03-loop-flag.asm -o loop8080.rom
go run ../retronet-8080/cmd/retronet-8080 -bin loop8080.rom -trace -steps 50
```
