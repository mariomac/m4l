package wasm

import (
	"fmt"
	"github.com/mariomac/msxmml/pkg/song"
	"time"

	"github.com/mariomac/msxmml/export/wasm/webaudio"
)

func ExportWasm(s *song.Song) {
	ctx := webaudio.WindowAudioContext()
	for _, c := range s.Channels {
		exportChannel(ctx, c)
	}
}

var flatEquivs = map[song.Pitch]song.Pitch{
	song.A: song.G,
	song.B: song.A,
	song.C: song.B,
	song.D: song.C,
	song.E: song.D,
	song.F: song.E,
	song.G: song.F,
}

func exportChannel(ctx *webaudio.AudioContext, c *song.Channel) {
	ch := NewChannel(ctx, c.Instrument)
	sixteenths := float64(0) // todo: consider higher?
	for _, nt := range c.Notes {
		if nt.Pitch != song.Silence {
			var pitch string
			switch nt.Halftone {
			case song.Sharp:
				pitch = fmt.Sprintf("%c#%d", nt.Pitch, nt.Octave)
			case song.Flat:
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
