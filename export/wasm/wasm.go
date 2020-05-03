package wasm

import (
	"fmt"
	"time"

	"github.com/mariomac/msxmml/export/wasm/webaudio"
	"github.com/mariomac/msxmml/song"
	"github.com/mariomac/msxmml/song/note"
)

func ExportWasm(s *song.Song) {
	ctx := webaudio.WindowAudioContext()
	for _, c := range s.Channels {
		exportChannel(ctx, c)
	}
}

var flatEquivs = map[note.Pitch]note.Pitch{
	note.A: note.G,
	note.B: note.A,
	note.C: note.B,
	note.D: note.C,
	note.E: note.D,
	note.F: note.E,
	note.G: note.F,
}

func exportChannel(ctx *webaudio.AudioContext, c *song.Channel) {
	ch := NewChannel(ctx, c.Instrument)
	sixteenths := float64(0) // todo: consider higher?
	for _, nt := range c.Notes {
		if nt.Pitch != note.Silence {
			var pitch string
			switch nt.Halftone {
			case note.Sharp:
				pitch = fmt.Sprintf("%c#%d", nt.Pitch, nt.Octave)
			case note.Flat:
				pitch = fmt.Sprintf("%c#%d", flatEquivs[nt.Pitch], nt.Octave)
			default:
				pitch = fmt.Sprintf("%c%d", nt.Pitch, nt.Octave)
			}
			ch.Play(Note{Pitch: pitch, Time: time.Duration(sixteenths/8.0 * float64(time.Second) * 120/bpm)})
		}
		length := 16.0 / float64(nt.Length)
		if nt.Tuplet >= 3 { // todo consider 5-tuples etc...
			length *= float64(nt.Tuplet-1) / float64(nt.Tuplet)
		}
		sixteenths += length
	}
}
