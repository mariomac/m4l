    include "rom/header.asm"
    include "bios.inc"
    include "vars.asm"
    include "psg.asm"
    
music_file: incbin "assets/song.bin"


;----- program start -----
main:
        ; A & B channels use raw volume, C uses envelope
        volA 0xF
        ld a, 15
        ld (a_volume), a
        volC 0xF
        ld (c_volume), a
        envelopeB
        ld a, 0b10000
        ld (b_volume), a
        envelopeCycle 1200
        envelopeShape 0
	channelSet 0b111000

        ; init loop addre
        ld hl, 2
        ld (music_ip), hl

        ld a, 1
        ld (wait_cnt), a

music_loop:
        ld a, (wait_cnt)
        dec a
        jp z, parse_instruction
        ; yield until next frame
        ld (wait_cnt), a
        halt
        jp music_loop

parse_instruction:
        ; here is where the instructions must be parsed
        ld hl, [music_ip]
        ld a, (hl)
        bit 7, a
        jp nz, b1xxxxxxx
b0xxxxxxx:
        bit 6, a
        jp nz, b01xxxxxx
b00xxxxxx:
        bit 5, a
        jp nz, b001xxxxx
b000xxxxx:
        bit 4, a
        jp z, wait 
set_envelope_cycle:   ; assuming a == 0
        jp music_loop
b001xxxxx:
        bit 4, a
        jp nz, set_tone_b
set_tone_a:     ; 0010xxxx
        jp music_loop
b01xxxxxx:
        bit 5, a
        jp nz, b011xxxxx
set_noise_div_rate: ; 010xxxxx
        jp music_loop
b011xxxxx:
        bit 4, a
        jp nz, set_tone_c
envelope_wave_shape: ; 0110xxxx        
        jp music_loop
b1xxxxxxx:
        bit 6, a
        jp nz, b11xxxxxx
enable_channels: ; 10xxxxxx
        jp music_loop

b11xxxxxx:
        bit 5, a
        jp nz, b111xxxxx
b110xxxxx:
        bit 4, a
        jp nz, set_volume_b      
set_volume_a: ; 1100xxxx
        jp music_loop          
b111xxxxx:
        bit 4, a
        jp nz, b1111xxxx
set_volume_c: ; 1110xxxx        
        jp music_loop
b1111xxxx: ; assuming bits 3&2 must be zero
        bit 1, a
        jp nz, set_envelope_c
b1111000x:
        bit 0, a
        jp nz, set_envelope_b
set_envelope_a: ; 11110000
        jp music_loop

        ld a, 60
        ld (wait_cnt), a
        halt
        jp music_loop

wait:
        jp music_loop        
set_tone_b:
        jp music_loop
set_tone_c:
        jp music_loop
set_volume_b:
        jp music_loop
set_envelope_c:
        jp music_loop        
set_envelope_b:
        jp music_loop

stuff:  jp stuff
    include "rom/tail.asm"
