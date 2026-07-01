# 06 - Stack E PSW

## Obiettivo
Mostrare l'uso dello stack 8080 per salvare registri e Program Status Word.
L'esempio e' una base per comprendere chiamate annidate e preservazione dello
stato.

## Lettura Del Codice
`LXI SP, 0x2400` inizializza lo stack pointer. `PUSH B` salva la coppia `BC`;
`PUSH PSW` salva accumulatore e flag. Dopo aver azzerato i registri, `POP PSW`
e `POP B` ripristinano i valori nello stesso ordine inverso.

## Stato Atteso
Le celle `0x2100`, `0x2101` e `0x2102` contengono rispettivamente `0x12`,
`0x34` e `0x56`. L'esempio evidenzia la disciplina LIFO dello stack.

## Comandi
```sh
go run ./cmd/retronet-asm build examples/i8080/06-stack-psw.asm -o stack8080.rom
go run ../retronet-8080/cmd/retronet-8080 -bin stack8080.rom -trace -steps 40
```
