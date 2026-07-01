; calcolatrice-addsub — legge "cifra operatore cifra" e calcola + o - (stadio B).
; Esempio:  echo 7+5 | retronet-4004 -io calc.rom   ->  12
;           echo 9-4 | retronet-4004 -io calc.rom   ->  05
.arch i4004
        LDM 0
        DCL
        RDR
        XCH R0          ; R0 = primo operando
        RDR
        XCH R3          ; R3 = codice operatore
        RDR
        XCH R1          ; R1 = secondo operando
        ; --- dispatch: operatore == '+' (10) ? ---
        LDM 10
        XCH R4
        LD R3
        STC
        SUB R4          ; A = operatore - 10
        JCN 0x4, do_add ; A==0 → addizione; altrimenti sottrazione
do_sub: STC             ; D = op1 - op2  (routine a cifra singola)
        TCS
        SUB R1
        ADD R0
        DAA             ; A = unità, C = no-prestito
        XCH R2
        LDM 0           ; decine = 0 (risultato 0-9 se op1 >= op2)
        JUN show
do_add: LD R0           ; D = op1 + op2
        CLC
        ADD R1
        DAA             ; A = unità, C = riporto
        XCH R2
        TCC             ; A = decine (riporto)
show:   WMP             ; display <- decine
        LD R2
        WMP             ; display <- unità
halt:   JUN halt