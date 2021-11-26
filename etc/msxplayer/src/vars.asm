;constants
music_status_playing: equ 0
music_status_stopped: equ 1

; public vars
main_ram: equ 0xE000
music_status: equ main_ram
wait_cnt: equ music_status + 1     ; frames before interpreting next instructions
music_ip: equ wait_cnt + 1  ; music instruction pointer (bytes)
loop_addr: equ music_ip + 2 ; byte address of the instruction loop
a_volume: equ loop_addr + 1 ; volume status include envelope
b_volume: equ a_volume + 1
c_volume: equ b_volume + 1


