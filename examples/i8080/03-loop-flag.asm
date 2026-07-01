.arch i8080

; Loop controllato da flag.
; Somma 5 + 4 + 3 + 2 + 1 in A. DCR aggiorna il flag Z, JNZ continua finche'
; il contatore B non arriva a zero.

        MVI B, 5        ; contatore
        MVI A, 0        ; accumulatore
loop:   ADD B           ; A += B
        DCR B           ; B--; Z=1 se B diventa zero
        JNZ loop        ; continua mentre Z=0
        STA 0x2001      ; risultato: 15 (0x0F)
        HLT
