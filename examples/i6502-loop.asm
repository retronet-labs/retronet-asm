.arch i6502
.orgbase $8000

start:
        LDX #$05
        LDA #$00
loop:   CLC
        ADC #$03
        DEX
        BNE loop
        STA $0200

        .org $FFFC
        .word start
