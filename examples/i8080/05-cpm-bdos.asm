.arch i8080
.com

; Programma CP/M .COM minimale.
; Usa CALL 0005h: C contiene la funzione BDOS, DE punta al buffer.

.equ BDOS       0x0005
.equ BDOS_PRINT 9
.equ BDOS_TERM  0

        LXI D, msg
        MVI C, BDOS_PRINT
        CALL BDOS
        MVI C, BDOS_TERM
        CALL BDOS

msg:    .byte "CIAO DA CP/M$"
