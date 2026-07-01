.arch i8086
.orgbase 0x7C00

; Stampa le cifre 0..9 usando CX e LOOP.
; BL contiene il carattere ASCII corrente.

        mov cx, 10
        mov bl, 0x30      ; '0'

print_digit:
        mov ah, 0x0E
        mov al, bl
        int 0x10
        inc bl
        loop print_digit

halt:   jmp halt

        .org 0x7DFE
        .byte 0x55, 0xAA
