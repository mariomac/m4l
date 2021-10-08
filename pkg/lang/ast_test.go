package lang

import (
	"strings"
	"testing"
	"time"

	"github.com/mariomac/msxmml/pkg/song"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO TEST that a non-closed tuplet returns error
// TODO TEST that closing a non-open tuplet returns error

/*
func TestTwoChannelParse(t *testing.T) {
	mml := `
@foo <- abcdefgo2<ab#4
     abcdeefga8

@1<-acbcdedgo2>ab#4
     abcdeebbfga38
`
	s, err := lang.Parse(lang.NewTokenizer(bytes.NewReader([]byte(mml))))
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
}*/

func TestInstrument(t *testing.T) {
	s, err := Parse(NewTokenizer(strings.NewReader(`
$voice := {
	wave: sine
	adsr: 30->100, 100->60, 200, 210
}
`)))
	require.NoError(t, err)
	require.Contains(t, s.Constants, "voice")
	require.Len(t, s.Constants["voice"], 1)
	voice := s.Constants["voice"][0]
	require.NotNil(t, voice.Instrument)
	assert.Equal(t, song.Instrument{
		Wave: "sine",
		Envelope: []song.TimePoint{
			{1, 30 * time.Millisecond},
			{0.6, 100 * time.Millisecond},
			{0.6, 200 * time.Millisecond},
			{0, 210 * time.Millisecond},
		},
	}, *voice.Instrument)

}
