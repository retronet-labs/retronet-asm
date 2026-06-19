.arch i4004
; calcolatrice-completa — virgola fissa a 2 decimali, 4 operatori (+ - * /).
; Ogni operando e' memorizzato come intero scalato per 100 (1.5 -> 150, 12 -> 1200).
;  + e -  : dirette (scalato +/- scalato = scalato).
;  *      : reg2 = (A*B)/100  (il prodotto e' scalato per 10000).
;  /      : reg0 *= 100, poi reg2 = reg0/reg1  (quoziente scalato per 100).
; Il risultato (reg2, scalato per 100) e' mostrato con la virgola 2 cifre da destra.
;   echo 1.5+2.25= | retronet-4004 -io calc.rom   ->  3.75
;   echo 7/2=      | retronet-4004 -io calc.rom   ->  3.50
; Input: '.' = nibble 15. Output: 11 = '-', 15 = '.'.
; RAM (banco 0): reg0=A/lavoro, reg1=B, reg2=risultato, reg3=scratch.
        LDM 0
        DCL
        JMS readop      ; A in reg0; tasto finale (operatore) in R8
        LD R8
        XCH R10         ; R10 = operatore
        JMS cp03        ; reg3 = A (parcheggio)
        JMS readop      ; B in reg0; tasto finale ('=') ignorato
        JMS cp01        ; reg1 = B
        JMS cp30        ; reg0 = A
; --- dispatch a 4 vie su R10 ---
        LDM 10
        XCH R5
        LD R10
        STC
        SUB R5
        JCN 0x4, do_add
        LDM 11
        XCH R5
        LD R10
        STC
        SUB R5
        JCN 0x4, do_sub
        LDM 12
        XCH R5
        LD R10
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
; --- sottrazione: reg2 = reg0 - reg1 (con segno) ---
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
        JCN 0x2, sb_pos ; carry==1 -> reg0 >= reg1
        JUN sb_neg
sb_pos: JUN disp
sb_neg: JMS ncomp30     ; reg2 = reg1 - reg0
        FIM R0, 0x10
        FIM R2, 0x30
        FIM R4, 0x20
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
        WMP             ; '-'
        JUN disp
; --- moltiplicazione: reg2 = (reg0*reg1)/100 ---
do_mul: JMS clr2
m_lp:   JMS nz1
        JCN 0x4, m_dn
        JMS add20       ; reg2 += reg0
        JMS dec1        ; reg1--
        JUN m_lp
m_dn:   JMS div10       ; /10
        JMS div10       ; /100
        JUN disp
; --- divisione: reg0 *= 100, poi reg2 = reg0 / reg1 ---
do_div: JMS nz1
        JCN 0x4, dv_zero
        JUN dv_calc
dv_zero: JMS clr2
        JUN disp
dv_calc: JMS mul10       ; reg0 *= 10
        JMS mul10       ; reg0 *= 100
        JMS clr2        ; quoziente reg2 = 0
        JMS ncomp31     ; reg3 = comp9(reg1)
dv_lp:  JMS subtc       ; reg0 -= reg1 ; A=1 se reg0 >= reg1
        JCN 0x2, dv_go  ; carry==1 -> sottrazione valida
        JUN disp        ; prestito -> fine
dv_go:  JMS inc2        ; quoziente++
        JUN dv_lp
; --- display: reg2 scalato per 100 -> "intero.decimali" ---
disp:   JMS pdisp
halt:   JUN halt
; ====================== subroutine ======================
; readop: legge un operando in reg0 (intero scalato per 100); tasto finale in R8
readop: JMS clr0
        LDM 0
        XCH R12         ; seenPoint = 0
        LDM 0
        XCH R13         ; fracCount = 0
ro_lp:  RDR
        XCH R4          ; R4 = tasto
        LDM 15
        XCH R5
        LD R4
        STC
        SUB R5
        JCN 0x4, ro_pt  ; tasto == 15 -> punto decimale
        LDM 10
        XCH R5
        LD R4
        STC
        SUB R5
        JCN 0x2, ro_end ; tasto >= 10 -> fine operando
        LD R12
        JCN 0x4, ro_add ; seenPoint == 0 -> aggiungi (parte intera)
        LDM 2
        XCH R5
        LD R13
        STC
        SUB R5
        JCN 0x4, ro_lp  ; fracCount == 2 -> ignora la cifra
        INC R13         ; fracCount++
ro_add: LD R4
        XCH R14         ; salva la cifra (mul10 azzera R4)
        JMS mul10       ; reg0 *= 10
        FIM R0, 0x00
        SRC R0
        LD R14
        WRM             ; reg0[char0] = cifra
        JUN ro_lp
ro_pt:  LDM 1
        XCH R12         ; seenPoint = 1
        JUN ro_lp
ro_end: LD R4
        XCH R8          ; tasto finale in R8
        LD R13
        XCH R9          ; fracCount in R9
rs_lp:  LDM 2
        XCH R5
        LD R9
        STC
        SUB R5
        JCN 0x4, rs_dn  ; R9 == 2 -> fine scala
        JMS mul10       ; reg0 *= 10
        INC R9
        JUN rs_lp
rs_dn:  BBL 0
; clr0: reg0 = 0
clr0:   FIM R4, 0x00
        FIM R6, 0x80
        LDM 0
clr0l:  SRC R4
        WRM
        INC R5
        ISZ R6, clr0l
        BBL 0
; clr2: reg2 = 0
clr2:   FIM R4, 0x20
        FIM R6, 0x80
        LDM 0
clr2l:  SRC R4
        WRM
        INC R5
        ISZ R6, clr2l
        BBL 0
; mul10: reg0 *= 10 (char_i <- char_{i-1} per i=7..1, char0 <- 0)
mul10:  FIM R2, 0x00
        FIM R0, 0x07
        FIM R6, 0x90
m10l:   LD R1
        DAC
        XCH R3
        SRC R2
        RDM
        SRC R0
        WRM
        LD R1
        DAC
        XCH R1
        ISZ R6, m10l
        FIM R0, 0x00
        SRC R0
        LDM 0
        WRM
        BBL 0
; div10: reg2 /= 10 (char_i <- char_{i+1} per i=0..6, char7 <- 0)
div10:  FIM R4, 0x20
        FIM R2, 0x21
        FIM R6, 0x90
d10l:   SRC R2
        RDM
        SRC R4
        WRM
        INC R5
        INC R3
        ISZ R6, d10l
        FIM R4, 0x27
        SRC R4
        LDM 0
        WRM
        BBL 0
; cp01: reg1 = reg0
cp01:   FIM R0, 0x00
        FIM R2, 0x10
        FIM R6, 0x80
c01l:   SRC R0
        RDM
        SRC R2
        WRM
        INC R1
        INC R3
        ISZ R6, c01l
        BBL 0
; cp03: reg3 = reg0
cp03:   FIM R0, 0x00
        FIM R2, 0x30
        FIM R6, 0x80
c03l:   SRC R0
        RDM
        SRC R2
        WRM
        INC R1
        INC R3
        ISZ R6, c03l
        BBL 0
; cp30: reg0 = reg3
cp30:   FIM R0, 0x00
        FIM R2, 0x30
        FIM R6, 0x80
c30l:   SRC R2
        RDM
        SRC R0
        WRM
        INC R1
        INC R3
        ISZ R6, c30l
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
; add20: reg2 += reg0
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
; dec1: reg1 -= 1 (solo con reg1 != 0)
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
; ncomp31: reg3 = comp9(reg1)
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
; ncomp30: reg3 = comp9(reg0)
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
; inc2: reg2 += 1
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
; pdisp: mostra reg2 scalato per 100 come "intero.decimali"
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
        WMP             ; unita' parte intera: sempre
        LDM 15
        WMP             ; '.'
        FIM R8, 0x21
        SRC R8
        RDM
        WMP             ; decimi
        FIM R8, 0x20
        SRC R8
        RDM
        WMP             ; centesimi
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
