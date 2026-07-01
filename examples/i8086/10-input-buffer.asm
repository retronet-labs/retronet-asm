.arch i8086
.orgbase 0x7C00

; Legge una riga fino a Enter, la salva in un buffer e poi la ristampa.
; Esempio didattico: non applica un limite alla lunghezza della riga.

        xor ax, ax
        mov ds, ax
        mov es, ax

        mov si, prompt
        call print
        mov di, buffer

read:   mov ah, 0x00
        int 0x16
        cmp al, 0x0D
        je  done
        stosb
        mov ah, 0x0E
        int 0x10
        jmp read

done:   mov al, 0
        stosb
        mov si, newline
        call print
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

prompt: .byte "Scrivi e premi Enter: ", 0
newline:.byte 0x0D, 0x0A, 0
buffer: .byte 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

        .org 0x7DFE
        .byte 0x55, 0xAA
