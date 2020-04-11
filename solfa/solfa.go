package solfa

import "fmt"

type Pitch uint8

const (
	A Pitch = 'a'
	B Pitch = 'b'
	C Pitch = 'c'
	D Pitch = 'd'
	E Pitch = 'e'
	F Pitch = 'f'
	G Pitch = 'g'
)

type Halftone uint8

const (
	None  Halftone = 0
	Sharp Halftone = '#' //increases pitch by one semitone
	Flat  Halftone = 'b' // lowers pitch by one semitone
	// Todo: consider others http://neilhawes.com/sstheory/theory17.htm
)

type Note struct {
	Pitch    Pitch
	Length   uint8  // as a divisor 1: whole note
	Halftone Halftone
	Octave   uint8
}

type channel struct {
	LastNote *Note
	Octave   int
	Tempo    int
}

// states of a turing machine
type parseState int

const (
	initial parseState = iota
	setPitch
	setHalftone
	setLength
	globalOctave
	setGlobalOctave
)
const (
	defaultLength = 4
)

func ParseChannel(tab []byte) ([]Note, error) {
	status := initial
	global := channel{
		Octave: 4,
		Tempo:  120,
	}
	var notes []Note

	var addNote = func(p Pitch) {
		n := Note{
			Pitch:    p,
			Length:   defaultLength,
			Octave:   uint8(global.Octave),
			Halftone: None,
		}
		notes = append(notes, n)
		global.LastNote = &n
	}

	for i, c := range tab {
		switch status {
		case initial:
			if pitch, ok := isPitch(c); ok {
				addNote(pitch)
				status = setPitch
			} else if isOctave(c) {
				status = globalOctave
			} else if inc, ok := isIncOctave(c); ok {
				global.Octave += inc
				status = initial
			} else {
				return notes, fmt.Errorf("unexpected character %c at position %d", c, i)
			}
		case setPitch:
			if pitch, ok := isPitch(c); ok {
				addNote(pitch)
				status = setPitch
			} else if isOctave(c) {
				status = globalOctave
			} else if inc, ok := isIncOctave(c); ok {
				global.Octave += inc
				status = initial
			} else if acc, ok := isHalftone(c); ok {
				global.LastNote.Halftone = acc
				status = setHalftone
			} else if d, ok := isDigit(c); ok {
				global.LastNote.Length = uint8(d)
				status = setLength
			} else {
				return notes, fmt.Errorf("unexpected character %c at position %d", c, i)
			}
		case setHalftone:
			if pitch, ok := isPitch(c); ok {
				addNote(pitch)
				status = setPitch
			} else if isOctave(c) {
				status = globalOctave
			} else if inc, ok := isIncOctave(c); ok {
				global.Octave += inc
				status = initial
			} else {
				return notes, fmt.Errorf("unexpected character %c at position %d", c, i)
			}
		case setLength:
			if pitch, ok := isPitch(c); ok {
				addNote(pitch)
				status = setPitch
			} else if isOctave(c) {
				status = globalOctave
			} else if inc, ok := isIncOctave(c); ok {
				global.Octave += inc
				status = initial
			} else
		}
	}
}

var pitches = [8]Pitch{A, B, C, D, E, F, G}

func isPitch(c byte) (Pitch, bool) {
	if c >= 'A' && c <= 'Z' {
		return pitches[c-'A'], true
	}
	if c >= 'a' && c <= 'z' {
		return pitches[c-'a'], true
	}
	return 0, false
}

func isHalftone(c byte) (Halftone, bool) {
	if c == '#' || c == '+' {
		return Sharp, true
	}
	if c == '-' {
		return Flat, true
	}
	return None, false
}

func isOctave(c byte) bool {
	return c == 'o' || c == 'O'
}

func isDigit(c byte) (int, bool) {
	if c >= '0' && c <= '9' {
		return int(c - '0'), true
	}
	return -1, false
}

func isIncOctave(c byte) (int, bool) {
	if c == '<' {
		return -1, true
	}
	if c == '>' {
		return 1, true
	}
	return 0, false
}
