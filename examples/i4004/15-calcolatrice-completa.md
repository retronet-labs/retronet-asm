# 15 - Calcolatrice Completa A Virgola Fissa

## Obiettivo
Chiudere la serie 4004 con una calcolatrice a quattro operatori e due decimali
fissi. Gli operandi sono rappresentati come interi scalati per 100.

## Lettura Del Codice
Il parser di input interpreta cifre e punto decimale. Addizione e sottrazione
operano direttamente sui valori scalati; moltiplicazione e divisione richiedono
rispettivamente una normalizzazione per 100 e una pre-moltiplicazione per 100.
Le routine di copia tra registri RAM organizzano i valori temporanei.

## Stato Atteso
Input `1.5+2.25=` produce `3.75`; input `7/2=` produce `3.50`. Il programma
mostra come vincoli severi di registri e nibble portino a una disciplina di
rappresentazione dei dati.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i4004/15-calcolatrice-completa.asm -o calc-completa.rom
printf '1.5+2.25=' | go run ../go-4004/cmd/retronet-4004 -io calc-completa.rom
```
