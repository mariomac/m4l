package lang

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

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
	next := func() Token {
		require.True(t, tok.Next())
		return tok.Get()
	}

	n := func(n string, r, c int) Token {
		return Token{Type: Note, Content: n, Submatch: []string{n, "", "", ""}, Row: r, Col: c}
	}

	assert.Equal(t, Token{Type: ChannelId, Submatch: []string{"0"}, Content: "@0", Row: 2, Col: 1}, next())
	assert.Equal(t, Token{Type: SendArrow, Submatch: []string{}, Content: "<-", Row: 2, Col: 4}, next())
	assert.Equal(t, n("a", 2, 7), next())
	assert.Equal(t, n("b", 2, 8), next())
	assert.Equal(t, n("c", 2, 9), next())
	assert.Equal(t, n("d", 2, 10), next())
	assert.Equal(t, n("e", 2, 11), next())
	assert.Equal(t, n("f", 2, 12), next())
	assert.Equal(t, n("g", 2, 13), next())
	assert.Equal(t, Token{Type: Octave, Submatch: []string{"1"}, Content: `o1`, Row: 2, Col: 14}, next())
	assert.Equal(t, Token{Type: OctaveStep, Submatch: []string{"<"}, Content: `<`, Row: 2, Col: 16}, next())
	assert.Equal(t, n("a", 2, 17), next())
	assert.Equal(t, Token{Type: Note, Content: "b#4.", Submatch: []string{"b", "#", "4", "."}, Row: 2, Col: 18}, next())
	assert.Equal(t, n("a", 3, 7), next())
	assert.Equal(t, n("b", 3, 8), next())
	assert.Equal(t, n("c", 3, 9), next())
	assert.Equal(t, n("d", 3, 10), next())
	assert.Equal(t, Token{Type: Note, Content: "e.", Submatch: []string{"e", "", "", "."}, Row: 3, Col: 11}, next())
	assert.Equal(t, Token{Type: Note, Content: "e8", Submatch: []string{"e", "", "8", ""}, Row: 3, Col: 13}, next())
	assert.Equal(t, n("f", 3, 15), next())
	assert.Equal(t, n("g", 3, 16), next())
	assert.Equal(t, Token{Type: Note, Content: "a16", Submatch: []string{"a", "", "16", ""}, Row: 3, Col: 17}, next())
	assert.Equal(t, Token{Type: ChannelId, Submatch: []string{"troloro"}, Content: "@troloro", Row: 5, Col: 1}, next())
	assert.Equal(t, Token{Type: SendArrow, Submatch: []string{}, Content: "<-", Row: 5, Col: 9}, next())
	assert.Equal(t, n("a", 5, 11), next())
	assert.Equal(t, n("b", 5, 12), next())
	assert.Equal(t, n("c", 5, 13), next())
	assert.Equal(t, n("d", 5, 14), next())
	assert.Equal(t, n("e", 5, 15), next())
	assert.Equal(t, n("d", 5, 16), next())
	assert.Equal(t, n("g", 5, 17), next())
	assert.Equal(t, Token{Type: Octave, Submatch: []string{"1"}, Content: `o1`, Row: 5, Col: 18}, next())
	assert.Equal(t, Token{Type: OctaveStep, Submatch: []string{"<"}, Content: `<`, Row: 5, Col: 20}, next())
	assert.Equal(t, n("a", 5, 21), next())
	assert.Equal(t, n("a", 6, 5), next())
	assert.Equal(t, n("b", 6, 6), next())
	assert.Equal(t, Token{Type: OpenKey, Submatch: []string{}, Content: `{`, Row: 6, Col: 7}, next())
	assert.Equal(t, n("c", 6, 8), next())
	assert.Equal(t, n("d", 6, 9), next())
	assert.Equal(t, n("e", 6, 10), next())
	assert.Equal(t, Token{Type: CloseTuple, Submatch: []string{"3"}, Content: `}3`, Row: 6, Col: 11}, next())
	assert.Equal(t, Token{Type: Note, Content: "e-3..", Submatch: []string{"e", "-", "3", ".."}, Row: 6, Col: 13}, next())
	assert.Equal(t, Token{Type: Separator, Submatch: []string{}, Content: `|`, Row: 6, Col: 19}, next())
	assert.Equal(t, n("a", 6, 21), next())
	assert.Equal(t, n("b", 6, 22), next())
	assert.Equal(t, n("c", 6, 23), next())

	assert.False(t, tok.Next())

}

func TestTokenizer_ConstantDefinitions(t *testing.T) {
	mml := `
$voice := {
    wave: sine
    adsr: 50->100, 100-> 70,200, 10
}
$intro := ab-4..

@ch1 <- $intro$intro$intro
loop:
@ch1 <- c
---
@ch1 <- d e 
`
	tok := NewTokenizer(bytes.NewReader([]byte(mml)))

	next := func() Token {
		require.True(t, tok.Next())
		return tok.Get()
	}
	assert.Equal(t, Token{Type: ConstDef, Submatch: []string{"voice"}, Content: "$voice :=", Row: 2, Col: 1}, next())
	assert.Equal(t, Token{Type: OpenKey, Content: "{", Submatch: []string{}, Row: 2, Col: 11}, next())
	assert.Equal(t, Token{Type: MapEntry, Content: "wave: sine", Submatch: []string{"wave", "sine"}, Row: 3, Col: 5}, next())
	assert.Equal(t, Token{Type: AdsrVector, Content: "adsr: 50->100, 100-> 70,200, 10",
		Submatch: []string{"50", "100", "100", "70", "200", "10"}, Row: 4, Col: 5}, next())
	assert.Equal(t, Token{Type: CloseKey, Content: "}", Submatch: []string{}, Row: 5, Col: 1}, next())
	assert.Equal(t, Token{Type: ConstDef, Content: "$intro :=", Submatch: []string{"intro"}, Row: 6, Col: 1}, next())
	assert.Equal(t, Token{Type: Note, Content: "a", Submatch: []string{"a", "", "", ""}, Row: 6, Col: 11}, next())
	assert.Equal(t, Token{Type: Note, Content: "b-4..", Submatch: []string{"b", "-", "4", ".."}, Row: 6, Col: 12}, next())
	assert.Equal(t, Token{Type: ChannelId, Content: "@ch1", Submatch: []string{"ch1"}, Row: 8, Col: 1}, next())
	assert.Equal(t, Token{Type: SendArrow, Content: "<-", Submatch: []string{}, Row: 8, Col: 6}, next())
	assert.Equal(t, Token{Type: ConstRef, Content: "$intro", Submatch: []string{"intro"}, Row: 8, Col: 9}, next())
	assert.Equal(t, Token{Type: ConstRef, Content: "$intro", Submatch: []string{"intro"}, Row: 8, Col: 15}, next())
	assert.Equal(t, Token{Type: ConstRef, Content: "$intro", Submatch: []string{"intro"}, Row: 8, Col: 21}, next())
	assert.Equal(t, Token{Type: LoopTag, Content: "loop:", Submatch: []string{}, Row: 9, Col: 1}, next())
	assert.Equal(t, Token{Type: ChannelId, Content: "@ch1", Submatch: []string{"ch1"}, Row: 10, Col: 1}, next())
	assert.Equal(t, Token{Type: SendArrow, Content: "<-", Submatch: []string{}, Row: 10, Col: 6}, next())
	assert.Equal(t, Token{Type: Note, Content: "c", Submatch: []string{"c", "", "", ""}, Row: 10, Col: 9}, next())
	assert.Equal(t, Token{Type: ChannelSync, Content: "---", Submatch: []string{}, Row: 11, Col: 1}, next())
	assert.Equal(t, Token{Type: ChannelId, Content: "@ch1", Submatch: []string{"ch1"}, Row: 12, Col: 1}, next())
	assert.Equal(t, Token{Type: SendArrow, Content: "<-", Submatch: []string{}, Row: 12, Col: 6}, next())
	assert.Equal(t, Token{Type: Note, Content: "d", Submatch: []string{"d", "", "", ""}, Row: 12, Col: 9}, next())
	assert.Equal(t, Token{Type: Note, Content: "e", Submatch: []string{"e", "", "", ""}, Row: 12, Col: 11}, next())

	assert.False(t, tok.Next())
}
