; divisione-bcd — divisione BCD a cifra singola per sottrazioni ripetute: 7 / 2 = 3 r 1.
; Primo programma con salto condizionale (JCN). C=1 dopo la sottrazione => valida
; (committa e Q++); C=0 => prestito, ci si ferma e resta il resto.
.arch i4004
        LDM 0
        DCL                 ; banco 0
        LDM 2
        XCH R1              ; R1 = divisore (Y) = 2
        LDM 7
        XCH R2              ; R2 = resto (R), parte dal dividendo = 7
        LDM 0
        XCH R3              ; R3 = quoziente (Q) = 0
loop:   STC                 ; --- sottrazione tentativa R - Y ---
        TCS
        SUB R1
        ADD R2
        DAA                 ; A = R - Y ; C = 1 se R >= Y (nessun prestito)
        JCN 0x2, commit     ; salta se C==1 (sottrazione valida)
        JUN done            ; C==0: R < Y → finito
commit: XCH R2              ; R = A (nuovo resto); committa solo se valida
        INC R3              ; Q++
        JUN loop
done:   FIM R0, 0x00        ; salva i risultati
        SRC R0
        LD R3
        WRM                 ; RAM[0][0][0] = quoziente
        FIM R0, 0x01
        SRC R0
        LD R2
        WRM                 ; RAM[0][0][1] = resto
halt:   JUN halt