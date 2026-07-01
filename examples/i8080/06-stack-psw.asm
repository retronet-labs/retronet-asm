.arch i8080

; Stack 8080: PUSH/POP di coppie registro e PSW.
; Imposta SP, salva A/B/C, sporca i registri e li ripristina.

        LXI SP, 0x2400
        MVI A, 0x12
        MVI B, 0x34
        MVI C, 0x56

        PUSH B          ; salva BC
        PUSH PSW        ; salva A e flags

        MVI A, 0
        MVI B, 0
        MVI C, 0

        POP PSW         ; ripristina A e flags
        POP B           ; ripristina BC

        STA 0x2100
        MOV A, B
        STA 0x2101
        MOV A, C
        STA 0x2102
        HLT
