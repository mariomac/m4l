package lang

import (
	"fmt"
	"strconv"

	"github.com/mariomac/msxmml/pkg/song"
	"github.com/mariomac/msxmml/pkg/song/note"
)

// tablature := (ID | NOTE | SILENCE | OCTAVE | INCOCT | DECOCT | tuplet | '|')+
func (p *Parser) tablatureNode() (song.Tablature, error) {
	t := song.Tablature{}
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case Note:
			if n, err := parseNote(tok); err != nil {
				return nil, err
			} else {
				t = append(t, song.TablatureItem{Note: &n})
			}
		case Silence:
			n := parseSilence(tok)
			t = append(t, song.TablatureItem{Note: &n})
		case Octave:
			o := parseOctave(tok)
			t = append(t, song.TablatureItem{SetOctave: &o})
		case IncOctave, DecOctave:
			o := parseOctaveStep(tok)
			t = append(t, song.TablatureItem{OctaveStep: &o})
		case OpenKey:
			if tu, err := p.tupletNode(); err != nil {
				return nil, err
			} else {
				t = append(t, tu...)
			}
		case Separator:
		// just ignore
		default:
			// end of tablature, return
			return t, nil
		}
		p.t.Next()
	}
	return t, nil
}

// A note should come represented by an array where
// 0: pitch - 1: halftone - 2: length - 3: dots
// todo: return t, error if a given note can't be sharp or flat
func parseNote(token Token) (note.Note, error) {
	n := note.Note{
		Pitch:    getPitch(token.Submatch[0][0]),
		Length:   defaultLength,
		Halftone: note.NoHalftone,
		Dots:     len(token.Submatch[3]),
	}
	// get halftone
	if len(token.Submatch[1]) > 0 {
		switch token.Submatch[1][0] {
		case '+', '#':
			n.Halftone = note.Sharp
		case '-':
			n.Halftone = note.Flat
		default:
			panic(fmt.Sprintf("wrong halftone %q! this is a bug", token.Submatch[1]))
		}
	}

	// get Length
	if len(token.Submatch[2]) > 0 {
		l, err := strconv.Atoi(token.Submatch[2])
		if err != nil {
			panic(fmt.Sprintf("wrong length for note: %#v! this is a bug. Err: %s",
				token, err.Error()))
		}
		if l < minLength || l > maxLength {
			return n, fmt.Errorf(
				"wrong note length: %d. Must be in range %d to %d", l, minLength, maxLength)
		}
		n.Length = l
	}
	return n, nil
}

var pitches = [8]note.Pitch{note.A, note.B, note.C, note.D, note.E, note.F, note.G}

func getPitch(c byte) note.Pitch {
	if c >= 'A' && c <= 'Z' {
		return pitches[c-'A']
	}
	if c >= 'a' && c <= 'z' {
		return pitches[c-'a']
	}
	panic(fmt.Sprintf("pitch can't be '%c'! this is a bug", c))
}

func parseSilence(token Token) note.Note {
	n := note.Note{Pitch: note.Silence}
	if len(token.Submatch[0]) == 0 {
		n.Length = defaultLength
		return n
	}
	length, err := strconv.Atoi(token.Submatch[0])
	if err != nil {
		panic(fmt.Sprintf("silence can't be %q! this is a bug", token.Submatch))
	}
	n.Length = length
	return n
}

func parseOctave(token Token) int {
	oct, err := strconv.Atoi(token.Submatch[0])
	if err != nil {
		panic(fmt.Sprintf("octave can't be %q! this is a bug", token.Submatch))
	}
	return oct
}

func parseOctaveStep(token Token) int {
	switch token.Content[0] {
	case '<':
		return -1
	case '>':
		return +1
	default:
		panic(fmt.Sprintf("invalid octave step %q! This is a bug", token.Content))
	}
}

func (p *Parser) tupletNode() (song.Tablature, error) {
	return nil, nil
}

/*
func parseApplyTuplet(token Token, c *song.Channel, startTupletIndex *int) error {
	if *startTupletIndex < 0 {
		return ParserError{token, "closing a non-opened tuple"}
	}
	nTuple, err := strconv.Atoi(token.Submatch[0])
	if err != nil {
		panic(fmt.Sprintf("invalid tuple number %q! This is a bug", token.Submatch[0]))
	}
	if nTuple < 3 {
		return ParserError{token, "tuplet number should be at least 3"}
	}
	for i := *startTupletIndex; i < len(c.Notes); i++ {
		c.Notes[i].Tuplet = nTuple
	}
	return nil
}
 */
