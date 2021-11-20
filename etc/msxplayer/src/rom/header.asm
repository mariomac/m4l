; Compilation address
	org 04000H

; ROM header
	db "AB"     ; magic number
	dw main  ; program execution address
	dw 0, 0, 0, 0, 0, 0