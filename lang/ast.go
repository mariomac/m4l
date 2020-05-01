package lang

import (
	"fmt"
	"regexp"

	"github.com/mariomac/msxmml/solfa"

	"github.com/mariomac/msxmml/song"
)

// Convention: tokenizer always receives a tokenizer with a token available, excepting the Root

func Root(t *Tokenizer) (*song.Song, error) {
	s := &song.Song{Channels: map[string]song.Channel{}}
	t.Next()
	for !t.EOF() {
		token := t.Get()
		switch token.Type {
		case Channel:
			ch, err := ChannelNode(t)
			if err != nil {
				return nil, err
			}
			s.Channels[ch.Name] = ch
		default:
			return nil, &ParserError{t, fmt.Errorf("unexpected input: %q. I was expecting a channel ID", string(token.Content))}
		}
	}

	return s, nil
}

func ChannelNode(t *Tokenizer) (song.Channel, error) {
	last := t.Get()
	c := song.Channel{Name: string(last.Content[1:])}

	if !t.Next() {
		return c, unexpectedError(t, "channel information", []byte("end of input"))
	}
	last = t.Get()
	if last.Type != ChannelSendArrow {
		return c, unexpectedError(t, "an arrow '<-'", last.Content)
	}
	if !t.Next() {
		return c, unexpectedError(t, "channel information", []byte("end of input"))
	}
	tabs, err := TablatureNode(t)
	if err != nil {
		return c, err
	}
	// TODO: keep last note/global config for each channel number
	// so octave and other data stays between successive commands
	c.Notes, err = solfa.Parse(tabs)
	if err != nil {
		return c, &ParserError{t: t, cause: fmt.Errorf("problem with channel tablature: %w", err)}
	}
	return c, nil
}

var tabRegex = regexp.MustCompile(`^(([a-zA-Z][+\-#]?\d*)|[<>]|\||(\{[^\{}]*}\d+))*$`)

func TablatureNode(t *Tokenizer) ([]byte, error) {
	var tablature []byte
	tok := t.Get()
	if tok.Type != String && !tabRegex.Match(tok.Content) {
		return nil, unexpectedError(t, "a music tablature", tok.Content)
	}
	tablature = append(tablature, tok.Content...)
	for t.Next() {
		tok := t.Get()
		if tok.Type != String && !tabRegex.Match(tok.Content) {
			return tablature, nil
		}
		tablature = append(tablature, tok.Content...)
	}

	return tablature, nil
}

type ParserError struct {
	t     *Tokenizer
	cause error
}

func (p *ParserError) Error() string {
	return fmt.Sprintf("%d:%d - %s", p.t.Row, p.t.Col, p.cause.Error())
}

func unexpectedError(t *Tokenizer, expected string, actual []byte) *ParserError {
	return &ParserError{
		t:     t,
		cause: fmt.Errorf("unexpected %q. I expected %s", string(actual), expected),
	}
}
