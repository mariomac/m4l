package lang

import (
	"strings"
	"testing"

	"github.com/mariomac/msxmml/pkg/song"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO TEST that a non-closed tuplet returns error
// TODO TEST that closing a non-open tuplet returns error

func TestInstrument(t *testing.T) {
	s, err := Parse(strings.NewReader(`
$voice := trumpet {
	wave: sine
	sordine: true
}
`))
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
	s, err := Parse(strings.NewReader(`
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
`))
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
	for n, exp := range []song.Pitch{song.A, song.B, song.C} {
		assert.Equal(t,
			&song.Note{Pitch: exp, Length: defaultLength},
			s.Constants["const"][n].Note)
	}
	// check @ch1 <- c1.d-2..e+4$const$const
	require.Len(t, s.Blocks[0].Channels, 1)
	require.Contains(t, s.Blocks[0].Channels, "ch1")
	require.Len(t, s.Blocks[0].Channels["ch1"].Items, 9)
	assert.Equal(t,
		&song.Note{Pitch: song.C, Length: 1, Dots: 1},
		s.Blocks[0].Channels["ch1"].Items[0].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.D, Length: 2, Dots: 2, Halftone: song.Flat},
		s.Blocks[0].Channels["ch1"].Items[1].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.E, Length: 4, Halftone: song.Sharp},
		s.Blocks[0].Channels["ch1"].Items[2].Note)
	/// constants unroll
	assert.Equal(t,
		&song.Note{Pitch: song.A, Length: 4},
		s.Blocks[0].Channels["ch1"].Items[3].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.B, Length: 4},
		s.Blocks[0].Channels["ch1"].Items[4].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.C, Length: 4},
		s.Blocks[0].Channels["ch1"].Items[5].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.A, Length: 4},
		s.Blocks[0].Channels["ch1"].Items[6].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.B, Length: 4},
		s.Blocks[0].Channels["ch1"].Items[7].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.C, Length: 4},
		s.Blocks[0].Channels["ch1"].Items[8].Note)
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
		&song.Silence{Length: 4},
		s.Blocks[1].Channels["ch1"].Items[1].Silence)
	assert.Equal(t,
		&song.Note{Pitch: song.A, Length: defaultLength},
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
		&song.Note{Pitch: song.A, Length: defaultLength},
		s.Blocks[1].Channels["ch2"].Items[1].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.C, Length: defaultLength},
		s.Blocks[1].Channels["ch2"].Items[2].Note)
	assert.Equal(t,
		2,
		*s.Blocks[1].Channels["ch2"].Items[3].SetOctave)
	assert.Equal(t,
		&song.Note{Pitch: song.D, Length: defaultLength},
		s.Blocks[1].Channels["ch2"].Items[4].Note)

	// check synced block after barrier
	// @ch1 <- {dec}3
	assert.Equal(t,
		&song.Note{Pitch: song.D, Tuplet: 3, Length: defaultLength},
		s.Blocks[2].Channels["ch1"].Items[0].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.E, Tuplet: 3, Length: defaultLength},
		s.Blocks[2].Channels["ch1"].Items[1].Note)
	assert.Equal(t,
		&song.Note{Pitch: song.C, Tuplet: 3, Length: defaultLength},
		s.Blocks[2].Channels["ch1"].Items[2].Note)
}

func TestParseTupletWithOctaveChange(t *testing.T) {
	s, err := Parse(strings.NewReader(`
@ch1 <- o4 (ab>c)3 a
`))
	require.NoError(t, err)
	it := s.Blocks[0].Channels["ch1"].Items
	require.Len(t, it, 6)
	require.NotNil(t, it[0].SetOctave)
	assert.Equal(t, 4, *it[0].SetOctave)
	require.NotNil(t, it[1].Note)
	require.Equal(t, song.Note{Pitch: song.A, Length: 4, Tuplet: 3}, *it[1].Note)
	require.NotNil(t, it[2].Note)
	require.Equal(t, song.Note{Pitch: song.B, Length: 4, Tuplet: 3}, *it[2].Note)
	require.NotNil(t, it[3].OctaveStep)
	require.Equal(t, 1, *it[3].OctaveStep)
	require.NotNil(t, it[4].Note)
	require.Equal(t, song.Note{Pitch: song.C, Length: 4, Tuplet: 3}, *it[4].Note)
	require.NotNil(t, it[5].Note)
	require.Equal(t, song.Note{Pitch: song.A, Length: 4}, *it[5].Note)
}