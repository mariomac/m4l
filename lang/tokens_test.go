package lang

import "testing"
import "github.com/stretchr/testify/assert"

func TestTokenizer_TwoChannels(t *testing.T) {
	mml := `
@0 <- abcdefgo1<ab4#
      abcdeefghwa3322

@1 <- abcdedgo1<ab4#
      abcdeebbfghwa3322
`
	tok := NewTokenizer([]byte(mml))
	var tokens []Token
	for hasNext := tok.Next() ; hasNext ; hasNext = tok.Next() {
		tokens = append(tokens, tok.Get())
	}
	assert.Equal(t, []Token{
		{Type: Channel, Channel: ChannelToken{Number: 0}},
		{Type: ChannelSendArrow},
		{Type: String, String: StringToken{Content: []byte("abcdefgo1<ab4#")}},
		{Type: String, String: StringToken{Content: []byte("abcdeefghwa3322")}},
		{Type: Channel, Channel: ChannelToken{Number: 1}},
		{Type: ChannelSendArrow},
		{Type: String, String: StringToken{Content: []byte("abcdedgo1<ab4#")}},
		{Type: String, String: StringToken{Content: []byte("abcdeebbfghwa3322")}}},
		tokens)
}
