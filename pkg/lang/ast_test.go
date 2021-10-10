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

func TestCompleteProgram(t *testing.T) {
	s, err := Parse(NewTokenizer(strings.NewReader(`
$voice := {
	wave: sine
	adsr: 30->100, 100->60, 200, 210
}
$const := abc

@ch1 <- c1.d-2..e+4$const$const
loop:
@ch1 <- r4a>
@ch2 <- aco2 | d
---
@ch1 <- {dec}3
`)))
	require.NoError(t, err)
	require.NotNil(t, s)

}
