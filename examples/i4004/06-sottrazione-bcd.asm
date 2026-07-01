; sottrazione-bcd — sottrazione BCD a cifra singola: 7 - 5 = 2.
; Il prestito vive nel flag C (1 = nessun prestito). Routine TCS+SUB+ADD+DAA.
.arch i4004
        LDM 0
        DCL                 ; banco RAM 0 (serve A=0 prima di DCL)
        FIM R0, 0x00
        SRC R0
        LDM 5
        XCH R1              ; R1 = 5  (sottraendo S)
        LDM 7
        XCH R2              ; R2 = 7  (minuendo M)
        ; D = M - S
        STC                 ; C = 1: nessun prestito in ingresso (cifra unità)
        TCS                 ; A = 10 (base BCD)
        SUB R1              ; A += ~S  (sottrae il sottraendo)
        ADD R2              ; A += M + riporto interno
        DAA                 ; correzione BCD; C = nessun-prestito in uscita
        WRM                 ; RAM[0][0][0] = 2
halt:   JUN halt