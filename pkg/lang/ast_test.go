package lang

import (
	"strings"
	"testing"

	"github.com/mariomac/msxmml/pkg/song/note"

	"github.com/mariomac/msxmml/pkg/song"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO TEST that a non-closed tuplet returns error
// TODO TEST that closing a non-open tuplet returns error

func TestInstrument(t *testing.T) {
	s, err := Parse(NewTokenizer(strings.NewReader(`
$voice := trumpet {
	wave: sine
	sordine: true
}
`)))
	require.NoError(t, err)
	require.Contains(t, s.Constants, "voice")
	require.Len(t, s.Constants["voice"], 1)
	voice := s.Constants["voice"][0]
	require.NotNil(t, voice.Instrument)
	assert.Equal(t, song.Instrument{
		Class: "trumpet",
		Properties: map[string]string{
			"wave":    "sine",
			"sordine": "true",
		},
	}, *voice.Instrument)
}

func TestCompleteProgram(t *testing.T) {
	s, err := Parse(NewTokenizer(strings.NewReader(`
$voice := triki {
	wave: sine
	adsr: traka
}
$const := abc

@ch1 <- c1.d-2..e+4$const$const
loop:
@ch1 <- v14r4a>
@ch2 <- v13aco2 | d
---
@ch1 <- (dec)3
`)))
	require.NoError(t, err)
	require.NotNil(t, s)
	// check overall tablature structure
	require.Len(t, s.Constants, 2)
	require.Len(t, s.Blocks, 3)
	// check $voice constant definition
	require.Len(t, s.Constants["voice"], 1)
	assert.Equal(t,
		&song.Instrument{
			Class: "triki",
			Properties: map[string]string{
				"wave": "sine",
				"adsr": "traka",
			},
		},
		s.Constants["voice"][0].Instrument)
	// check $const constant definition
	require.Len(t, s.Constants["const"], 3)
	for n, exp := range []note.Pitch{note.A, note.B, note.C} {
		assert.Equal(t,
			&note.Note{Pitch: exp, Length: defaultLength},
			s.Constants["const"][n].Note)
	}
	// check @ch1 <- c1.d-2..e+4$const$const
	require.Len(t, s.Blocks[0].Channels, 1)
	require.Contains(t, s.Blocks[0].Channels, "ch1")
	require.Len(t, s.Blocks[0].Channels["ch1"].Items, 5)
	assert.Equal(t,
		&note.Note{Pitch: note.C, Length: 1, Dots: 1},
		s.Blocks[0].Channels["ch1"].Items[0].Note)
	assert.Equal(t,
		&note.Note{Pitch: note.D, Length: 2, Dots: 2, Halftone: note.Flat},
		s.Blocks[0].Channels["ch1"].Items[1].Note)
	assert.Equal(t,
		&note.Note{Pitch: note.E, Length: 4, Halftone: note.Sharp},
		s.Blocks[0].Channels["ch1"].Items[2].Note)
	cref := "const"
	assert.Equal(t, &cref, s.Blocks[0].Channels["ch1"].Items[3].ConstantRef)
	assert.Equal(t, &cref, s.Blocks[0].Channels["ch1"].Items[4].ConstantRef)

	//	check loop label
	assert.Equal(t, 1, s.LoopIndex)

	// check synced block ch1 and ch2 channel statements
	require.Len(t, s.Blocks[1].Channels, 2)
	// @ch1 <- v14r4a>
	require.Len(t, s.Blocks[1].Channels["ch1"].Items, 4)
	assert.Equal(t,
		14,
		*s.Blocks[1].Channels["ch1"].Items[0].Volume)
	assert.Equal(t,
		&note.Note{Pitch: note.Silence, Length: 4},
		s.Blocks[1].Channels["ch1"].Items[1].Note)
	assert.Equal(t,
		&note.Note{Pitch: note.A, Length: defaultLength},
		s.Blocks[1].Channels["ch1"].Items[2].Note)
	assert.Equal(t,
		1,
		*s.Blocks[1].Channels["ch1"].Items[3].OctaveStep)

	// @ch2 <- v13aco2 | d
	require.Len(t, s.Blocks[1].Channels["ch2"].Items, 5)
	assert.Equal(t,
		13,
		*s.Blocks[1].Channels["ch2"].Items[0].Volume)
	assert.Equal(t,
		&note.Note{Pitch: note.A, Length: defaultLength},
		s.Blocks[1].Channels["ch2"].Items[1].Note)
	assert.Equal(t,
		&note.Note{Pitch: note.C, Length: defaultLength},
		s.Blocks[1].Channels["ch2"].Items[2].Note)
	assert.Equal(t,
		2,
		*s.Blocks[1].Channels["ch2"].Items[3].SetOctave)
	assert.Equal(t,
		&note.Note{Pitch: note.D, Length: defaultLength},
		s.Blocks[1].Channels["ch2"].Items[4].Note)

	// check synced block after barrier
	// @ch1 <- {dec}3
	assert.Equal(t,
		&note.Note{Pitch: note.D, Tuplet: 3, Length: defaultLength},
		s.Blocks[2].Channels["ch1"].Items[0].Note)
	assert.Equal(t,
		&note.Note{Pitch: note.E, Tuplet: 3, Length: defaultLength},
		s.Blocks[2].Channels["ch1"].Items[1].Note)
	assert.Equal(t,
		&note.Note{Pitch: note.C, Tuplet: 3, Length: defaultLength},
		s.Blocks[2].Channels["ch1"].Items[2].Note)
}
