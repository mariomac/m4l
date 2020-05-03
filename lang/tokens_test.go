package lang

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer_TwoChannels(t *testing.T) {
	mml := `
@0 <- abcdefgo1<ab#4.
      abcde.e8fga16

@troloro<-abcdedgo1<a
    ab{cde}3e-3.. | abc
`
	tok := NewTokenizer(bytes.NewReader([]byte(mml)))
	var tokens []Token
	for tok.Next() {
		tokens = append(tokens, tok.Get())
	}

	n := func(n string, r, c int) Token {
		return Token{Type: Note, Content: n, Submatch: []string{n, "", "", ""}, Row: r, Col: c}
	}
	assert.Equal(t, []Token{
		{Type: ChannelID, Submatch: []string{"0"}, Content: "@0", Row: 2, Col: 1},
		{Type: ChannelSendArrow, Submatch: []string{}, Content: "<-", Row: 2, Col: 4},
		n("a", 2, 7), n("b", 2, 8), n("c", 2, 9), n("d", 2, 10), n("e", 2, 11), n("f", 2, 12), n("g", 2, 13),
		{Type: Octave, Submatch: []string{"1"}, Content: `o1`, Row: 2, Col: 14},
		{Type: DecOctave, Submatch: []string{}, Content: `<`, Row: 2, Col: 16},
		n("a", 2, 17),
		{Type: Note, Content: "b#4.", Submatch: []string{"b", "#", "4", "."}, Row: 2, Col: 18},
		n("a", 3, 7), n("b", 3, 8), n("c", 3, 9), n("d", 3, 10),
		{Type: Note, Content: "e.", Submatch: []string{"e", "", "", "."}, Row: 3, Col: 11},
		{Type: Note, Content: "e8", Submatch: []string{"e", "", "8", ""}, Row: 3, Col: 13},
		n("f", 3, 15), n("g", 3, 16),
		{Type: Note, Content: "a16", Submatch: []string{"a", "", "16", ""}, Row: 3, Col: 17},
		{Type: ChannelID, Submatch: []string{"troloro"}, Content: "@troloro", Row: 5, Col: 1},
		{Type: ChannelSendArrow, Submatch: []string{}, Content: "<-", Row: 5, Col: 9},
		n("a", 5, 11), n("b", 5, 12), n("c", 5, 13), n("d", 5, 14), n("e", 5, 15), n("d", 5, 16), n("g", 5, 17),
		{Type: Octave, Submatch: []string{"1"}, Content: `o1`, Row: 5, Col: 18},
		{Type: DecOctave, Submatch: []string{}, Content: `<`, Row: 5, Col: 20},
		n("a", 5, 21), n("a", 6, 5), n("b", 6, 6),
		{Type: OpenSection, Submatch: []string{}, Content: `{`, Row: 6, Col: 7},
		n("c", 6, 8), n("d", 6, 9), n("e", 6, 10),
		{Type: CloseTuplet, Submatch: []string{"3"}, Content: `}3`, Row: 6, Col: 11},
		{Type: Note, Content: "e-3..", Submatch: []string{"e", "-", "3", ".."}, Row: 6, Col: 13},
		{Type: Separator, Submatch: []string{}, Content: `|`, Row: 6, Col: 19},
		n("a", 6, 21), n("b", 6, 22), n("c", 6, 23)},
		tokens)
}
