.arch i8080

        MVI A, 0x48      ; 'H'
        OUT 1
        MVI A, 0x49      ; 'I'
        OUT 1
        HLT
