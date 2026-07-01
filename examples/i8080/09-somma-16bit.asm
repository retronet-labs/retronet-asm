.arch i8080

; Aritmetica a 16 bit con coppie registro.
; Calcola 0x1234 + 0x0102 = 0x1336 e salva il risultato little-endian.

        LXI H, 0x1234
        LXI D, 0x0102
        DAD D           ; HL = HL + DE
        SHLD result
        HLT

result: .byte 0, 0
