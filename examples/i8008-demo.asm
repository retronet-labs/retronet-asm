.arch i8008
; demo i8008: istruzioni a 1 byte senza operandi.
; Mostra trasferimenti, ALU-registro, increment/decrement, rotate e l'arresto.
; (Senza immediati i registri partono da 0: qui conta la codifica, non i valori.)
;   retronet-asm build examples/i8008-demo.asm -o demo.rom
;   retronet-8008 -bin demo.rom -disasm 8     # ri-disassembla gli stessi mnemonici
        LBA             ; B <- A
        LCB             ; C <- B
        ADB             ; A <- A + B
        ADC             ; A <- A + C
        INB             ; B++
        DCC             ; C--
        RLC             ; ruota A a sinistra
        HLT
