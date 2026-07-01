.arch i8080
.com

; Programma .COM CP/M-like: stampa un prompt, legge un carattere e lo ristampa.

.equ BDOS        0x0005
.equ BDOS_READ   1
.equ BDOS_OUT    2
.equ BDOS_PRINT  9
.equ BDOS_TERM   0

        LXI D, prompt
        MVI C, BDOS_PRINT
        CALL BDOS

        MVI C, BDOS_READ
        CALL BDOS        ; A = carattere letto

        MOV E, A
        MVI C, BDOS_OUT
        CALL BDOS

        MVI C, BDOS_TERM
        CALL BDOS

prompt: .byte "Premi un tasto: $"
