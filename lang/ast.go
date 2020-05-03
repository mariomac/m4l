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

// channel: CHANNELID <- tablature
func (p *Parser) channelNode(s *song.Song) error {
	eofErr := func() error {
		return fmt.Errorf("unexpected EOF at position %d:%d", p.t.row, p.t.col)
	}
	last := p.t.Get()
	chName := last.Submatch[0]
	c, ok := s.Channels[chName]
	if !ok {
		c = &song.Channel{Name: chName, Instrument: song.DefaultInstrument}
		c.Status.Octave = defaultOctave
		s.Channels[chName] = c
	}
	if !p.t.Next() {
		return eofErr()
	}
	last = p.t.Get()
	if last.Type != ChannelSendArrow {
		return SyntaxError{t: last}
	}
	if !p.t.Next() {
		return eofErr()
	}
	if err := p.tablatureNode(c); err != nil {
		return err
	}
	return nil
}

type SyntaxError struct {
	t Token
}

func (p SyntaxError) Error() string {
	return fmt.Sprintf("%d:%d - Unexpected %q", p.t.Row, p.t.Col, p.t.Content)
}
