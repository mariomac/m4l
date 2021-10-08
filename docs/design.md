# Design document
## Syntax
```
# comments start with sharp #
# constants are defined with a let. They can be instruments or tablature pieces
# variables start with $ and assigning an instrument or tablature uses the `:=`symbol
$instrument1 := {
    wave: sine
}
$instrument2 := {
    wave: square
    adsr: 5->100, 20->60, 25, 30
}
$piece := o4 e8e8 r8 c8 e

# set channels instruments, combine variables and tablature literals. Constants are read with an $

@channel1 <- $instrument1 $piece
@channel2 <- $instrument2 r16 $piece

# sync barrier. Music doesn't continue until all channels have finished (two dash at least) 

--

# loop can include channels and sync barriers. It is an infinite loop, so it does not have sense
# nest loops or put anything after the loop 
# loop tag also acts as a synced block
loop:

@channel3 <- a1 b2 c3 c4 c5
```

## Grammar

```
program := constantDef* statement* ('loop:' statement*)?

constantDef := ID ':=' (instrumentDef | tablature+)

instrumentDef := '{' mapEntry* ('adsr:' adsrVector)? mapEntry* '}'

tablature := (ID | NOTE | SILENCE | OCTAVE | INCOCT | DECOCT | tuplet | '|')+
tuplet := '{' tablature '}' NUM

ID := $(\w)+

statement := channelFill | SYNC 
channelFill := CHANNEL_ID '<-' tablature+
SYNC := '-'*

```

## Binary compilation

```
# music start byte
# block for instruments
# block for 
# block for music
    # address of loop
    # other blocks for music
```
