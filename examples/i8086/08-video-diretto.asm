.arch i8086
.orgbase 0x7C00

; Scrittura diretta in memoria video MDA: segmento B000h, offset 0000.
; Ogni word e' carattere ASCII nel byte basso e attributo nel byte alto.

        mov ax, 0xB000
        mov es, ax
        mov di, 0

        mov ax, 0x074F    ; 'O' con attributo 07h
        stosw
        mov ax, 0x074B    ; 'K'
        stosw

halt:   jmp halt

        .org 0x7DFE
        .byte 0x55, 0xAA
