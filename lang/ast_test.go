package lang

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTwoChannelParse(t *testing.T) {
	mml := `
@0 <- abcdefgo2<ab#4
      abcdeefghwa8

@1 <- abcdedgo2>ab4#
      abcdeebbfghwa38
`
	s, err := Root(NewTokenizer([]byte(mml)))
	require.NoError(t, err)
	assert.Len(t, s, 2)
	assert.Contains(t, s, 0)
	assert.Contains(t, s, 1)
	ch := s.Channels[0]
	assert.Equal(t, 0, ch.Number)
	assert.Len(t, ch.Notes, 10)

}
