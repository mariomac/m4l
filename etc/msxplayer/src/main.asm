    include "rom/header.asm"
    include "bios.inc"
    include "vars.asm"
    include "psg.asm"
    
    

;----- program start -----
main:
        ; enable channel a
        envelopeA
        envelopeCycle 1200
	channelSet 0b111110
        ld a, 1
        ld (wait_cnt), a

music_loop:
        ld a, (wait_cnt)
        dec a
        jp nz, yield
        
        ; here is where the instruction must be parsed        
        noteA 0x1, 0x40
        envelopeShape 0
        
        ld a, 60
        ld (wait_cnt), a
        halt
        jp music_loop

yield:
        ld (wait_cnt), a
        halt
        jp music_loop

stuff:  jp stuff

    include "rom/tail.asm"
