.arch i8086

; Boot sector dimostrativo per retronet-pc: legge i tasti col servizio tastiera
; del BIOS (INT 16h, AH=00h, attesa bloccante) e li ristampa col teletype
; (INT 10h, AH=0Eh), in loop. Prova l'input da tastiera dell'emulatore.

.orgbase 0x7C00

read:   mov ah, 0x00      ; attendi un tasto
        int 0x16          ; AL = ASCII
        mov ah, 0x0E      ; funzione teletype
        int 0x10          ; stampa AL
        jmp read

        .org 0x7DFE       ; riempi fino all'offset 510
        .byte 0x55, 0xAA  ; firma di boot
