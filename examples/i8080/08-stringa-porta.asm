.arch i8080

; Stampa una stringa terminata da zero dalla memoria alla porta 1.
; HL scorre il buffer, M legge memoria[HL].

        LXI H, msg
loop:   MOV A, M
        CPI 0
        JZ done
        OUT 1
        INX H
        JMP loop
done:   HLT

msg:    .byte "8080", 0
