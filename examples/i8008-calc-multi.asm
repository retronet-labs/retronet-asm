.arch i8008
; calcolatrice 8008 multi-cifra (binaria, 8 bit: operandi e risultato 0..255).
; Legge "N op N" dal terminale ASCII (INP 0), calcola + - * / in binario e stampa
; il risultato in decimale (OUT 8). L'8008 e' a 8 bit e non ha DAA: aritmetica
; binaria (contrasto con la calcolatrice 4004, che e' BCD).
;   retronet-asm build examples/i8008-calc-multi.asm -o calc.rom
;   retronet-8008 -bin calc.rom -terminal-input "12*12=" -steps 2000000   ->  144
; Limiti: valori e risultato 0..255; la sottrazione assume N1 >= N2; overflow ignorato.
.equ ZERO  0x30         ; '0'
.equ NINE1 0x3A         ; ':' = subito dopo '9'
.equ PLUS  0x2B
.equ MINUS 0x2D
.equ STAR  0x2A
; --- input dei due operandi ---
        CAL readnum    ; B = opA, A = operatore
        LDA             ; D = operatore
        LEB             ; E = opA
        CAL readnum    ; B = opB, A = '=' (ignorato)
; --- dispatch sull'operatore (opA=E, opB=B) ---
        LAD
        CPI PLUS
        JTZ do_add
        CPI MINUS
        JTZ do_sub
        CPI STAR
        JTZ do_mul
        JMP do_div      ; default: '/'
do_add: LAE
        ADB             ; A = opA + opB
        JMP show
do_sub: LAE
        SUB             ; A = opA - opB
        JMP show
do_mul: LCI 0           ; C = prodotto
m_lp:   LAB
        CPI 0
        JTZ m_dn        ; opB esaurito
        LAC
        ADE             ; A = C + opA
        LCA             ; C += opA
        DCB             ; opB--
        JMP m_lp
m_dn:   LAC
        JMP show
do_div: LCI 0           ; C = quoziente
        LAB
        CPI 0
        JTZ d_dn        ; divisore 0 -> risultato 0
d_lp:   LAE
        SUB             ; A = resto - opB
        JTC d_dn        ; prestito (resto < opB) -> fine
        LEA             ; resto -= opB
        INC             ; quoziente++
        JMP d_lp
d_dn:   LAC
        JMP show
; --- display decimale di A (0..255): centinaia C, decine D, unita' E ---
show:   LCI 0
sh_h:   CPI 100
        JTC sh_hd
        SUI 100
        INC             ; centinaia++
        JMP sh_h
sh_hd:  LDI 0
sh_t:   CPI 10
        JTC sh_td
        SUI 10
        IND             ; decine++
        JMP sh_t
sh_td:  LEA             ; E = unita'
        LBI 0           ; B = flag "gia' stampata una cifra"
        LAC             ; centinaia
        CPI 0
        JTZ p_t
        ADI ZERO
        OUT 8
        LBI 1
p_t:    LAB
        CPI 0
        JFZ p_tn        ; gia' stampata -> stampa le decine comunque
        LAD
        CPI 0
        JTZ p_u         ; decine 0 e niente prima -> salta
p_tn:   LAD
        ADI ZERO
        OUT 8
p_u:    LAE
        ADI ZERO
        OUT 8           ; unita': sempre
halt:   HLT
; --- readnum: B = numero letto da INP 0; A = carattere terminatore ---
readnum: LBI 0
rn_lp:  INP 0
        CPI ZERO
        JTC rn_end      ; A < '0' -> terminatore
        CPI NINE1
        JFC rn_end      ; A >= ':' -> terminatore
        SUI ZERO        ; A = cifra 0..9
        LLA             ; L = cifra
        LAB
        ADA             ; 2B
        LHA             ; H = 2B
        ADA             ; 4B
        ADA             ; 8B
        ADH             ; 8B + 2B = 10B
        ADL             ; + cifra
        LBA             ; B = B*10 + cifra
        JMP rn_lp
rn_end: RET
