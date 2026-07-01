.arch i8080

; Aritmetica elementare su registri 8080.
; Calcola 5 + 7, lascia il risultato in A e lo copia in memoria a 0x2000.

        MVI B, 5        ; B = primo addendo
        MVI C, 7        ; C = secondo addendo
        MOV A, B        ; A = B: l'ALU lavora sempre sull'accumulatore
        ADD C           ; A = A + C = 12
        STA 0x2000      ; memoria[0x2000] = 12
        HLT
