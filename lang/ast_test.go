package lang

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTwoChannelParse(t *testing.T) {
	mml := `
@0 <- abcdefgo2<ab#4
      abcdeefga8

@1<-acbcdedgo2>ab#4
      abcdeebbfga38
`
	s, err := Root(NewTokenizer([]byte(mml)))
	require.NoError(t, err)
	assert.Len(t, s.Channels, 2)
	assert.Contains(t, s.Channels, 0)
	assert.Contains(t, s.Channels, 1)
	ch := s.Channels[0]
	assert.Equal(t, 0, ch.Number)
	assert.Len(t, ch.Notes, 18)
	ch = s.Channels[1]
	assert.Equal(t, 1, ch.Number)
	assert.Len(t, ch.Notes, 21)

}
