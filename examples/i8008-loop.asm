.arch i8008
; somma COUNT + ... + 1 in A con un loop: B e' il contatore, JFZ ripete finche'
; B non arriva a zero. Con COUNT=5 il risultato e' A = 15 (0x0F). Mostra anche .equ.
;   retronet-asm build examples/i8008-loop.asm -o loop.rom
;   retronet-8008 -bin loop.rom            # dump finale: A=0x0F, Halted=true
.equ COUNT 5
        LBI COUNT       ; B = contatore
        LAI 0           ; A = accumulatore
loop:   ADB             ; A += B
        DCB             ; B-- (aggiorna il flag Zero)
        JFZ loop        ; se Zero e' falso (B != 0) ripeti
        HLT
