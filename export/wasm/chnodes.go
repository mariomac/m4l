package wasm

import (
	"github.com/mariomac/msxmml/pkg/song"
	"time"

	"github.com/mariomac/msxmml/export/wasm/webaudio"
)

// TODO: share and configure
const bpm = 180

type ChannelNodes struct {
	ctx   *webaudio.AudioContext
	inst  song.Instrument
	notes *webaudio.OscillatorNode
}

//mirar si zurula con un solo oscilador y usando "frequency.setValueAtTime"
func NewChannel(ctx *webaudio.AudioContext, i song.Instrument) *ChannelNodes {
	return &ChannelNodes{
		ctx:  ctx,
		inst: i,
	}
}

type Note struct {
	Pitch string
	Time  time.Duration
}

func (cn *ChannelNodes) Play(n Note) {
	if cn.notes == nil {
		cn.notes = cn.ctx.NoteNodes(freqs[n.Pitch], cn.inst)
	}
	cn.notes.TriggerEnvelope(freqs[n.Pitch], n.Time)
}
