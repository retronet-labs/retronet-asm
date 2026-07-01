# 03 - Moltiplicazione Per Addizioni Ripetute

## Obiettivo
Costruire `3 x 4` usando solo addizione e controllo di ciclo. L'esempio mette in
relazione i limiti hardware del 4004 con la costruzione algoritmica.

## Lettura Del Codice
`FIM R0, 0x03` mette l'addendo in `R1`. La coppia `R2/R3` punta alla cella RAM
di destinazione tramite `SRC R2`. Il registro `R4` contiene `12`, cioe' `16 - 4`:
`ISZ` incrementa il registro e ripete finche' non torna a zero, ottenendo quattro
iterazioni.

## Stato Atteso
La cella RAM selezionata riceve `12`, codificato come nibble esadecimale `0xC`.
Il programma chiarisce che il ciclo `ISZ` del 4004 lavora per conteggio modulo
16, non con un confronto esplicito.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/03-moltiplicazione.asm -o mul.rom
go run ../go-4004/cmd/retronet-4004 -dump-ram mul.rom
```
