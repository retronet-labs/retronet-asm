.arch i8008
; somma 5 + 4 + 3 + 2 + 1 in A con un loop: B e' il contatore, JFZ ripete
; finche' B non arriva a zero. Risultato atteso: A = 15 (0x0F), B = 0.
;   retronet-asm build examples/i8008-loop.asm -o loop.rom
;   retronet-8008 -bin loop.rom            # dump finale: A=0x0F, Halted=true
        LBI 5           ; B = contatore
        LAI 0           ; A = accumulatore
loop:   ADB             ; A += B
        DCB             ; B-- (aggiorna il flag Zero)
        JFZ loop        ; se Zero e' falso (B != 0) ripeti
        HLT
