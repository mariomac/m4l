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

	item, name := reader.Next()
	assert.Equal(t, "a", name)
	assert.Equal(t, song.TablatureItem{Instrument: &song.Instrument{Class: "psg"}}, item)
	item, name = reader.Next()
	assert.Equal(t, "a", name)
	assert.Equal(t, song.TablatureItem{Silence: &song.Silence{Length: 1}}, item)
	item, name = reader.Next()
	assert.Equal(t, "b", name)
	assert.Equal(t, song.TablatureItem{Instrument: &song.Instrument{Class: "psg2"}}, item)
	item, name = reader.Next()
	assert.Equal(t, "b", name)
	assert.Equal(t, song.TablatureItem{Note: &song.Note{Pitch: song.C, Length: 2}}, item)
	item, name = reader.Next()
	assert.Equal(t, "c", name)
	assert.Equal(t, song.TablatureItem{SetOctave: &one}, item)
	item, name = reader.Next()
	assert.Equal(t, "c", name)
	assert.Equal(t, song.TablatureItem{Silence: &song.Silence{Length: 2}}, item)
	item, name = reader.Next()
	assert.Equal(t, "b", name)
	assert.Equal(t, song.TablatureItem{OctaveStep: &one}, item)
	item, name = reader.Next()
	assert.Equal(t, "b", name)
	assert.Equal(t, song.TablatureItem{Silence: &song.Silence{Length: 2}}, item)
	item, name = reader.Next()
	assert.Equal(t, "c", name)
	assert.Equal(t, song.TablatureItem{Note: &song.Note{Pitch: song.D, Length: 4}}, item)
	item, name = reader.Next()
	assert.Equal(t, "c", name)
	assert.Equal(t, song.TablatureItem{Note: &song.Note{Pitch: song.E, Length: 4}}, item)

	item, name = reader.Next()
	assert.Equal(t, "a", name)
	assert.Equal(t, song.TablatureItem{Note: &song.Note{Pitch: song.A, Length: 4}}, item)
	item, name = reader.Next()
	assert.Equal(t, "b", name)
	assert.Equal(t, song.TablatureItem{Note: &song.Note{Pitch: song.F, Length: 4}}, item)

	item, name = reader.Next()
	assert.Equal(t, "a", name)
	assert.Equal(t, song.TablatureItem{Note: &song.Note{Pitch: song.B, Length: 4}}, item)

	// end of tablature. Nothing else to read
	_, name = reader.Next()
	assert.Empty(t, name)
}

