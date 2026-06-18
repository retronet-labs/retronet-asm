.arch i4004
; calcolatrice-completa — calcolatrice multi-cifra con i 4 operatori (+ - * /).
; Legge "OPER1 op OPER2 =", ciascun operando fino a 4 cifre dalla tastiera,
; calcola in BCD a 8 cifre nella RAM e mostra il risultato con soppressione
; degli zeri iniziali.
;   echo 12*12= | retronet-4004 -io calc.rom   ->  144
;   echo 144/12= | retronet-4004 -io calc.rom  ->  12
; Layout RAM (banco 0): reg0=A, reg1=B, reg2=risultato, reg3=scratch.
; La sottrazione usa il complemento a 10: M - S = M + comp9(S) + 1 (scarta il
; riporto finale), riusando l'addizionatore BCD. Il riporto finale vale 1 se
; M >= S: la divisione lo usa come confronto.
        LDM 0
        DCL
; --- input primo operando in R0..R3 (R0=MSD) ---
        LDM 0
        XCH R0
        LDM 0
        XCH R1
        LDM 0
        XCH R2
        LDM 0
        XCH R3
readA:  RDR
        XCH R4
        LDM 10
        XCH R5
        LD R4
        STC
        SUB R5          ; C=1 se tasto >= 10 (operatore)
        JCN 0x2, opA
        LD R1
        XCH R0
        LD R2
        XCH R1
        LD R3
        XCH R2
        LD R4
        XCH R3
        JUN readA
opA:    LD R4
        XCH R8          ; R8 = operatore (10..13)
; --- deposita A in RAM reg0 (little-endian: char0=R3 ... char3=R0) ---
        FIM R10, 0x00
        SRC R10
        LD R3
        WRM
        FIM R10, 0x01
        SRC R10
        LD R2
        WRM
        FIM R10, 0x02
        SRC R10
        LD R1
        WRM
        FIM R10, 0x03
        SRC R10
        LD R0
        WRM
; --- input secondo operando in R0..R3 ---
        LDM 0
        XCH R0
        LDM 0
        XCH R1
        LDM 0
        XCH R2
        LDM 0
        XCH R3
readB:  RDR
        XCH R4
        LDM 10
        XCH R5
        LD R4
        STC
        SUB R5          ; C=1 se tasto >= 10 ('=')
        JCN 0x2, opB
        LD R1
        XCH R0
        LD R2
        XCH R1
        LD R3
        XCH R2
        LD R4
        XCH R3
        JUN readB
; --- deposita B in RAM reg1 ---
opB:    FIM R10, 0x10
        SRC R10
        LD R3
        WRM
        FIM R10, 0x11
        SRC R10
        LD R2
        WRM
        FIM R10, 0x12
        SRC R10
        LD R1
        WRM
        FIM R10, 0x13
        SRC R10
        LD R0
        WRM
; --- dispatch a 4 vie sull'operatore in R8 ---
        LDM 10
        XCH R5
        LD R8
        STC
        SUB R5
        JCN 0x4, do_add ; == '+'
        LDM 11
        XCH R5
        LD R8
        STC
        SUB R5
        JCN 0x4, do_sub ; == '-'
        LDM 12
        XCH R5
        LD R8
        STC
        SUB R5
        JCN 0x4, do_mul ; == '*'
        JUN do_div      ; altrimenti '/'
; --- somma: reg2 = reg0 + reg1 (8 cifre) ---
do_add: FIM R0, 0x00
        FIM R2, 0x10
        FIM R4, 0x20
        FIM R6, 0x80
        CLC
a_lp:   SRC R0
        RDM
        SRC R2
        ADM
        DAA
        SRC R4
        WRM
        INC R1
        INC R3
        INC R5
        ISZ R6, a_lp
        JUN disp
; --- sottrazione: reg2 = reg0 - reg1 (complemento a 10, assume A>=B) ---
do_sub: JMS ncomp31     ; reg3 = comp9(reg1)
        FIM R0, 0x00    ; M = reg0
        FIM R2, 0x30    ; reg3 = comp9(B)
        FIM R4, 0x20    ; dest = reg2
        FIM R6, 0x80
        STC             ; +1 (completa il complemento a 10)
sb_lp:  SRC R2
        RDM
        SRC R0
        ADM
        DAA
        SRC R4
        WRM
        INC R1
        INC R3
        INC R5
        ISZ R6, sb_lp
        JUN disp        ; riporto finale scartato
; --- moltiplicazione: reg2 = reg0 * reg1 (addizioni ripetute) ---
do_mul: JMS clr2        ; reg2 = 0
m_lp:   JMS nz1         ; A = 1 se reg1 != 0
        JCN 0x4, m_dn   ; reg1 == 0 -> fine
        JMS add20       ; reg2 += reg0
        JMS dec1        ; reg1 -= 1
        JUN m_lp
m_dn:   JUN disp
; --- divisione: reg2 = reg0 / reg1 (sottrazioni ripetute) ---
do_div: JMS clr2        ; quoziente reg2 = 0
        JMS nz1         ; divisore != 0?
        JCN 0x4, d_dn   ; reg1 == 0 -> risultato 0
        JMS ncomp31     ; reg3 = comp9(reg1) (una volta sola)
d_lp:   JMS subtc       ; reg0 = reg0 - reg1 ; A = 1 se reg0 >= reg1
        JCN 0x4, d_dn   ; prestito -> stop
        JMS inc2        ; quoziente++
        JUN d_lp
d_dn:   JUN disp
; --- display: reg2 char7..char0 con soppressione zeri iniziali ---
disp:   LDM 0
        XCH R6          ; started = 0
        FIM R8, 0x27
        SRC R8
        RDM
        JMS pdig
        FIM R8, 0x26
        SRC R8
        RDM
        JMS pdig
        FIM R8, 0x25
        SRC R8
        RDM
        JMS pdig
        FIM R8, 0x24
        SRC R8
        RDM
        JMS pdig
        FIM R8, 0x23
        SRC R8
        RDM
        JMS pdig
        FIM R8, 0x22
        SRC R8
        RDM
        JMS pdig
        FIM R8, 0x21
        SRC R8
        RDM
        JMS pdig
        ; cifra delle unità (char0): sempre stampata (anche se 0)
        FIM R8, 0x20
        SRC R8
        RDM
        WMP
halt:   JUN halt
; ====================== subroutine ======================
; pdig: stampa la cifra in A sopprimendo gli zeri iniziali (flag R6)
pdig:   XCH R7
        LD R6
        JCN 0x4, pchk
        LD R7
        WMP
        BBL 0
pchk:   LD R7
        JCN 0x4, pskip
        LDM 1
        XCH R6
        LD R7
        WMP
        BBL 0
pskip:  BBL 0
; clr2: reg2 = 0 (8 cifre)
clr2:   FIM R4, 0x20
        FIM R6, 0x80
        LDM 0
clr2l:  SRC R4
        WRM
        INC R5
        ISZ R6, clr2l
        BBL 0
; nz1: A = 1 se reg1 != 0, altrimenti A = 0
nz1:    FIM R2, 0x10
        FIM R6, 0x80
nz1l:   SRC R2
        RDM
        JCN 0x4, nz1z   ; cifra == 0 -> continua
        BBL 1
nz1z:   INC R3
        ISZ R6, nz1l
        BBL 0
; add20: reg2 += reg0 (8 cifre BCD)
add20:  FIM R0, 0x00
        FIM R4, 0x20
        FIM R6, 0x80
        CLC
a20l:   SRC R0
        RDM
        SRC R4
        ADM
        DAA
        WRM
        INC R1
        INC R5
        ISZ R6, a20l
        BBL 0
; dec1: reg1 -= 1 (chiamata solo con reg1 != 0) — prestito sui soli zeri di coda
dec1:   FIM R2, 0x10
        FIM R6, 0x80
        LDM 1
        XCH R7          ; R7 = prestito da sottrarre (1)
d1l:    SRC R2
        RDM
        STC
        SUB R7          ; A = cifra - prestito ; C=1 se nessun prestito
        JCN 0x2, d1stop
        LDM 9
        WRM             ; cifra era 0 -> 9, continua il prestito
        INC R3
        ISZ R6, d1l
        BBL 0
d1stop: WRM             ; scrive cifra-1; le cifre superiori restano invariate
        BBL 0
; ncomp31: reg3 = complemento a 9 di reg1 (reg3[i] = 9 - reg1[i])
ncomp31: FIM R2, 0x10
        FIM R4, 0x30
        FIM R6, 0x80
nc1l:   LDM 9
        STC
        SRC R2
        SBM             ; A = 9 + ~reg1[i] + 1 = 9 - reg1[i] (mod 16)
        SRC R4
        WRM
        INC R3
        INC R5
        ISZ R6, nc1l
        BBL 0
; subtc: reg0 = reg0 + reg3 + 1 (= reg0 - reg1 via comp.a10); A = 1 se reg0 >= reg1
subtc:  FIM R0, 0x00
        FIM R2, 0x30
        FIM R6, 0x80
        STC             ; +1
stcl:   SRC R2
        RDM
        SRC R0
        ADM
        DAA
        WRM
        INC R1
        INC R3
        ISZ R6, stcl
        JCN 0x2, stok
        BBL 0
stok:   BBL 1
; inc2: reg2 += 1 (incremento BCD a 8 cifre)
inc2:   FIM R4, 0x20
        FIM R6, 0x80
        STC
i2l:    SRC R4
        RDM
        XCH R7
        TCC
        ADD R7
        DAA
        WRM
        INC R5
        ISZ R6, i2l
        BBL 0
