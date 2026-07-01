.arch i8008
; calcolatrice 8008 (binaria, una cifra per operando). Legge "C op C" dal
; terminale ASCII (INP porta 0), calcola + - * / in binario e stampa il
; risultato in decimale (OUT porta 8). L'8008 non ha DAA: si lavora in binario.
;   retronet-asm build examples/i8008/04-calcolatrice.asm -o calc.rom
;   retronet-8008 -bin calc.rom -terminal-input "6*7=" -steps 100000   ->  42
; Nota: la sottrazione assume A >= B (niente segno in questa prima versione).
.equ ZERO  0x30         ; '0'
.equ PLUS  0x2B         ; '+'
.equ MINUS 0x2D         ; '-'
.equ STAR  0x2A         ; '*'
; --- input: cifra1 -> B, operatore -> D, cifra2 -> C ---
        INP 0
        SUI ZERO        ; A = cifra1 - '0'
        LBA             ; B = opA
        INP 0
        LDA             ; D = operatore
        INP 0
        SUI ZERO
        LCA             ; C = opB
; --- dispatch sull'operatore (CPI non modifica A) ---
        LAD
        CPI PLUS
        JTZ do_add
        CPI MINUS
        JTZ do_sub
        CPI STAR
        JTZ do_mul
        JMP do_div      ; default: '/'
; --- somma: A = B + C ---
do_add: LAB
        ADC
        JMP show
; --- sottrazione: A = B - C (assume B >= C) ---
do_sub: LAB
        SUC
        JMP show
; --- moltiplicazione: B * C per addizioni ripetute ---
do_mul: LEI 0           ; E = prodotto
m_lp:   LAC
        CPI 0
        JTZ m_dn        ; C == 0 -> fine
        LAE
        ADB             ; A = E + B
        LEA             ; E += B
        DCC             ; C--
        JMP m_lp
m_dn:   LAE
        JMP show
; --- divisione: B / C per sottrazioni ripetute ---
do_div: LEI 0           ; E = quoziente
        LAC
        CPI 0
        JTZ d_dn        ; divisore 0 -> risultato 0
d_lp:   LAB
        SUC             ; A = B - C
        JTC d_dn        ; prestito (B < C) -> fine
        LBA             ; B = B - C
        INE             ; quoziente++
        JMP d_lp
d_dn:   LAE
        JMP show
; --- display decimale del valore in A (0..99), zero iniziale soppresso ---
show:   LEI 0           ; E = decine
s_lp:   CPI 10
        JTC s_un        ; A < 10 -> resto = unita'
        SUI 10
        INE             ; decine++
        JMP s_lp
s_un:   LCA             ; C = unita'
        LAE
        CPI 0
        JTZ s_only      ; decine == 0 -> stampa solo le unita'
        ADI ZERO
        OUT 8           ; stampa la cifra delle decine
s_only: LAC
        ADI ZERO
        OUT 8           ; stampa la cifra delle unita'
halt:   HLT
