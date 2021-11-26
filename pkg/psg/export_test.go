package psg

import (
	"strings"
	"testing"

	"github.com/mariomac/msxmml/pkg/lang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportNotes(t *testing.T) {
	song, err := lang.Parse(strings.NewReader(`
@ch1 <- |a1|  |  |  |b2|  |c4|  |
@ch2 <- |d2|  |e1|  |  |  |f4.|  |
@ch3 <- |a4|b1|  |  |  |d2|  |  |
`))
	require.NoError(t, err)
	songBytes, err := Export(song)
	require.NoError(t, err)
	expected := append([]byte{0, 0}, encodeInstructions([]instruction{
		{Type: channels, Data: 0b111_110},
		{Type: toneA, Data: 0xfe},
		{Type: channels, Data: 0b111_100},
		{Type: toneB, Data: 0x17d},
		{Type: channels, Data: 0b111_000},
		{Type: toneC, Data: 0xfe},
		{Type: wait, Data: 30},
		{Type: toneC, Data: 0xe3},
		{Type: wait, Data: 30},
		{Type: toneB, Data: 0x153}, // e1
		{Type: wait, Data: 31},
		{Type: wait, Data: 29},
		{Type: toneA, Data: 0xE3}, //b2
		{Type: wait, Data: 30},
		{Type: toneC, Data: 0x17d},
		{Type: wait, Data: 30},
		{Type: toneA, Data: 0x1ac},
		{Type: toneB, Data: 0x140},
		{Type: wait, Data: 30},
		{Type: wait, Data: 15}, // at the end of the block, syncing to the added dot
		{Type: end},
	})...)

	assert.Equal(t, expected, songBytes)
}

func TestOctaves(t *testing.T) {
	song, err := lang.Parse(strings.NewReader(`
tempo 60
psg.hz 50
@ch1 <- o5a>b
@ch2 <- c<d
`))
	require.NoError(t, err)
	songBytes, err := Export(song)
	require.NoError(t, err)
	expected := append([]byte{0, 0}, encodeInstructions([]instruction{
		{Type: channels, Data: 0b111_110},
		{Type: toneA, Data: 0x7F}, // octave 5 a
		{Type: channels, Data: 0b111_100},
		{Type: toneB, Data: 0x1AC}, // octave 4 c
		{Type: wait, Data: 31},
		{Type: wait, Data: 19},
		{Type: toneA, Data: 0x39},  // octave 6 b
		{Type: toneB, Data: 0x2FA}, // octave 3 d
		{Type: wait, Data: 31},
		{Type: wait, Data: 19},
		{Type: end},
	})...)
	assert.Equal(t, expected, songBytes)
}

func TestSilences(t *testing.T) {
	song, err := lang.Parse(strings.NewReader(`
@a <- r1 a r2 b r4 c r8
`))
	require.NoError(t, err)
	songBytes, err := Export(song)
	require.NoError(t, err)
	expected := append([]byte{0, 0}, encodeInstructions([]instruction{
		{Type: wait, Data: 31}, // 4 beats waiting
		{Type: wait, Data: 31},
		{Type: wait, Data: 31},
		{Type: wait, Data: 27},

		{Type: channels, Data: 0b111_110},
		{Type: toneA, Data: 0xfe},
		{Type: wait, Data: 30}, // wait for the note

		{Type: channels, Data: 0b111_111},
		{Type: wait, Data: 31}, // 2 beats silence waiting
		{Type: wait, Data: 29},

		{Type: channels, Data: 0b111_110},
		{Type: toneA, Data: 0xE3},
		{Type: wait, Data: 30},

		{Type: channels, Data: 0b111_111},
		{Type: wait, Data: 30}, // 1 beat silence waiting

		{Type: channels, Data: 0b111_110},
		{Type: toneA, Data: 0x1ac},
		{Type: wait, Data: 30},

		{Type: channels, Data: 0b111_111},
		{Type: wait, Data: 15}, // 1/2 beat silence waiting
		{Type: end},
	})...)
	assert.Equal(t, expected, songBytes)

}

func TestParseTupletWithOctaveChange(t *testing.T) {
	s, err := lang.Parse(strings.NewReader(`
@ch1 <- (a>aa)3 a
`))
	require.NoError(t, err)
	songBytes, err := Export(s)
	require.NoError(t, err)
	expected := append([]byte{0, 0}, encodeInstructions([]instruction{
		{Type: channels, Data: 0b111_110},
		{Type: toneA, Data: 0xfe}, // octave 4
		{Type: wait, Data: 20},
		{Type: toneA, Data: 0x7F}, // octave 5
		{Type: wait, Data: 20},
		{Type: toneA, Data: 0x7F}, // octave 5
		{Type: wait, Data: 20},
		{Type: toneA, Data: 0x7F}, // octave 5
		{Type: wait, Data: 30},
		{Type: end},
	})...)
	assert.Equal(t, expected, songBytes)
}

func TestExportLoop(t *testing.T) {
	s, err := lang.Parse(strings.NewReader(`
@ch1 <- a b
loop:
@ch1 <- > c d
`)) // todo: document that the > won't increase the octave undefinitely
	require.NoError(t, err)
	songBytes, err := Export(s)
	require.NoError(t, err)
	expected := append([]byte{9, 0}, encodeInstructions([]instruction{
		{Type: channels, Data: 0b111_110},		// songBytes[2],
		{Type: toneA, Data: 0xFE}, // o4 a		   songBytes[3]
		{Type: wait, Data: 30},                 // songBytes[5],
		{Type: toneA, Data: 0xE3}, // o4 b         songBytes[6]
		{Type: wait, Data: 30},				    // songBytes[8],
		{Type: toneA, Data: 0xD6}, // o5 c         songBytes[9] <-- loop here!
		{Type: wait, Data: 30},
		{Type: toneA, Data: 0xBE}, // o5 d
		{Type: wait, Data: 30},
		{Type: end},
	})...)
	assert.Equal(t, expected, songBytes)

}