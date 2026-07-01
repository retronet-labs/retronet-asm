.arch i8086
.orgbase 0x7C00

; Boot sector con mini-menu.
; Stampa una scelta, legge un tasto con INT 16h e usa una subroutine per
; stampare stringhe terminate da zero tramite INT 10h/AH=0Eh.

        xor ax, ax
        mov ds, ax
        mov si, menu
        call print

read:   mov ah, 0x00
        int 0x16
        cmp al, 0x31      ; '1'
        je  one
        cmp al, 0x32      ; '2'
        je  two
        mov si, invalid
        call print
        jmp read

one:    mov si, msg1
        call print
        jmp halt

two:    mov si, msg2
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

menu:   .byte "1) SOMMA  2) USCITA", 0x0D, 0x0A, 0
invalid:.byte "Scelta non valida", 0x0D, 0x0A, 0
msg1:   .byte "Hai scelto SOMMA", 0x0D, 0x0A, 0
msg2:   .byte "Fine", 0x0D, 0x0A, 0

        .org 0x7DFE
        .byte 0x55, 0xAA
