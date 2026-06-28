.arch i8086

; Dimostrazione degli operandi in memoria del backend i8086: stampa un messaggio
; leggendolo con l'indirizzamento indicizzato [msg+bx] (base BX + spiazzamento
; simbolico), invece di LODSB. Boot sector per retronet-pc.

.orgbase 0x7C00

        xor ax, ax
        mov ds, ax        ; DS = 0
        mov bx, 0         ; indice
        mov ah, 0x0E      ; teletype
print:  mov al, [msg+bx]  ; AL = [DS:msg+BX]  (operando in memoria indicizzato)
        cmp al, 0
        je  halt
        int 0x10
        inc bx
        jmp print
halt:   jmp halt

msg:    .byte "Indexed memory OK!", 0

        .org 0x7DFE
        .byte 0x55, 0xAA
