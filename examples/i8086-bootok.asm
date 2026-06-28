.arch i8086

; Boot sector dimostrativo per retronet-pc: stampa un messaggio usando il
; servizio teletype del BIOS (INT 10h, AH=0Eh) e poi si ferma in loop.
; Il BIOS carica il settore a 0000:7C00 e vi salta: .orgbase 0x7C00 fa risolvere
; le label all'indirizzo di caricamento, senza padding iniziale.

.orgbase 0x7C00

        xor ax, ax        ; DS = 0 per leggere il messaggio
        mov ds, ax
        mov si, msg
        mov ah, 0x0E      ; funzione teletype
print:  lodsb             ; AL = [DS:SI], SI++
        cmp al, 0         ; fine stringa?
        je  halt
        int 0x10          ; stampa AL
        jmp print
halt:   jmp halt

msg:    .byte "RETRONET-PC: BOOT OK!", 0

        .org 0x7DFE       ; riempi fino all'offset 510
        .byte 0x55, 0xAA  ; firma di boot
