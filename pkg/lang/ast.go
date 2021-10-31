package lang

import (
	"fmt"
	"io"

	"github.com/mariomac/msxmml/pkg/song"
)

const (
	defaultOctave = 4
	minOctave     = 0
	maxOctave     = 8
	minLength     = 1
	maxLength     = 64
	defaultLength = 4
	maxVolume     = 15
)

func (p *Parser) eofErr() error {
	return UnexpecedEofError{Row: p.t.row, Col: p.t.col}
}

type Parser struct {
	t *Tokenizer
}

// Convention: tokenizer always receives a tokenizer with a token available, excepting the Root
// program := constantDef* statement* ('loop:' statement*)?
func Parse(reader io.ReadSeeker) (*song.Song, error) {
	props, lines, err := parseHeader(reader)
	if err != nil {
		return nil, err
	}

	t := NewTokenizer(reader, lines)
	p := &Parser{
		t: t,
	}
	s := &song.Song{
		Properties: props,
		Constants:  map[string]song.Tablature{},
		LoopIndex:  -1,
	}

	if err := parseBody(s, p); err != nil {
		return nil, err
	}
	return s, nil
}

func parseBody(s *song.Song, p *Parser) error {
	s.AddSyncedBlock()
	p.t.Next()
	for !p.t.EOF() {
		token := p.t.Get()
		switch token.Type {
		case ConstDef:
			if err := p.constantDefNode(s); err != nil {
				return err
			}
		case LoopTag:
			if err := p.loopNode(s); err != nil {
				return err
			}
		case ChannelSync:
			s.AddSyncedBlock()
			p.t.Next()
		case ChannelId:
			if err := p.channelFillNode(s); err != nil {
				return err
			}
		default:
			return SyntaxError{t: token}
		}
	}
	return nil
}

// constantDef := ID ':=' (instrumentDef | tablature+)
func (p *Parser) constantDefNode(s *song.Song) error {
	tok := p.t.Get()
	id := tok.getConstDefId()
	if _, ok := s.Constants[id]; ok {
		return RedefinitionError{tok}
	}
	if !p.t.Next() {
		return p.eofErr()
	}
	tok = p.t.Get()
	switch tok.Type {
	case OpenInstrument:
		inst, err := p.instrumentDefinitionNode(tok.getInstrumentClass())
		if err != nil {
			return err
		}
		s.Constants[id] = song.Tablature{{Instrument: &inst}}
	default:
		tabl, err := p.tablatureNode(s, false)
		if err != nil {
			return err
		}
		s.Constants[id] = tabl
	}
	// not running p.t.Next as it was the last statement in both instrument and tablature nodes
	return nil
}

func (p *Parser) loopNode(s *song.Song) error {
	if s.LoopIndex >= 0 {
		return ParserError{t: p.t.Get(), msg: "duplicate 'loop:' tag"}
	}
	s.LoopIndex = len(s.Blocks)
	s.AddSyncedBlock()
	p.t.Next()
	return nil
}

// instrumentDef := class '{' mapEntry* ('adsr:' adsrVector)? mapEntry* '}'
func (p *Parser) instrumentDefinitionNode(class string) (song.Instrument, error) {
	inst := song.Instrument{
		Class:      class,
		Properties: map[string]string{},
	}
	if !p.t.Next() {
		return inst, p.eofErr()
	}
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case MapEntry:
			k, v := tok.getMapKeyValue()
			inst.Properties[k] = v
		case CloseInstrument:
			p.t.Next()
			return inst, nil
		default:
			return inst, SyntaxError{tok}
		}
		p.t.Next()
	}
	return inst, nil
}

// tablature := (ID | NOTE | SILENCE | OCTAVE | INCOCT | DECOCT | tuplet | '|')+
func (p *Parser) tablatureNode(s *song.Song, allowConstants bool) (song.Tablature, error) {
	t := song.Tablature{}
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case ConstRef:
			if !allowConstants {
				return nil, ParserError{t: tok, msg: "can't refer constants inside constants"}
			}
			id := tok.getConstRefId()
			if items, ok := s.Constants[id]; !ok {
				return nil, ParserError{t: tok, msg: fmt.Sprintf("constant %q not defined", id)}
			} else { // expand constant as notes
				t = append(t, items...)
			}
		case Note:
			if n, err := tok.getNote(); err != nil {
				return nil, err
			} else {
				t = append(t, song.TablatureItem{Note: &n})
			}
		case Volume:
			n := tok.getVolume()
			if n > maxVolume {
				return t, ParserError{t: tok, msg: fmt.Sprintf("max volume is 16 (was: %d)", n)}
			}
			t = append(t, song.TablatureItem{Volume: &n})
		case Silence:
			n := tok.getSilence()
			t = append(t, song.TablatureItem{Note: &n})
		case Octave:
			o := tok.getOctave()
			t = append(t, song.TablatureItem{SetOctave: &o})
		case OctaveStep:
			o := tok.getOctaveStep()
			t = append(t, song.TablatureItem{OctaveStep: &o})
		case OpenTuple:
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

// tuplet := '(' (NOTE|OCTAVE|INCOCT|DECOCT) + ')' NUM
func (p *Parser) tupletNode() (song.Tablature, error) {
	if !p.t.Next() {
		return nil, p.eofErr()
	}
	t := song.Tablature{}
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case Note:
			if n, err := tok.getNote(); err != nil {
				return nil, err
			} else {
				t = append(t, song.TablatureItem{Note: &n})
			}
		case Volume:
			n := tok.getVolume()
			if n > maxVolume {
				return t, ParserError{t: tok, msg: fmt.Sprintf("max volume is 16 (was: %d)", n)}
			}
			t = append(t, song.TablatureItem{Volume: &n})
		case Silence:
			n := tok.getSilence()
			t = append(t, song.TablatureItem{Note: &n})
		case Octave:
			o := tok.getOctave()
			t = append(t, song.TablatureItem{SetOctave: &o})
		case OctaveStep:
			o := tok.getOctaveStep()
			t = append(t, song.TablatureItem{OctaveStep: &o})
		case CloseTuple:
			tn := tok.getTupletNumber()
			for n := range t {
				if t[n].Note != nil {
					t[n].Note.Tuplet = tn
				}
			}
			p.t.Next()
			return t, nil
		case Separator:
		// just ignore
		case NoMatch:
			return nil, SyntaxError{tok}
		default:
			if p.t.EOF() {
				return nil, p.eofErr()
			}
			return nil, SyntaxError{tok}
		}
		p.t.Next()
	}
	return nil, p.eofErr()
}

func (p *Parser) channelFillNode(s *song.Song) error {
	tok := p.t.Get()
	channelId := tok.getChannelId()
	if !p.t.Next() {
		return p.eofErr()
	}
	tok = p.t.Get()
	if tok.Type != SendArrow {
		return SyntaxError{t: tok}
	}
	if !p.t.Next() {
		return p.eofErr()
	}
	tab, err := p.tablatureNode(s, true)
	if err != nil {
		return err
	}
	// tablature might be empty. Return error or just accept it?
	s.AddItems(channelId, tab...)

	// not advancing the tokenizer. After a tablature, the token points to the next statement
	return nil
}
