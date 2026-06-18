.arch i4004
; calcolatrice-completa — multi-cifra, 4 operatori (+ - * /). La sottrazione
; gestisce il segno negativo (A<B); la divisione produce 2 cifre decimali
; (divisione lunga su resto*10).
;   echo 10/3=  | retronet-4004 -io calc.rom   ->  3.33
;   echo 3-5=   | retronet-4004 -io calc.rom   ->  -2
;   echo 99*99= | retronet-4004 -io calc.rom   ->  9801
; In uscita il nibble 11 = '-' e il nibble 15 = '.' (mappati dal -io).
; Layout RAM (banco 0): reg0=A/resto, reg1=B, reg2=risultato, reg3=scratch.
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
        SUB R5
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
; --- deposita A in RAM reg0 (little-endian) ---
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
        SUB R5
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
        JCN 0x4, do_add
        LDM 11
        XCH R5
        LD R8
        STC
        SUB R5
        JCN 0x4, do_sub
        LDM 12
        XCH R5
        LD R8
        STC
        SUB R5
        JCN 0x4, do_mul
        JUN do_div
; --- somma: reg2 = reg0 + reg1 ---
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
; --- sottrazione: reg2 = reg0 - reg1 (complemento a 10) ---
do_sub: JMS ncomp31
        FIM R0, 0x00
        FIM R2, 0x30
        FIM R4, 0x20
        FIM R6, 0x80
        STC
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
        JCN 0x2, sb_pos ; carry==1 -> reg0 >= reg1 (risultato >= 0)
        JUN sb_neg      ; reg0 < reg1 -> negativo
sb_pos: JUN disp
sb_neg: JMS ncomp30     ; reg3 = comp9(reg0); ricalcola reg2 = reg1 - reg0
        FIM R0, 0x10    ; M = reg1 (B)
        FIM R2, 0x30    ; comp9(reg0)
        FIM R4, 0x20    ; dest = reg2
        FIM R6, 0x80
        STC
sn_lp:  SRC R2
        RDM
        SRC R0
        ADM
        DAA
        SRC R4
        WRM
        INC R1
        INC R3
        INC R5
        ISZ R6, sn_lp
        LDM 11
        WMP             ; segno '-'
        JUN disp
; --- moltiplicazione: reg2 = reg0 * reg1 (addizioni ripetute) ---
do_mul: JMS clr2
m_lp:   JMS nz1
        JCN 0x4, m_dn
        JMS add20
        JMS dec1
        JUN m_lp
m_dn:   JUN disp
; --- divisione decimale: parte intera in reg2, poi 2 decimali on-the-fly ---
do_div: JMS clr2          ; quoziente intero reg2 = 0
        JMS nz1           ; divisore != 0?
        JCN 0x4, dz       ; ==0 -> "0.00"
        JMS ncomp31       ; reg3 = comp9(divisore)
di_lp:  JMS subtc         ; reg0 -= reg1 ; A=1 se reg0 >= reg1
        JCN 0x4, di_dn    ; prestito -> fine parte intera
        JMS inc2          ; quoziente intero ++
        JUN di_lp
di_dn:  JMS addb01        ; ripristina il resto: reg0 += reg1
        JMS pdisp         ; stampa la parte intera (reg2)
        LDM 15
        WMP               ; virgola
        JMS dfrac         ; 1a cifra decimale
        JMS dfrac         ; 2a cifra decimale
        JUN halt
dz:     JMS pdisp         ; reg2 == 0 -> "0"
        LDM 15
        WMP
        LDM 0
        WMP
        LDM 0
        WMP
        JUN halt
; --- display intero (reg2) con soppressione zeri ---
disp:   JMS pdisp
halt:   JUN halt
; ====================== subroutine ======================
; pdisp: stampa reg2 (char7..char1 con soppressione zeri, char0 sempre)
pdisp:  LDM 0
        XCH R6
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
        FIM R8, 0x20
        SRC R8
        RDM
        WMP             ; unità: sempre
        BBL 0
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
; nz1: A = 1 se reg1 != 0
nz1:    FIM R2, 0x10
        FIM R6, 0x80
nz1l:   SRC R2
        RDM
        JCN 0x4, nz1z
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
; addb01: reg0 += reg1 (8 cifre BCD)
addb01: FIM R0, 0x00
        FIM R2, 0x10
        FIM R6, 0x80
        CLC
ab_l:   SRC R2
        RDM
        SRC R0
        ADM
        DAA
        WRM
        INC R1
        INC R3
        ISZ R6, ab_l
        BBL 0
; dec1: reg1 -= 1 (chiamata solo con reg1 != 0)
dec1:   FIM R2, 0x10
        FIM R6, 0x80
        LDM 1
        XCH R7
d1l:    SRC R2
        RDM
        STC
        SUB R7
        JCN 0x2, d1stop
        LDM 9
        WRM
        INC R3
        ISZ R6, d1l
        BBL 0
d1stop: WRM
        BBL 0
; ncomp31: reg3 = complemento a 9 di reg1
ncomp31: FIM R2, 0x10
        FIM R4, 0x30
        FIM R6, 0x80
nc1l:   LDM 9
        STC
        SRC R2
        SBM
        SRC R4
        WRM
        INC R3
        INC R5
        ISZ R6, nc1l
        BBL 0
; ncomp30: reg3 = complemento a 9 di reg0 (per la sottrazione col segno)
ncomp30: FIM R2, 0x00
        FIM R4, 0x30
        FIM R6, 0x80
nc0l:   LDM 9
        STC
        SRC R2
        SBM
        SRC R4
        WRM
        INC R3
        INC R5
        ISZ R6, nc0l
        BBL 0
; subtc: reg0 = reg0 + reg3 + 1 (= reg0 - reg1); A = 1 se reg0 >= reg1
subtc:  FIM R0, 0x00
        FIM R2, 0x30
        FIM R6, 0x80
        STC
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
; mul10: reg0 *= 10 (char_i <- char_{i-1} per i=7..1, char0 <- 0)
mul10:  FIM R2, 0x00    ; sorgente pair (R2=reg0, R3=indice)
        FIM R0, 0x07    ; destinazione pair (R0=reg0, R1=indice=7)
        FIM R6, 0x90    ; contatore 16-7 = 9 -> 7 spostamenti
m10l:   LD R1
        DAC
        XCH R3          ; R3 = R1 - 1 (indice sorgente)
        SRC R2
        RDM
        SRC R0
        WRM             ; char[R1] = char[R1-1]
        LD R1
        DAC
        XCH R1          ; R1--
        ISZ R6, m10l
        FIM R0, 0x00
        SRC R0
        LDM 0
        WRM             ; char0 = 0
        BBL 0
; dfrac: calcola e stampa una cifra decimale (resto reg0 *= 10, conta i reg1)
dfrac:  JMS mul10       ; reg0 *= 10
        LDM 0
        XCH R9          ; R9 = cifra decimale
df_lp:  JMS subtc       ; reg0 -= reg1 ; A=1 se valido
        JCN 0x4, df_dn  ; prestito -> stop
        LD R9
        IAC
        XCH R9          ; cifra++
        JUN df_lp
df_dn:  JMS addb01      ; ripristina il resto reg0 += reg1
        LD R9
        WMP             ; stampa la cifra decimale
        BBL 0
