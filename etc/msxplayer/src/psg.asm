REG0_A_NOTE_L: EQU 0
REG1_A_NOTE_H: EQU 1
REG2_B_NOTE_L: EQU 2
REG3_B_NOTE_H: EQU 3
REG4_C_NOTE_L: EQU 4
REG5_C_NOTE_H: EQU 5
REG6_NOISE_FREQ: EQU 6
REG7_CHANNEL_SET:   EQU 7
REG8_A_VOLUME: EQU 8
REG9_B_VOLUME: EQU 9
REG10_C_VOLUME: EQU 10
REG11_ENV_CYCLE_L: EQU 11
REG12_ENV_CYCLE_H: EQU 12
REG13_ENV_SHAPE: EQU 13

// https://github.com/Konamiman/MSX2-Technical-Handbook/blob/master/md/Chapter5a.md#1-psg-and-sound-output

; enbles tone and disables noise
; argument: address of the variable containing channels data
channelSet: MACRO set
	ld	a, REG7_CHANNEL_SET
        ld      e, set
        CALL    BIOS_WRTPSG
ENDM
; vol functions disable envelope
volA: MACRO vol
        ld      a, REG8_A_VOLUME
        ld      e, vol
        call    BIOS_WRTPSG
ENDM
volB: MACRO vol
        ld      a, REG9_B_VOLUME
        ld      e, vol
        call    BIOS_WRTPSG
ENDM
volC: MACRO vol
        ld      a, REG10_C_VOLUME
        ld      e, vol
        call    BIOS_WRTPSG
ENDM

; cycle: 16 bits
envelopeCycle: MACRO cycle
        ld      a, REG11_ENV_CYCLE_L
        ld      e, cycle AND 0xF
        call    BIOS_WRTPSG
        ld      a, REG12_ENV_CYCLE_H
        ld      e, cycle >> 8
        call    BIOS_WRTPSG
ENDM

; pattern: 4 bits
envelopeShape: MACRO shape        
        ld      a, REG13_ENV_SHAPE
        ld      e, shape
        call    BIOS_WRTPSG
ENDM

envelopeA: MACRO
        ld      a, REG8_A_VOLUME
        ld      e, 0b10000
        call    BIOS_WRTPSG
ENDM
envelopeB: MACRO
        ld      a, REG9_B_VOLUME
        ld      e, 0b10000
        call    BIOS_WRTPSG
ENDM
envelopeC: MACRO
        ld      a, REG10_C_VOLUME
        ld      e, 0b10000
        call    BIOS_WRTPSG
ENDM

noteA: MACRO hi, lo
        ld      a, REG0_A_NOTE_L
        ld      e, lo
        call    BIOS_WRTPSG
        ld      a, REG1_A_NOTE_H
        ld      e, hi
        call    BIOS_WRTPSG
ENDM

noteB: MACRO hi, lo
        ld      a, REG2_B_NOTE_L
        ld      e, lo
        call    BIOS_WRTPSG
        ld      a, REG3_B_NOTE_H
        ld      e, hi
        call    BIOS_WRTPSG
ENDM
