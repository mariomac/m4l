package solfa

type Length float64

// USA name for note duration
const (
	Whole        Length = 1
	Half         Length = 2
	Quarter      Length = 4
	Eigth        Length = 8
	Sixteenth    Length = 16
	ThirtySecond Length = 32
	SixtyFourth  Length = 64
)

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

type Accidental uint8

const (
	None  Accidental = 0
	Sharp Accidental = '#' //increases pitch by one semitone
	Flat  Accidental = 'b' // lowers pitch by one semitone
	// Todo: consider others http://neilhawes.com/sstheory/theory17.htm
)

type Note struct {
	Pitch      Pitch
	Length     Length
	Accidental Accidental
	Octave     uint8
}

type channel struct {
	LastNote Note
	Octave uint8
	Tempo int
}

// states of a turing machine
type parseState int
const (
	initial parseState = iota
	setPitch
	setAccident
	setLocalOctave
	globalOctave
	setGlobalOctave
)
const (
	defaultLength = Quarter
)
func ParseChannel(tab []byte) ([]Note, error) {
	status := initial
	global := channel{
		Octave: 4,
		Tempo: 120,
	}
	var notes []Note
	for _, c := range tab {
		switch status {
		case initial:
			if p, ok := isPitch(c) ; ok {
				global.LastNote = Note{
					Pitch: p,
					Length: defaultLength,
					Octave: global.Octave,
					Accidental: None,
				}
			}
		}
	}
}

var pitches = [8]Pitch{A, B, C, D, E, F, G}
func isPitch(c byte) (Pitch, bool) {
	if c >= 'A' && c <= 'Z' {
		return pitches[c - 'A'], true
	}
	if c >= 'a' && c <= 'z' {
		return pitches[c - 'a'], true
	}
	return 0, false
}
