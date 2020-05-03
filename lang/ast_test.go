package lang

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTwoChannelParse(t *testing.T) {
	mml := `
@foo <- abcdefgo2<ab#4
     abcdeefga8

@1<-acbcdedgo2>ab#4
     abcdeebbfga38
`
	s, err := Parse(NewTokenizer(bytes.NewReader([]byte(mml))))
	require.NoError(t, err)
	require.Len(t, s.Channels, 2)
	require.Contains(t, s.Channels, "foo")
	ch := s.Channels["foo"]
	assert.Equal(t, "foo", ch.Name)
	assert.Len(t, ch.Notes, 18)

	require.Contains(t, s.Channels, "1")
	ch = s.Channels["1"]
	assert.Equal(t, "1", ch.Name)
	assert.Len(t, ch.Notes, 21)

}
