; add — addizione binaria semplice: 4 + 3 = 7 nell'accumulatore.

.arch i4004

        LDM 3
        XCH R1          ; R1 = 3
        LDM 4           ; A = 4
        ADD R1          ; A = 4 + 3 = 7
halt:   JUN halt
