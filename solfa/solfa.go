package solfa

import (
	"fmt"
	"regexp"
	"strconv"
)

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
	Length   uint // as a divisor 1: whole note
	Tuplet   uint // e.g. 3 means this note is part of a triplet
	Halftone Halftone
	Octave   uint
}

type channel struct {
	Octave uint
	Tempo  uint
}

const (
	minOctave     = 1
	maxOctave     = 8
	minLength     = 1
	maxLength     = 64
	defaultLength = 4
)

var regex = regexp.MustCompile(`^\{[^\{}]*}\d+`)

var tokenizer = regexp.MustCompile(`^(` +
	`(([a-zA-Z])([+\-#]?)(\d*))` + // note, octave, silence...
	`|([<>])` + // increase/decrease octave
	`|\s+|\|` + // spaces or vertical bars are ignored
	`|(\{[^\{}]*}\d+)` + // anything into brackets followed by a number: tuplet
	`)`)

func Parse(tab []byte) ([]Note, error) {
	global := channel{
		Octave: 4,
		Tempo:  120,
	}
	return parseSubstring(tab, &global)
}

func parseSubstring(tab []byte, global *channel) ([]Note, error) {

	var notes []Note

	index := 0
	for len(tab) > 0 {
		token := tokenizer.Find(tab)
		if token == nil {
			return nil, fmt.Errorf("at position %d: unexpected charecter '%c'", index, tab[0])
		}
		tab = tab[len(token):]
		switch token[0] {
		case 'o', 'O':
			if err := parseOctave(token, global); err != nil {
				return nil, fmt.Errorf("at position %d: %w", index, err)
			}
		case 'r', 'R':
			note, err := parseSilence(token)
			if err != nil {
				return nil, fmt.Errorf("at position %d: %w", index, err)
			}
			notes = append(notes, note)
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g',
			'A', 'B', 'C', 'D', 'E', 'F', 'G':
			note, err := parseNote(token, global)
			if err != nil {
				return nil, fmt.Errorf("at position %d: %w", index, err)
			}
			notes = append(notes, note)
		case '<':
			if global.Octave == minOctave {
				return nil, fmt.Errorf("at position %d: can't set octave lower than %d", index, minOctave)
			}
			global.Octave--
		case '>':
			if global.Octave == maxOctave {
				return nil, fmt.Errorf("at position %d: can't set an octave greater than %d", index, maxOctave)
			}
			global.Octave++
		case '{':
			tNotes, err := parseTuplet(token, global)
			if err != nil {
				return nil, fmt.Errorf("tuplet starting at %d: %w", index, err)
			}
			notes = append(notes, tNotes...)
		case '|':
			// just ignore
		default:
			return nil, fmt.Errorf("unknown char at position %d: %c", index, token)
		}
		index += len(token)
	}
	return notes, nil
}

var tuplet = regexp.MustCompile(`^\{(.*)}(\d+)$`)

func parseTuplet(token []byte, c *channel) ([]Note, error) {
	sm := tuplet.FindSubmatch(token)
	if sm == nil {
		panic(fmt.Sprintf("wrong format for tuplet: %q! this is a bug", string(token)))
	}

	tNotes, err := parseSubstring(sm[1], c)
	if err != nil {
		return nil, err
	}

	nTuple, err := strconv.Atoi(string(sm[2]))
	if err != nil {
		return nil, fmt.Errorf("tuple needs a proper number suffix: %w", err)
	}

	for i := range tNotes {
		tNotes[i].Tuplet = uint(nTuple)
	}

	return tNotes, nil
}

func parseSilence(token []byte) (Note, error) {
	n := Note{Pitch: Silence}
	if len(token) == 1 {
		n.Length = defaultLength
		return n, nil
	}
	length, err := strconv.Atoi(string(token[1:]))
	if err != nil {
		return n, fmt.Errorf("wrong format for silence: %q. It must be an 'R' followed by a number", string(token))
	}
	n.Length = uint(length)
	return n, nil
}

var note = regexp.MustCompile(`^([+\-#]?)(\d*)$`)

// todo: return error if a given note can't be sharp or flat
func parseNote(token []byte, c *channel) (Note, error) {
	sm := note.FindSubmatch(token[1:])
	if sm == nil {
		panic(fmt.Sprintf("wrong format for note: %q! this is a bug", string(token)))
	}
	n := Note{
		Pitch:    getPitch(token[0]),
		Length:   defaultLength,
		Octave:   c.Octave,
		Halftone: NoHalftone,
	}
	// read Length
	if len(sm[2]) > 0 {
		l, err := strconv.Atoi(string(sm[2]))
		if err != nil {
			panic(fmt.Sprintf("wrong length for note: %q! this is a bug. Err: %s",
				string(token), err.Error()))
		}
		if l < minLength || l > maxLength {
			return Note{}, fmt.Errorf(
				"wrong note length: %d. Must be in range %d to %d", l, minLength, maxLength)
		}
		n.Length = uint(l)
	}
	// read halftone
	if len(sm[1]) > 0 {
		switch sm[1][0] {
		case '+', '#':
			n.Halftone = Sharp
		case '-':
			n.Halftone = Flat
		default:
			panic(fmt.Sprintf("wrong halftone '%c'! this is a bug", sm[1][0]))
		}
	}
	return n, nil
}

func parseOctave(token []byte, c *channel) error {
	i, err := strconv.Atoi(string(token[1:]))
	if err != nil {
		return err
	}
	if i < minOctave || i > maxOctave {
		return fmt.Errorf("octave value must be 1 to 8")
	}
	c.Octave = uint(i)
	return nil
}

var pitches = [8]Pitch{A, B, C, D, E, F, G}

func getPitch(c byte) Pitch {
	if c >= 'A' && c <= 'Z' {
		return pitches[c-'A']
	}
	if c >= 'a' && c <= 'z' {
		return pitches[c-'a']
	}
	panic(fmt.Sprintf("pitch can't be '%c'! this is a bug", c))
}
