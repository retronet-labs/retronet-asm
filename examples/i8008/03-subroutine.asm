.arch i8008
; chiamata di subroutine: dbl raddoppia A. Qui 9 -> 18 (0x12).
; Mostra CAL/RET e lo stack hardware dell'8008 (il PC corrente e' in cima).
;   retronet-asm build examples/i8008/03-subroutine.asm -o sub.rom
;   retronet-8008 -bin sub.rom             # dump finale: A=0x12, Halted=true
        LAI 9           ; A = 9
        CAL dbl         ; A = dbl(A)
        HLT
dbl:    ADA             ; A += A  (= A * 2)
        RET
