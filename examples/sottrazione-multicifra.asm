; sottrazione-multicifra — sottrazione BCD multi-cifra: 52 - 27 = 25.
; Routine per cifra TCS->SBM->ADM->DAA in loop su RAM; il prestito si propaga
; nel flag C (1 = nessun prestito), come il riporto ma al contrario.
.arch i4004
        LDM 0
        DCL
        ; M = 52 nel registro RAM 0: char0=2 (unità), char1=5 (decine)
        FIM R0, 0x00
        SRC R0
        LDM 2
        WRM
        FIM R0, 0x01
        SRC R0
        LDM 5
        WRM
        ; S = 27 nel registro RAM 1
        FIM R2, 0x10
        SRC R2
        LDM 7
        WRM
        FIM R2, 0x11
        SRC R2
        LDM 2
        WRM
        ; puntatori (M=reg0, S=reg1, risultato=reg2) a char 0; contatore = 16-2
        FIM R0, 0x00
        FIM R2, 0x10
        FIM R4, 0x20
        FIM R6, 0xE0
        STC                 ; C = 1: nessun prestito sulla cifra meno significativa
loop:   TCS                 ; A = 10/9 secondo il prestito in ingresso
        SRC R2
        SBM                 ; A += ~S(RAM)  → sottrae la cifra del sottraendo
        SRC R0
        ADM                 ; A += M(RAM) + riporto interno
        DAA                 ; correzione BCD; C = nessun-prestito in uscita
        SRC R4
        WRM                 ; scrivi la cifra del risultato
        INC R1
        INC R3
        INC R5
        ISZ R6, loop        ; ripeti per 2 cifre
halt:   JUN halt