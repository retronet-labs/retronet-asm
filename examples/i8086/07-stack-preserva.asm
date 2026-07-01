.arch i8086
.orgbase 0x7C00

; Subroutine che preserva i registri usati internamente con PUSH/POP.

        xor ax, ax
        mov ds, ax

        mov si, msg1
        call print_preserve
        mov si, msg2
        call print_preserve

halt:   jmp halt

print_preserve:
        push ax
        push si
        mov ah, 0x0E
print_next:
        lodsb
        cmp al, 0
        je  print_done
        int 0x10
        jmp print_next
print_done:
        pop si
        pop ax
        ret

msg1:   .byte "Prima riga", 0x0D, 0x0A, 0
msg2:   .byte "Seconda riga", 0x0D, 0x0A, 0

        .org 0x7DFE
        .byte 0x55, 0xAA
