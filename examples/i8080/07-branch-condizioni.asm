.arch i8080

; Branch condizionali dopo un confronto.
; Classifica A=7 rispetto alla soglia 5 e salva 1 se A >= 5, altrimenti 0.

        MVI A, 7
        CPI 5           ; confronta A con 5: Carry=1 se A < 5
        JC  minore
        MVI A, 1
        JMP done
minore: MVI A, 0
done:   STA 0x2103
        HLT
