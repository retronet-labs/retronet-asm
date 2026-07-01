.arch i6502
.orgbase $8000

; Modi di indirizzamento 6502: immediato, assoluto indicizzato e indiretto
; indicizzato. Il puntatore in zero page $00/$01 viene caricato con l'indirizzo
; della tabella.

reset:  LDX #$00
        LDY #$02

        LDA #<tabella
        STA $00
        LDA #>tabella
        STA $01

        LDA tabella,X   ; absolute,X: primo elemento
        STA $0200
        LDA ($00),Y     ; indirect indexed: terzo elemento
        STA $0201

halt:   JMP halt

tabella:
        .byte $11, $22, $33, $44

        .org $FFFC
        .word reset
