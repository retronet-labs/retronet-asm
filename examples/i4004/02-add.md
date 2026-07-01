# 02 - Addizione Binaria

## Obiettivo
Mostrare la prima operazione ALU sul 4004: sommare due nibble in binario senza
accesso alla RAM e senza correzione decimale.

## Lettura Del Codice
`LDM 3` carica il primo addendo, `XCH R1` lo deposita in un registro generale e
`LDM 4` prepara il secondo addendo. `ADD R1` calcola `4 + 3` nell'accumulatore.
Il risultato resta in `A`, perche' l'esercizio isola il comportamento dell'ALU.

## Stato Atteso
Alla label `halt`, l'accumulatore contiene `7`. Il flag Carry resta non
essenziale perche' la somma non supera il range del nibble.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/02-add.asm -o add.rom
go run ../go-4004/cmd/retronet-4004 -trace add.rom
```
