package song

import (
	"time"

	"github.com/mariomac/msxmml/pkg/song/note"
)

type Song struct {
	Constants map[string]*TablatureItem
	Channels  map[string]*Channel
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
	Status struct {
		Octave     int
		Instrument Instrument
	}
	Items []TablatureItem
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
