.arch i8086
.orgbase 0x7C00

; Boot sector minimo: codice valido, loop stabile e firma 55 AA.
; Serve come base per capire .orgbase e la posizione della firma.

start:  cli
halt:   hlt
        jmp halt

        .org 0x7DFE
        .byte 0x55, 0xAA
