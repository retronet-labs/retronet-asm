; calcolatrice-add — somma di due cifre lette dalla tastiera (stadio A).
; Firmware per retronet-4004 in modalità -io: legge 2 tasti con RDR, somma in
; BCD, e manda al display (WMP) le due cifre del risultato (decine, unità).
;
;   retronet-asm build examples/calcolatrice-add.asm -o calc.rom
;   echo 75 | retronet-4004 -io calc.rom        ->  display: 12
.arch i4004
        LDM 0
        DCL             ; banco RAM 0
        RDR             ; A = primo tasto
        XCH R0          ; R0 = primo operando
        RDR             ; A = secondo tasto
        XCH R1          ; R1 = secondo operando
        LD R0
        CLC
        ADD R1          ; A = R0 + R1
        DAA             ; A = unità, C = riporto
        XCH R2          ; salva la cifra unità
        TCC             ; A = riporto -> cifra decine
        WMP             ; display <- decine
        LD R2
        WMP             ; display <- unità
halt:   JUN halt
