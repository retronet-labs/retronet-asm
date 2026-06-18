.arch i4004
; calcolatrice-cifre — input multi-cifra (fino a 4) con display a soppressione
; di zeri iniziali. Si digitano le cifre, poi '=' per mostrare.
;   retronet-asm build examples/calcolatrice-cifre.asm -o calc.rom
;   echo 308= | retronet-4004 -io calc.rom        ->  display: 308
        LDM 0
        DCL
        LDM 0
        XCH R0          ; azzera il registro d'ingresso R0..R3
        LDM 0
        XCH R1
        LDM 0
        XCH R2
        LDM 0
        XCH R3
read:   RDR             ; A = tasto
        XCH R4          ; R4 = tasto
        LDM 14
        XCH R5
        LD R4
        STC
        SUB R5          ; A = tasto - 14  (==0 se '=')
        JCN 0x4, show   ; tasto == 14 ('=') -> mostra
        ; shift a sinistra: R0<-R1, R1<-R2, R2<-R3, R3<-nuova cifra
        LD R1
        XCH R0
        LD R2
        XCH R1
        LD R3
        XCH R2
        LD R4
        XCH R3
        JUN read
show:   LDM 0
        XCH R6          ; R6 = flag "ho gia' stampato una cifra" (0 = no)
        LD R0
        JMS pdig
        LD R1
        JMS pdig
        LD R2
        JMS pdig
        LD R3
        JMS pdig
        LD R6
        JCN 0x4, zero   ; nessuna cifra stampata (tutto zero) -> stampa "0"
        JUN halt
zero:   LDM 0
        WMP
halt:   JUN halt
; --- subroutine: stampa la cifra in A sopprimendo gli zeri iniziali ---
pdig:   XCH R7          ; R7 = cifra
        LD R6
        JCN 0x4, pchk   ; started == 0 -> controlla se e' zero iniziale
        LD R7
        WMP             ; gia' iniziato: stampa sempre
        BBL 0
pchk:   LD R7
        JCN 0x4, pskip  ; cifra == 0 e non iniziato -> salta (zero iniziale)
        LDM 1
        XCH R6          ; started = 1
        LD R7
        WMP
        BBL 0
pskip:  BBL 0