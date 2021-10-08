package lang

import (
	"fmt"
	"github.com/mariomac/msxmml/pkg/lang"
	"github.com/mariomac/msxmml/pkg/song"
	"github.com/mariomac/msxmml/pkg/song/note"
	"strconv"
)

// tablature: (octave|note|pause|....)+
func (p *lang.Parser) tablatureNode(c *song.Channel) error {
	if !p.t.Next() {
		return p.eofErr()
	}
	startTupletIndex := -1
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case lang.Note:
			if n, err := parseNote(tok, c); err != nil {
				return err
			} else {
				c.Notes = append(c.Notes, n)
			}
		case lang.Silence:
			if n, err := parseSilence(tok); err != nil {
				return err
			} else {
				c.Notes = append(c.Notes, n)
			}
		case lang.Octave:
			if err := parseOctave(tok, c); err != nil {
				return err
			}
		case lang.IncOctave, lang.DecOctave:
			if err := parseOctaveStep(tok, c); err != nil {
				return err
			}
		case OpenSection:
			if startTupletIndex >= 0 {
				return lang.ParserError{tok, "can't open a tuple inside a tuple"}
			}
			startTupletIndex = len(c.Notes)
		case CloseTuplet:
			if err := parseApplyTuplet(tok, c, &startTupletIndex) ; err != nil {
				return err
			}
			startTupletIndex = -1
		case lang.Separator:
		// just ignore
		default:
			// end of tablature, return
			return nil
		}
		p.t.Next()
	}
	return nil
}

// A note should come represented by an array where
// 0: pitch - 1: halftone - 2: length - 3: dots
// todo: return error if a given note can't be sharp or flat
func parseNote(token lang.Token, c *song.Channel) (note.Note, error) {
	n := note.Note{
		Pitch:    getPitch(token.Submatch[0][0]),
		Length:   lang.defaultLength,
		Octave:   c.Status.Octave,
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
		if l < lang.minLength || l > lang.maxLength {
			return n, fmt.Errorf(
				"wrong note length: %d. Must be in range %d to %d", l, lang.minLength, lang.maxLength)
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

func parseSilence(token lang.Token) (note.Note, error) {
	n := note.Note{Pitch: note.Silence}
	if len(token.Submatch[0]) == 0 {
		n.Length = lang.defaultLength
		return n, nil
	}
	length, err := strconv.Atoi(token.Submatch[0])
	if err != nil {
		panic(fmt.Sprintf("silence can't be %q! this is a bug", token.Submatch))
	}
	n.Length = length
	return n, nil
}

func parseOctave(token lang.Token, c *song.Channel) error {
	oct, err := strconv.Atoi(token.Submatch[0])
	if err != nil {
		panic(fmt.Sprintf("silence can't be %q! this is a bug", token.Submatch))
	}
	if err := assertOctave(oct); err != nil {
		return lang.ParserError{token, err.Error()}
	}
	c.Status.Octave = oct
	return nil
}

func parseOctaveStep(token lang.Token, c *song.Channel) error {
	oct := c.Status.Octave
	switch token.Content[0] {
	case '<':
		oct--
	case '>':
		oct++
	default:
		panic(fmt.Sprintf("invalid octave step %q! This is a bug", token.Content))
	}
	if err := assertOctave(oct); err != nil {
		return lang.ParserError{token, err.Error()}
	}
	c.Status.Octave = oct
	return nil
}

func assertOctave(oct int) error {
	if oct < lang.minOctave || oct > lang.maxOctave {
		return fmt.Errorf("octave must be in range [%d..%d]. Actual: %d", lang.minOctave, lang.maxOctave, oct)
	}
	return nil
}

func parseApplyTuplet(token lang.Token, c *song.Channel, startTupletIndex *int) error {
	if *startTupletIndex < 0 {
		return lang.ParserError{token, "closing a non-opened tuple"}
	}
	nTuple, err := strconv.Atoi(token.Submatch[0])
	if err != nil {
		panic(fmt.Sprintf("invalid tuple number %q! This is a bug", token.Submatch[0]))
	}
	if nTuple < 3 {
		return lang.ParserError{token, "tuplet number should be at least 3"}
	}
	for i := *startTupletIndex ; i < len(c.Notes) ; i++ {
		c.Notes[i].Tuplet = nTuple
	}
	return nil
}

