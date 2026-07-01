; moltiplicazione — 3 x 4 = 12 tramite addizioni ripetute (loop con ISZ).
; Equivalente a retronet-4004/examples/moltiplicazione: assemblato produce
; gli stessi byte di testdata/moltiplicazione.rom.

.arch i4004

        LDM 0
        DCL                 ; banco RAM 0
        FIM R0, 0x03        ; R1 = 3 (addendo)
        FIM R2, 0x00        ; indirizzo RAM 0x00
        SRC R2
        LDM 12              ; contatore = 16 - 4
        XCH R4
loop:   ADD R1              ; A += 3
        ISZ R4, loop        ; ripeti 4 volte
        WRM                 ; RAM[0][0][0] = 12
halt:   JUN halt
