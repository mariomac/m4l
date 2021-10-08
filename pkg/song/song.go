package song

import (
	"time"

	"github.com/mariomac/msxmml/pkg/song/note"
)

const defaultOctave = 4

// TODO: test all
type Song struct {
	Constants map[string]Tablature
	Blocks    []SyncedBlock
}

type Tablature []TablatureItem

func (s *Song) AddSyncedBlock() {
	s.Blocks = append(s.Blocks, SyncedBlock{Channels: map[string]*Channel{}})
}

func (s *Song) endBlock() *SyncedBlock {
	if len(s.Blocks) == 0 {
		return nil
	}
	return &s.Blocks[len(s.Blocks)-1]
}

func (s *Song) AddItems(channelName string, items ...TablatureItem) {
	if len(s.Blocks) == 0 {
		s.AddSyncedBlock()
	}
	ch, ok := s.endBlock().Channels[channelName]
	if !ok || ch == nil {
		ch = &Channel{}
		s.endBlock().Channels[channelName] = ch
	}
	ch.Items = append(ch.Items, items...)
}

// TablatureItem pseudo-union type: whatever you can find in a tablature
type TablatureItem struct {
	Instrument  *Instrument
	VariableRef *string
	Note        *note.Note
	SetOctave   *int
	IncOctave   *int // negative: decrements
}

type Channel struct {
	Items []TablatureItem
}

// SyncedBlock contains channels that sound at the same time. The SyncedBlock hasn't finished
// until all the channels finish
type SyncedBlock struct {
	Channels map[string]*Channel
}

type Instrument struct {
	Wave     string
	Envelope []TimePoint // attack decay sustain release
}

type TimePoint struct {
	Val  float64
	Time time.Duration
}

var DefaultInstrument = Instrument{
	Wave: "square",
	Envelope: []TimePoint{
		{1, 50 * time.Millisecond},
		{0.7, 100 * time.Millisecond},
		{0.7, 200 * time.Millisecond},
		{0, 250 * time.Millisecond},
	},
}
