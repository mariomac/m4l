package lang

import (
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
	s := &song.Song{}
	s.AddSyncedBlock()

	p.t.Next()
	for !p.t.EOF() {
		token := p.t.Get()
		switch token.Type {
		case ConstName:
			if err := p.constantDefNode(s); err != nil {
				return nil, err
			}
		case LoopTag:
			if err := p.loopNode(s); err != nil {
				return nil, err
			}
		case ChannelId, ChannelSync:
			if err := p.statementNode(s); err != nil {
				return nil, err
			}
		default:
			return nil, SyntaxError{t: token}
		}
		p.t.Next()
	}
	return s, nil
}

// constantDef := ID ':=' (instrumentDef | tablature+)
func (p *Parser) constantDefNode(s *song.Song) error {
	tok := p.t.Get()
	id := tok.getConstID()
	if _, ok := s.Constants[id]; ok {
		return RedefinitionError{tok}
	}
	if !p.t.Next() {
		return p.eofErr()
	}
	tok = p.t.Get()
	if tok.Type != Assign {
		return SyntaxError{tok}
	}
	tok = p.t.Get()
	switch tok.Type {
	case OpenKey:
		inst, err := p.instrumentDefinitionNode(s)
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
	return nil
}

// ('loop:' statement*)
func (p *Parser) loopNode(s *song.Song) error {
	return nil
}

// statement := channelFill | SYNC
func (p *Parser) statementNode(s *song.Song) error {
	return nil
}

func (p *Parser) instrumentDefinitionNode(s *song.Song) (song.Instrument, error) {
	return nil
}

func (p *Parser) tablatureNode(s *song.Song) (song.Tablature, error) {
	return nil
}
