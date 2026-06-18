; somma-multicifra — addizione BCD multi-cifra: 47 + 58 = 105.
; Cifre in RAM (little-endian), loop su tre puntatori paralleli (A, B,
; risultato); il riporto si propaga nel flag C tra le iterazioni, TCC
; trasforma l'ultimo riporto nella cifra delle centinaia. Stessi byte di
; retronet-4004/testdata/somma-multicifra.rom.

        LDM 0
        DCL
        ; A = 47 nel registro RAM 0: char0 = 7 (unità), char1 = 4 (decine)
        FIM R0, 0x00
        SRC R0
        LDM 7
        WRM
        FIM R0, 0x01
        SRC R0
        LDM 4
        WRM
        ; B = 58 nel registro RAM 1
        FIM R2, 0x10
        SRC R2
        LDM 8
        WRM
        FIM R2, 0x11
        SRC R2
        LDM 5
        WRM
        ; puntatori a char 0 (A=reg0, B=reg1, risultato=reg2) e contatore = 16-2
        FIM R0, 0x00
        FIM R2, 0x10
        FIM R4, 0x20
        FIM R6, 0xE0
        CLC
loop:   SRC R0
        RDM                 ; A = cifra di A
        SRC R2
        ADM                 ; A += cifra di B + riporto
        DAA                 ; correzione BCD
        SRC R4
        WRM                 ; scrivi la cifra del risultato
        INC R1
        INC R3
        INC R5
        ISZ R6, loop        ; ripeti per 2 cifre
        TCC                 ; ultimo riporto -> cifra centinaia
        SRC R4
        WRM
halt:   JUN halt
