package song

import (
	"github.com/mariomac/msxmml/pkg/song/note"
	"time"
)

type Song struct {
	Instruments map[string]Instrument

	Channels map[string]*Channel
}

type TablatureItem struct {
	VariableRef *string
	Note *note.Note
}

type Channel struct {
	Status struct {
		Octave int
	}
	Name       string
	Notes      []note.Note
	Instrument Instrument
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
