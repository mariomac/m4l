package lang

import (
	"fmt"
	"strings"

	"github.com/mariomac/msxmml/pkg/song"
)

const (
	defaultOctave = 4
	minOctave     = 0
	maxOctave     = 8
	minLength     = 1
	maxLength     = 64
	defaultLength = 4
)

func (p *Parser) eofErr() error {
	return UnexpecedEofError{Row: p.t.row, Col: p.t.col}
}

type Parser struct {
	t *Tokenizer
}

// Convention: tokenizer always receives a tokenizer with a token available, excepting the Root
// program := constantDef* statement* ('loop:' statement*)?
func Parse(t *Tokenizer) (*song.Song, error) {
	p := &Parser{
		t: t,
	}
	s := &song.Song{
		Constants: map[string]song.Tablature{},
		LoopIndex: -1,
	}
	s.AddSyncedBlock()

	p.t.Next()
	for !p.t.EOF() {
		token := p.t.Get()
		switch token.Type {
		case ConstDef:
			if err := p.constantDefNode(s); err != nil {
				return nil, err
			}
		case LoopTag:
			if err := p.loopNode(s); err != nil {
				return nil, err
			}
		case ChannelSync:
			s.AddSyncedBlock()
			p.t.Next()
		case ChannelId:
			if err := p.channelFillNode(s); err != nil {
				return nil, err
			}
		default:
			return nil, SyntaxError{t: token}
		}
	}
	return s, nil
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
	case OpenKey:
		inst, err := p.instrumentDefinitionNode()
		if err != nil {
			return err
		}
		s.Constants[id] = song.Tablature{{Instrument: &inst}}
	default:
		tabl, err := p.tablatureNode(s)
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

// instrumentDef := '{' mapEntry* ('adsr:' adsrVector)? mapEntry* '}'
func (p *Parser) instrumentDefinitionNode() (song.Instrument, error) {
	inst := song.Instrument{}
	if !p.t.Next() {
		return inst, p.eofErr()
	}
	definedAdsr, definedWave := false, false
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case AdsrVector:
			if definedAdsr {
				return inst, ParserError{tok, "defining ADSR envelope twice"}
			}
			definedAdsr = true
			inst.Envelope = tok.getAdsr()
		case MapEntry:
			switch strings.ToLower(tok.getMapKey()) {
			case "adsr":
				return inst, ParserError{tok, "adsr should have a value like: 20->100, 50->80, 100, 120"}
			case "wave":
				if definedWave {
					return inst, ParserError{tok, "wave is defined twice"}
				}
				definedWave = true
				// todo: maybe validate wave values?
				inst.Wave = tok.getWave()
			default:
				return inst, ParserError{tok, "only 'adsr' and 'wave' properties are allowed"}
			}
		case CloseKey:
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
func (p *Parser) tablatureNode(s *song.Song) (song.Tablature, error) {
	t := song.Tablature{}
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case ConstRef:
			// todo: test that unexisting const id returns an error
			id := tok.getConstRefId()
			if _, ok := s.Constants[id]; !ok {
				return nil, ParserError{t: tok, msg: fmt.Sprintf("constant %q not defined", id)}
			}
			t = append(t, song.TablatureItem{ConstantRef: &id})
		case Note:
			if n, err := tok.getNote(); err != nil {
				return nil, err
			} else {
				t = append(t, song.TablatureItem{Note: &n})
			}
		case Silence:
			n := tok.getSilence()
			t = append(t, song.TablatureItem{Note: &n})
		case Octave:
			o := tok.getOctave()
			t = append(t, song.TablatureItem{SetOctave: &o})
		case OctaveStep:
			o := tok.getOctaveStep()
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

// tuplet := '{' (NOTE|OCTAVE|INCOCT|DECOCT) + '}' NUM
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
		case AnyString:
			return nil, SyntaxError{tok}
		default:
			return nil, p.eofErr()
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
	tab, err := p.tablatureNode(s)
	if err != nil {
		return err
	}
	// tablature might be empty. Return error or just accept it?
	s.AddItems(channelId, tab...)

	// not advancing the tokenizer. After a tablature, the token points to the next statement
	return nil
}
