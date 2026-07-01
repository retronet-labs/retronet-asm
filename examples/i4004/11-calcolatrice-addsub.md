# 11 - Calcolatrice: Addizione E Sottrazione

## Obiettivo
Estendere la calcolatrice a due operatori. Il programma legge `cifra operatore
cifra` e sceglie tra somma e sottrazione.

## Lettura Del Codice
L'operatore viene confrontato con il codice della somma tramite `STC`, `SUB` e
`JCN`. Il ramo `do_add` usa la sequenza BCD della somma; il ramo `do_sub` usa la
sequenza con `TCS`. Il blocco `show` riunifica l'output.

## Stato Atteso
Input come `7+5` produce `12`, mentre `9-4` produce `05`. Il programma mostra il
pattern classico dispatch-routine-output.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/11-calcolatrice-addsub.asm -o calc-addsub.rom
printf '7+5' | go run ../go-4004/cmd/retronet-4004 -io calc-addsub.rom
```
