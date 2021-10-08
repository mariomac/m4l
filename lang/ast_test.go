package lang

import (
	"bytes"
	"github.com/mariomac/msxmml/pkg/lang"
	"github.com/mariomac/msxmml/pkg/song"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO TEST that a non-closed tuplet returns error
// TODO TEST that closing a non-open tuplet returns error

func TestTwoChannelParse(t *testing.T) {
	mml := `
@foo <- abcdefgo2<ab#4
     abcdeefga8

@1<-acbcdedgo2>ab#4
     abcdeebbfga38
`
	s, err := Parse(lang.NewTokenizer(bytes.NewReader([]byte(mml))))
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

func TestInstruments(t *testing.T) {
	mml := `
@voice {
	wave: sine
	adsr: 30->100, 100->60, 200, 210
}
@voice <- abc
@another <- cda
`

	s, err := Parse(lang.NewTokenizer(bytes.NewReader([]byte(mml))))
	require.NoError(t, err)
	require.Len(t, s.Channels, 2)
	require.Contains(t, s.Channels, "voice")
	ch := s.Channels["voice"]
	assert.Equal(t, "voice", ch.Name)
	assert.Len(t, ch.Notes, 3)
	assert.Equal(t, song.Instrument{
		Wave: "sine",
		Envelope: []song.TimePoint{
			{1, 30 * time.Millisecond},
			{0.6, 100 * time.Millisecond},
			{0.6, 200 * time.Millisecond},
			{0, 210 * time.Millisecond},
		},
	}, ch.Instrument)

	require.Contains(t, s.Channels, "another")
	ch = s.Channels["another"]
	assert.Equal(t, "another", ch.Name)
	assert.Len(t, ch.Notes, 3)
	assert.Equal(t, song.DefaultInstrument, ch.Instrument)

}
