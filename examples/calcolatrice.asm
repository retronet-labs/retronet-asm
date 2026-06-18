; calcolatrice — calcolatrice a una cifra con 4 operatori (stadio B completo).
; Legge "cifra operatore cifra", calcola + - * / e mostra il risultato (2 cifre).
; Interattiva:  echo 6*7 | retronet-4004 -io calc.rom   ->  42
.arch i4004
        LDM 0
        DCL
        RDR
        XCH R0          ; R0 = primo operando
        RDR
        XCH R3          ; R3 = codice operatore
        RDR
        XCH R1          ; R1 = secondo operando
        ; --- dispatch a 4 vie (confronto = STC + SUB + JCN A==0) ---
        LDM 10
        XCH R4
        LD R3
        STC
        SUB R4
        JCN 0x4, do_add ; operatore == '+' (10)
        LDM 11
        XCH R4
        LD R3
        STC
        SUB R4
        JCN 0x4, do_sub ; == '-' (11)
        LDM 12
        XCH R4
        LD R3
        STC
        SUB R4
        JCN 0x4, do_mul ; == '*' (12)
        JUN do_div      ; altrimenti '/' (13)
do_add: LD R0
        CLC
        ADD R1
        DAA
        XCH R2          ; unità
        TCC             ; decine = riporto
        JUN show
do_sub: STC
        TCS
        SUB R1
        ADD R0
        DAA
        XCH R2          ; unità
        LDM 0           ; decine = 0
        JUN show
do_mul: LDM 0
        XCH R5          ; unità del prodotto = 0
        LDM 0
        XCH R6          ; decine del prodotto = 0
        LDM 0
        XCH R8          ; registro zero (per propagare il riporto)
        STC
        LDM 0
        SUB R1
        XCH R7          ; R7 = 16 - op2 (contatore: op2 addizioni)
mul_lp: LD R5
        CLC
        ADD R0
        DAA
        XCH R5          ; unità += op1
        LD R6
        ADD R8
        DAA
        XCH R6          ; decine += riporto
        ISZ R7, mul_lp
        LD R5
        XCH R2          ; unità
        LD R6           ; decine
        JUN show
do_div: LD R0
        XCH R4          ; resto = dividendo (op1)
        LDM 0
        XCH R5          ; quoziente = 0
div_lp: STC
        TCS
        SUB R1
        ADD R4
        DAA             ; A = resto - divisore, C = no-prestito
        JCN 0x2, div_ok ; C==1 → sottrazione valida
        JUN div_dn
div_ok: XCH R4          ; resto = A
        INC R5          ; quoziente++
        JUN div_lp
div_dn: LD R5
        XCH R2          ; unità = quoziente
        LDM 0           ; decine = 0
show:   WMP             ; display <- decine
        LD R2
        WMP             ; display <- unità
halt:   JUN halt