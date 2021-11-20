; Padding to make the file size a multiple of 16K
; (Alternatively, include macros.asm and use ALIGN 4000H)
	ds -$ & 3FFFH
