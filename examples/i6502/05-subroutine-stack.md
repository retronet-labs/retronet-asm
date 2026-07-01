# 05 - Subroutine E Stack

## Obiettivo
Usare `JSR`, `RTS`, `PHA` e `PLA` per calcolare `A * 3`. L'esempio mostra lo
stack come supporto sia al controllo di flusso sia ai dati temporanei.

## Lettura Del Codice
Il chiamante carica `5` in `A` e invoca `triple`. La subroutine salva il valore
originale con `PHA`, raddoppia `A` con `ASL A`, conserva il doppio in zero page,
recupera l'originale con `PLA` e somma con `ADC $20`.

## Stato Atteso
La cella `$0202` contiene `$0F`. Il valore dello stack pointer viene
inizializzato esplicitamente per rendere il comportamento riproducibile.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i6502/05-subroutine-stack.asm -o stack6502.rom
go run ../retronet-6502/cmd/retronet-6502 -bin stack6502.rom -load 0x8000 -trace -steps 80
```
