.arch i4004
; calcolatrice-due: due operandi multi-cifra + operatore (+ e -).
; Input a registro a scorrimento R0..R3 (R0=MSD) come C1; i due operandi
; vengono depositati in RAM little-endian (char0=unità) e sommati/sottratti
; col loop multi-cifra (riporto/prestito nel flag C), come somma-multicifra.
;   echo 47+58= | retronet-4004 -io calc.rom   ->  105
;   echo 52-27= | retronet-4004 -io calc.rom   ->  25
        LDM 0
        DCL
; --- input primo operando in R0..R3 ---
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
; --- dispatch operatore (C2: + e -) ---
        LDM 10
        XCH R5
        LD R8
        STC
        SUB R5
        JCN 0x4, do_add ; op==10 '+'
        JUN do_sub      ; altrimenti '-'
; --- somma multi-cifra (4 cifre + 5a cifra dal riporto) ---
do_add: FIM R0, 0x00
        FIM R2, 0x10
        FIM R4, 0x20
        FIM R6, 0xC0    ; contatore = 16-4
        CLC
addlp:  SRC R0
        RDM
        SRC R2
        ADM
        DAA
        SRC R4
        WRM
        INC R1
        INC R3
        INC R5
        ISZ R6, addlp
        TCC
        SRC R4
        WRM             ; 5a cifra = riporto finale
        JUN disp
; --- sottrazione multi-cifra (TCS/SBM/ADM/DAA) ---
do_sub: FIM R0, 0x00
        FIM R2, 0x10
        FIM R4, 0x20
        FIM R6, 0xC0
        STC             ; nessun prestito sulla cifra meno significativa
sublp:  TCS
        SRC R2
        SBM
        SRC R0
        ADM
        DAA
        SRC R4
        WRM
        INC R1
        INC R3
        INC R5
        ISZ R6, sublp
        FIM R4, 0x24    ; 5a cifra del risultato = 0
        SRC R4
        LDM 0
        WRM
        JUN disp
; --- display risultato (reg2, char4=MSD .. char0) con soppressione zeri ---
disp:   LDM 0
        XCH R6          ; started = 0
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
        JMS pdig
        LD R6
        JCN 0x4, dzero  ; tutto zero -> "0"
        JUN halt
dzero:  LDM 0
        WMP
halt:   JUN halt
; --- subroutine pdig: stampa la cifra in A sopprimendo gli zeri iniziali ---
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