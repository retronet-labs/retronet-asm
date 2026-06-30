.arch i6502
.orgbase $8000

start:
        SED
        CLC
        LDA #$45
        ADC #$55
        STA $0201

        .org $FFFC
        .word start
