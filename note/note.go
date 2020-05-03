package note

type Pitch uint8

const (
	Silence Pitch = 0
	A       Pitch = 'a'
	B       Pitch = 'b'
	C       Pitch = 'c'
	D       Pitch = 'd'
	E       Pitch = 'e'
	F       Pitch = 'f'
	G       Pitch = 'g'
)

type Halftone uint8

const (
	NoHalftone Halftone = 0
	Sharp      Halftone = '#' //increases pitch by one semitone
	Flat       Halftone = '-' // lowers pitch by one semitone
	// Todo: consider others http://neilhawes.com/sstheory/theory17.htm
)

type Note struct {
	Pitch    Pitch
	Length   int // as a divisor 1: whole note
	Tuplet   int // e.g. 3 means this note is part of a triplet
	Halftone Halftone
	Octave   int
	Dots     int // number of dots
}
