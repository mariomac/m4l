package lang

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var a Token = Token{Type: Note, Content: []byte(`a`)}

func TestTokenizer_TwoChannels(t *testing.T) {
	mml := `
@0 <- abcdefgo1<ab#4.
      abcde.e8fg

@troloro<-abcdedgo1<a
	ab{cde}3e-3.. | abc
`
	tok := NewTokenizer(bytes.NewReader([]byte(mml)))
	var tokens []Token
	for tok.Next() {
		tokens = append(tokens, tok.Get())
	}

	n := func(n string) Token {
		return Token{Type: Note, Content: []byte(n)}
	}
	assert.Equal(t, []Token{
		{Type: ChannelID, Content: []byte("@0")},
		{Type: ChannelSendArrow, Content: []byte("<-")},
		n("a"), n("b"), n("c"), n("d"), n("e"), n("f"), n("g"),
		{Type: Octave, Content: []byte(`o1`)},
		{Type: DecOctave, Content: []byte(`<`)},
		n("a"), n("b#4."), n("a"), n("b"), n("c"), n("d"), n("e."), n("e8"), n("f"), n("g"),
		{Type: ChannelID, Content: []byte("@troloro")},
		{Type: ChannelSendArrow, Content: []byte("<-")},
		n("a"), n("b"), n("c"), n("d"), n("e"), n("d"), n("g"),
		{Type: Octave, Content: []byte(`o1`)},
		{Type: DecOctave, Content: []byte(`<`)},
		n("a"), n("a"), n("b"),
		{Type: OpenSection, Content: []byte(`{`)},
		n("c"), n("d"), n("e"),
		{Type: CloseTuplet, Content: []byte(`}3`)},
		n("e-3.."),
		{Type: Separator, Content: []byte(`|`)},
		n("a"), n("b"), n("c")},
		tokens)
}
