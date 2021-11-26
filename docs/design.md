# Design document

## Syntax

```
; comments start with ;
; constants are defined with a let. They can be instruments or tablature pieces

; program header must start with a set of "key value" properties

; tempo, specified as the number of beats per minute (1 quarter/crotchet == 1 beat)
tempo 120
; for retro-machines, the destination refresh rate (50 or 60 Hz) must be specified
; to properly calculate the tempo
psg.hz 60

# variables start with $ and assigning an instrument or tablature uses the `:=`symbol
$instrument1 := psg {
    pattern: 2
}
$instrument2 := psg {
    frequency: 123
    noise: 12
    cycle: 23
    pattern: 4
}
$piece := o4 e8e8 r8 c8 e

; set channels instruments, combine variables and tablature literals. Constants are read with an $

@channel1 <- $instrument1 $piece
@channel2 <- $instrument2 r16 $piece

; sync barrier. Music doesn't continue until all channels have finished (two dash at least) 

--

; loop can include channels and sync barriers. It is an infinite loop, so it does not have sense
; nest loops or put anything after the loop 
; loop tag also acts as a synced block
loop:

@channel3 <- a1 b2 c3 c4 c5
```

## Grammar

```
program := header? constantDef* statement* ('loop:' statement*)?
header := (KEY VAL\n)* 

constantDef := ID ':=' (instrumentDef | tablature)

instrumentDef := CLASS '{' mapEntry* '}'

tablature := (ID | NOTE | SILENCE | OCTAVE | INCOCT | DECOCT | tuplet | '|')+

tuplet := '(' (NOTE|OCTAVE|INCOCT|DECOCT) + ')' NUM

ID := $(\w)+

statement := channelFill | SYNC 
channelFill := CHANNEL_ID '<-' tablature+
SYNC := '-'*

```

## Binary compilation for MSX PSG

`
loop start byte (2 bytes) == 0 if no loop
tablature instructions (variable-byte encoding)
`

### Tablature instruction set

Instructions have variable bit size

* `00000000 hhhhhhhh llllllll` set envelope cycle
* `000nnnnn` wait N refresh signals (up to 32), n must be != 0
* `0010hhhh llllllll` Set tone CH A
* `0011hhhh llllllll` Set tone CH B
* `0111hhhh llllllll` Set tone CH C
* `0110nnnn` Envelope wave shape
* `010nnnnn` Set noise div rate.
* `10cbaCBA` Enable channels
  - `cba` noise on channels c,b,a
  - `CBA` tone on channels c,b,a
* `1100vvvv` set volume A
* `1101vvvv` set volume B
* `1110vvvv` set volume C
* `11110000` set envelope for A (ignore volume)
* `11110001` set envelope for B (ignore volume)
* `11110010` set envelope for C (ignore volume)
* `11111xxx` song end (finish or jump to loop)