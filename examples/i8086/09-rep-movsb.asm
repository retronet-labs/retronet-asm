.arch i8086
.orgbase 0x7C00

; Copia una stringa con REP MOVSB, poi stampa il buffer copiato.

        xor ax, ax
        mov ds, ax
        mov es, ax

        mov si, msg
        mov di, buffer
        mov cx, 6
        rep
        movsb
        mov byte [di], 0

        mov si, buffer
        call print

halt:   jmp halt

print:  mov ah, 0x0E
print_next:
        lodsb
        cmp al, 0
        je  print_done
        int 0x10
        jmp print_next
print_done:
        ret

msg:    .byte "COPIA!"
buffer: .byte 0, 0, 0, 0, 0, 0, 0

        .org 0x7DFE
        .byte 0x55, 0xAA
