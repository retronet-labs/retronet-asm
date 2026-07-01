# Esempi Didattici

Gli esempi sono organizzati per architettura e ordinati per difficolta'
crescente. Ogni programma ha due file:

- `nome.asm`: sorgente assemblabile con commenti locali.
- `nome.md`: spiegazione didattica, stato atteso, comando di build e comando
  per provarlo sull'emulatore corrispondente.

## Flusso Di Lavoro

```sh
go run ./cmd/retronet-asm build examples/i4004/01-hello4004.asm -o out.rom
```

Senza `-o`, il file generato prende il nome del sorgente con estensione `.rom`.
Gli output generati non sono versionati.

## Percorsi

| Cartella | Architettura | Focus |
|----------|--------------|-------|
| `i4004/` | Intel 4004 | nibble, RAM, BCD, calcolatrici |
| `i8008/` | Intel 8008 | registri nel mnemonico, loop, terminale |
| `i8080/` | Intel 8080 | accumulator machine, stack, memoria, CP/M `.COM` |
| `i6502/` | MOS/NMOS 6502 | vettori, zero page, branch, stack |
| `i8086/` | Intel 8086/8088 | boot sector PC, BIOS, stack, video, stringhe |

## Ordine Consigliato

1. `i4004`: partire da `01-hello4004` e arrivare alle calcolatrici. Questa serie
   mostra come costruire aritmetica e controllo partendo da una CPU a 4 bit.
2. `i8008`: leggere i mnemonici come codifica dei registri e confrontare
   l'aritmetica binaria con il BCD del 4004.
3. `i8080`: osservare la forma piu' regolare dell'ISA e il passaggio a CP/M.
4. `i6502`: studiare address mode, zero page e vettori hardware.
5. `i8086`: concludere con boot sector real mode e servizi BIOS.

## Note Di Esecuzione

- Gli esempi `i4004` producono ROM per l'emulatore 4004 e usano spesso RAM o
  porte virtuali.
- Gli esempi `i8008` e molti esempi `i8080` sono ROM nude. `i8080/05-cpm-bdos.asm`
  e `i8080/10-cpm-echo.asm` generano programmi `.COM`.
- Gli esempi `i6502` includono il vettore RESET a `$FFFC`; per questo le ROM
  sono piu' grandi dei programmi effettivi.
- Gli esempi `i8086` sono boot sector da 512 byte con firma `55 AA`; il backend
  puo' comunque assemblare anche codice real mode non incapsulato come boot sector.
