.arch i6502
.orgbase $8000

; ROM 6502 minima con vettore di reset.
; Inizializza lo stack, scrive un valore in memoria e resta in un loop stabile.

reset:  SEI             ; disabilita interrupt mascherabili
        CLD             ; modalita' binaria per ADC/SBC
        LDX #$FF
        TXS             ; stack pointer = $FF
        LDA #$01
        STA $0200       ; segnale osservabile in memoria
halt:   JMP halt

        .org $FFFC
        .word reset     ; vettore RESET little-endian
