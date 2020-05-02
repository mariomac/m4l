package wasm

import (
	"fmt"
	"time"

	"github.com/mariomac/msxmml/export/wasm/webaudio"
	"github.com/mariomac/msxmml/solfa"
	"github.com/mariomac/msxmml/song"
)

func ExportWasm(s *song.Song) {
	ctx := webaudio.WindowAudioContext()
	for _, c := range s.Channels {
		exportChannel(ctx, c)
	}
}

var flatEquivs = map[solfa.Pitch]solfa.Pitch {
	solfa.A: solfa.G,
	solfa.B: solfa.A,
	solfa.C: solfa.B,
	solfa.D: solfa.C,
	solfa.E: solfa.D,
	solfa.F: solfa.E,
	solfa.G: solfa.F,
}

func exportChannel(ctx *webaudio.AudioContext, c song.Channel) {
	ch := NewChannel(ctx, Instrument{
		adsr: webaudio.ADSR{
			{1, 50 * time.Millisecond},
			{0.7, 100 * time.Millisecond},
			{0.7, 200 * time.Millisecond},
			{0, 250 * time.Millisecond},
		}})
	sixteenths := float64(0) // todo: consider higher?
	for _, note := range c.Notes {
		if note.Pitch != solfa.Silence {
			var pitch string
			switch note.Halftone {
			case solfa.Sharp:
				pitch = fmt.Sprintf("%c#%d", note.Pitch, note.Octave)
			case solfa.Flat:
				pitch = fmt.Sprintf("%c#%d", flatEquivs[note.Pitch], note.Octave)
			default:
				pitch = fmt.Sprintf("%c%d", note.Pitch, note.Octave)
			}
			ch.Play(Note{Pitch: pitch, Time: time.Duration(sixteenths/8.0 * float64(time.Second) * 120/bpm)})
		}
		length := 16.0 / float64(note.Length)
		if note.Tuplet >= 3 { // todo consider 5-tuples etc...
			length *= float64(note.Tuplet-1) / float64(note.Tuplet)
		}
		sixteenths += length
	}
}
