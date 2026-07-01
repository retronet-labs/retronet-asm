# 01 - Demo Istruzioni Base 8008

## Obiettivo
Mostrare istruzioni 8008 a un byte senza introdurre ancora dati immediati o
indirizzi. L'esempio serve soprattutto per familiarizzare con i mnemonici
compatibili con il disassembler di `retronet-8008`.

## Lettura Del Codice
I mnemonici 8008 codificano i registri nel nome: `LBA` significa `B <- A`,
`LCB` significa `C <- B`. Le operazioni `ADB` e `ADC` usano sempre
l'accumulatore come destinazione. `INB` e `DCC` modificano registri e flag.

## Stato Atteso
Il valore numerico non e' il centro dell'esempio, perche' i registri partono da
zero. L'obiettivo e' verificare forma e codifica delle istruzioni.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8008/01-demo.asm -o demo.rom
go run ../retronet-8008/cmd/retronet-8008 -bin demo.rom -disasm 8
```
