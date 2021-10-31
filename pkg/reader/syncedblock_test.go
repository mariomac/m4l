package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mariomac/msxmml/pkg/song"
)

func TestSyncedBlockReader(t *testing.T) {
	one := 1
	sb := song.SyncedBlock{Channels: map[string]*song.Channel{}}
	sb.Channels["a"] = &song.Channel{Items: []song.TablatureItem{
		{Instrument: &song.Instrument{Class: "psg"}},
		{Silence: &song.Silence{Length: 1}},
		{Note: &song.Note{Pitch: song.A, Length: 4}},
		{Note: &song.Note{Pitch: song.B, Length: 4}},
	}}
	sb.Channels["b"] = &song.Channel{Items: []song.TablatureItem{
		{Instrument: &song.Instrument{Class: "psg2"}},
		{Note: &song.Note{Pitch: song.C, Length: 2}},
		{OctaveStep: &one},
		{Silence: &song.Silence{Length: 2}},
		{Note: &song.Note{Pitch: song.F, Length: 4}},
	}}
	sb.Channels["c"] = &song.Channel{Items: []song.TablatureItem{
		{SetOctave: &one},
		{Silence: &song.Silence{Length: 2}},
		{Note: &song.Note{Pitch: song.D, Length: 4}},
		{Note: &song.Note{Pitch: song.E, Length: 4}},
	}}
	reader := NewSyncedBlock(sb)

	read := map[string][]song.TablatureItem{}
	for i := 0; i < 6; i++ {
		ti, ch := reader.Next()
		read[ch] = append(read[ch], ti)
	}
	assert.Equal(t, read, map[string][]song.TablatureItem{
		"a": {
			song.TablatureItem{Instrument: &song.Instrument{Class: "psg"}},
			song.TablatureItem{Silence: &song.Silence{Length: 1}},
		},
		"b": {
			song.TablatureItem{Instrument: &song.Instrument{Class: "psg2"}},
			song.TablatureItem{Note: &song.Note{Pitch: song.C, Length: 2}},
		},
		"c": {
			song.TablatureItem{SetOctave: &one},
			song.TablatureItem{Silence: &song.Silence{Length: 2}},
		},
	})

	read = map[string][]song.TablatureItem{}
	for i := 0; i < 3; i++ {
		ti, ch := reader.Next()
		read[ch] = append(read[ch], ti)
	}
	assert.Equal(t, read, map[string][]song.TablatureItem{
		"b": {
			song.TablatureItem{OctaveStep: &one},
			song.TablatureItem{Silence: &song.Silence{Length: 2}},
		},
		"c": {song.TablatureItem{Note: &song.Note{Pitch: song.D, Length: 4}}},
	})
	ti, ch := reader.Next()
	assert.Equal(t, "c", ch)
	assert.Equal(t, song.TablatureItem{Note: &song.Note{Pitch: song.E, Length: 4}}, ti)
	read = map[string][]song.TablatureItem{}
	for i := 0; i < 2; i++ {
		ti, ch = reader.Next()
		read[ch] = append(read[ch], ti)
	}
	assert.Equal(t, read, map[string][]song.TablatureItem{
		"a": {song.TablatureItem{Note: &song.Note{Pitch: song.A, Length: 4}}},
		"b": {song.TablatureItem{Note: &song.Note{Pitch: song.F, Length: 4}}},
	})
	ti, ch = reader.Next()
	assert.Equal(t, "a", ch)
	assert.Equal(t, song.TablatureItem{Note: &song.Note{Pitch: song.B, Length: 4}}, ti)

	// end of tablature. Nothing else to read
	_, ch = reader.Next()
	assert.Empty(t, ch)
}

