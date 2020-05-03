package lang

import (
	"fmt"

	"github.com/mariomac/msxmml/song"
)

const (
	defaultOctave = 4
	minOctave     = 0
	maxOctave     = 8
	minLength     = 1
	maxLength     = 64
	defaultLength = 4
)

type SyntaxError struct {
	t Token
}

func (p SyntaxError) Error() string {
	return fmt.Sprintf("%d:%d - Unexpected %q", p.t.Row, p.t.Col, p.t.Content)
}

func (p *Parser) eofErr() error {
	return fmt.Errorf("unexpected EOF at position %d:%d", p.t.row, p.t.col)
}

type Parser struct {
	t *Tokenizer
}

// Convention: tokenizer always receives a tokenizer with a token available, excepting the Root

// song: channel+
func Parse(t *Tokenizer) (*song.Song, error) {
	p := &Parser{
		t: t,
	}
	s := &song.Song{Channels: map[string]*song.Channel{}}
	p.t.Next()
	for !p.t.EOF() {
		token := p.t.Get()
		switch token.Type {
		case ChannelID:
			if err := p.channelNode(s); err != nil {
				return nil, err
			}
		default:
			return nil, SyntaxError{t: token}
		}
	}
	return s, nil
}

// channel: CHANNELID ( '<-' tablature | '{' instumentDefinition '}' )
func (p *Parser) channelNode(s *song.Song) error {
	last := p.t.Get()
	chName := last.Submatch[0]
	c, ok := s.Channels[chName]
	if !ok {
		c = &song.Channel{Name: chName, Instrument: song.DefaultInstrument}
		c.Status.Octave = defaultOctave
		s.Channels[chName] = c
	}
	if !p.t.Next() {
		return p.eofErr()
	}
	last = p.t.Get()
	switch last.Type {
	case ChannelSendArrow:
		if err := p.tablatureNode(c); err != nil {
			return err
		}
	case OpenSection:
		if err := p.instrumentDefinitionNode(c); err != nil {
			return err
		}
	default:
		return SyntaxError{t: last}
	}
	return nil
}

type ParserError struct {
	t   Token
	msg string
}

func (t ParserError) Error() string {
	return fmt.Sprintf("At %d,%d: %s", t.t.Row, t.t.Col, t.msg)
}
