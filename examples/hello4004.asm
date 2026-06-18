; hello4004 — scrive il valore 5 sulla porta di output del banco RAM 0.
; È il "ciao mondo" del 4004: nessun calcolo, solo un output e l'arresto.

        LDM 0
        DCL             ; seleziona il banco RAM 0
        LDM 5
        WMP             ; Port[0] = 5
halt:   JUN halt        ; arresto (salto su se stesso)
