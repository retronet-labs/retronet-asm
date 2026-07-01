# 03 - Decimal Mode

## Obiettivo
Mostrare `ADC` in modalita' decimale BCD. L'esempio calcola `$45 + $55` con
`SED`, ottenendo `$00` con Carry impostato.

## Lettura Del Codice
`SED` abilita il decimal mode del 6502. `CLC` azzera il Carry in ingresso.
`LDA #$45` e `ADC #$55` eseguono la somma come cifre decimali impacchettate,
non come puro binario. Il risultato basso viene salvato in `$0201`.

## Stato Atteso
`$0201` contiene `$00`; il Carry rappresenta la centinaia decimale. L'esempio va
letto come confronto con il BCD del 4004: qui la modalita' decimale e' uno stato
della CPU.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i6502/03-decimal.asm -o decimal6502.rom
go run ../retronet-6502/cmd/retronet-6502 -bin decimal6502.rom -load 0x8000 -trace -steps 30
```
