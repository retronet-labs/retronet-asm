.arch i8080

; Subroutine e accesso sequenziale alla memoria.
; HL punta a una tabella di tre byte; la subroutine sum3 li somma in A.

        LXI H, numeri   ; HL = indirizzo del primo elemento
        CALL sum3
        STA totale      ; salva 2 + 4 + 6 = 12
        HLT

sum3:   MVI B, 3        ; tre elementi da leggere
        MVI A, 0
loop:   ADD M           ; A += memoria[HL]
        INX H           ; HL punta all'elemento successivo
        DCR B
        JNZ loop
        RET

numeri: .byte 2, 4, 6
totale: .byte 0
