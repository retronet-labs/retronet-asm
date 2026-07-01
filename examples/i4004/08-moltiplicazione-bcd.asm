; moltiplicazione-bcd — M x N per addizioni ripetute: 25 x 5 = 125.
; Loop esterno N volte; a ogni giro somma M (BCD 3 cifre) all'accumulatore P.
; P parte da 0 (la RAM è azzerata all'avvio).
.arch i4004
        LDM 0
        DCL
        ; M = 25 in reg0: char0=5 (unità), char1=2 (decine), char2=0
        FIM R0, 0x00
        SRC R0
        LDM 5
        WRM
        FIM R0, 0x01
        SRC R0
        LDM 2
        WRM
        ; contatore esterno: R7 = 16 - N = 16 - 5 = 11 (0x0B)
        FIM R6, 0x0B
outer:  FIM R0, 0x00        ; puntatore M -> char0
        FIM R2, 0x10        ; puntatore P -> char0 (reg1)
        FIM R4, 0xD0        ; contatore interno = 16 - 3 = 13 (3 cifre)
        CLC                 ; nessun riporto a inizio somma
inner:  SRC R2
        RDM                 ; A = cifra di P
        SRC R0
        ADM                 ; A += cifra di M + riporto
        DAA                 ; correzione BCD
        SRC R2
        WRM                 ; rimetti la cifra in P
        INC R1              ; avanza puntatore M
        INC R3              ; avanza puntatore P
        ISZ R4, inner       ; ripeti per 3 cifre
        ISZ R7, outer       ; ripeti N volte
halt:   JUN halt