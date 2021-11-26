    include "rom/header.asm"
    include "bios.inc"
    include "vars.asm"
    include "psg.asm"
    
music_file: incbin "assets/song.bin"

;----- program start -----
main:

init_defaults:
        ; A & B & C channels use raw volume,
        volA 0xF
        volB 0xF
        volC 0xF
        ld a, 15
        ld (a_volume), a  ; todo: not needed variables?
        ld (b_volume), a        
        ld (c_volume), a

	channelSet 0b111000

        ; init instruction pointer, skipping loop address (2 bytes)
        ld hl, music_file
        inc hl
        inc hl
        ld [music_ip], hl
        ld a, music_status_playing
        ld [music_status], a
        ld a, 1
        ld [wait_cnt], a

music_loop:
        ld a, [music_status]
        cp music_status_stopped
        jp z, music_loop        ; todo: do a ret
        ; yield until next frame. TODO: remove
        halt
        ld a, [wait_cnt]
        dec a
        ld [wait_cnt], a
        cp 0
        jp nz, music_loop       ; todo: do a ret

parse_instruction:
        ld hl, [music_ip]       ; read instruction and increase instruction pointer
        ld a, (hl)
        inc hl
        ld [music_ip], hl
        bit 7, a
        jp nz, b1xxxxxxx
b0xxxxxxx:
        bit 6, a
        jp nz, b01xxxxxx
b00xxxxxx:
        bit 5, a
        jp nz, b001xxxxx
b000xxxxx:
        ; there is no 0001xxxx encoding so we directly check if
        ; we set the envelope (a == 0) or wait some minutes
        cp 0
        jp nz, wait 
set_envelope_cycle:   ; assuming a == 0
        jp parse_instruction
b001xxxxx:
        bit 4, a
        jp nz, set_tone_b
set_tone_a:     ; 0010xxxx
        ld      e, a
        ld      a, REG1_A_NOTE_H
        call    BIOS_WRTPSG
        ld      hl, [music_ip]               ; read next entry from stack pointer in e
        ld      e, (hl)
        inc     hl
        ld      [music_ip], hl
        ld      a, REG0_A_NOTE_L
        call    BIOS_WRTPSG
        jp      parse_instruction
b01xxxxxx:
        bit 5, a
        jp nz, b011xxxxx
set_noise_div_rate: ; 010xxxxx
        jp parse_instruction
b011xxxxx:
        bit 4, a
        jp nz, set_tone_c
envelope_wave_shape: ; 0110xxxx        
        jp parse_instruction
b1xxxxxxx:
        bit 6, a
        jp nz, b11xxxxxx
enable_channels: ; 10xxxxxx
        ld      e, a
        ld	a, REG7_CHANNEL_SET
        CALL    BIOS_WRTPSG
        jp parse_instruction
b11xxxxxx:
        bit 5, a
        jp nz, b111xxxxx
b110xxxxx:
        bit 4, a
        jp nz, set_volume_b      
set_volume_a: ; 1100xxxx
        jp parse_instruction          
b111xxxxx:
        bit 4, a
        jp nz, b1111xxxx
set_volume_c: ; 1110xxxx        
        jp parse_instruction
b1111xxxx: 
        bit 3, a
        jp nz, end_song
b11110xxx: ; assuming bit 3 is zero
        bit 1, a
        jp nz, set_envelope_c
b1111000x:
        bit 0, a
        jp nz, set_envelope_b
set_envelope_a: ; 11110000
        jp parse_instruction
wait:
        and 0b00011111                          ; remove instruction code and keep wait cycles
        ld (wait_cnt), a
        jp music_loop        
set_tone_b:
        ld      e, a
        ld      a, REG3_B_NOTE_H
        call    BIOS_WRTPSG
        ld      hl, [music_ip]               ; read next entry from stack pointer in e
        ld      e, (hl)
        inc     hl
        ld      [music_ip], hl
        ld      a, REG2_B_NOTE_L
        call    BIOS_WRTPSG
        jp parse_instruction
set_tone_c:
        ld      e, a
        ld      a, REG5_C_NOTE_H
        call    BIOS_WRTPSG
        ld      hl, [music_ip]               ; read next entry from stack pointer in e
        ld      e, (hl)
        inc     hl
        ld      [music_ip], hl
        ld      a, REG4_C_NOTE_L
        call    BIOS_WRTPSG
        jp parse_instruction
set_volume_b:
        jp parse_instruction
set_envelope_c:
        jp parse_instruction        
set_envelope_b:
        jp parse_instruction
end_song:
        ; check if the loop address is zero. If so, the song ends,
        ; otherwise it loops the music_ip to that address)
        ld      bc, [music_file]        ; ip = music_file + music_file[1:0]
        ld      hl, music_file
        add     hl, bc
        ld      [music_ip], hl
        ld      a, 1                    ; reset wait timer
        ld      [wait_cnt], a
        
        ; if return address is 0, we stop the music later, if not,
        ; let's return to the music loop
        ld      a, b
        cp      0
        jp      nz, music_loop
        ld      a, c
        cp      0
        jp      nz, music_loop
        
        ; todo, if there is loop address, update music_ip
        ; set music status as stopped, disable all channels, and return to loop
        ld      a, music_status_stopped
        ld      [music_status], a
        ld      e, 0b10111111
        ld	a, REG7_CHANNEL_SET
        CALL    BIOS_WRTPSG
        jp      music_loop

stuff:  jp stuff
    include "rom/tail.asm"
