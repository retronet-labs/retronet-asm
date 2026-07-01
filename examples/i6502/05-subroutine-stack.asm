.arch i6502
.orgbase $8000

; Subroutine con stack: triple calcola A * 3.
; PHA conserva l'operando originale, ASL produce 2*A, poi ADC somma l'originale.

reset:  LDX #$FF
        TXS
        CLD
        LDA #$05
        JSR triple
        STA $0202       ; risultato: $0F
halt:   JMP halt

triple: PHA             ; salva A originale sullo stack
        ASL A           ; A = 2*A
        STA $20         ; scratch in zero page
        PLA             ; A = valore originale
        CLC
        ADC $20         ; A = A + 2*A
        RTS

        .org $FFFC
        .word reset
