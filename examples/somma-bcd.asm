; somma-bcd — calcolatrice BCD a cifra singola: 7 + 5.
; Operandi in RAM, somma con ADM, correzione con DAA; la cifra unità (2) va in
; RAM e su Port[0], il riporto resta nel flag C. Stessi byte di
; retronet-4004/testdata/somma-bcd.rom.

.arch i4004

        LDM 0
        DCL
        FIM R0, 0x00        ; cella 0
        SRC R0
        LDM 7
        WRM                 ; RAM[0][0][0] = 7  (operando A)
        FIM R0, 0x01        ; cella 1
        SRC R0
        LDM 5
        WRM                 ; RAM[0][0][1] = 5  (operando B)
        FIM R0, 0x00
        SRC R0
        RDM                 ; A = 7
        CLC                 ; nessun riporto in ingresso
        FIM R0, 0x01
        SRC R0
        ADM                 ; A = 7 + 5 = 12 (0xC)
        DAA                 ; A = 2, C = 1  (correzione BCD)
        FIM R0, 0x02        ; cella 2
        SRC R0
        WRM                 ; RAM[0][0][2] = 2  (cifra unità)
        WMP                 ; Port[0] = 2
halt:   JUN halt
