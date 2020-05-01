package wasm

import (
	"log"
	"time"

	"github.com/mariomac/msxmml/export/wasm/webaudio"
)

// TODO: share and configure
const bpm = 180

type Instrument struct {
	adsr  webaudio.ADSR
}

type ChannelNodes struct {
	ctx   *webaudio.AudioContext
	inst  Instrument
	notes map[string]*webaudio.OscillatorNode
}

mirar si zurula con un solo oscilador y usando "frequency.setValueAtTime"
func NewChannel(ctx *webaudio.AudioContext, i Instrument) *ChannelNodes {
	return &ChannelNodes{
		ctx:   ctx,
		inst: i,
		notes: map[string]*webaudio.OscillatorNode{},
	}
}

type Note struct {
	Pitch string
	Time  time.Duration
}

func (cn *ChannelNodes) Play(n Note) {
	log.Printf("pitch: %s (freq: %f)", n.Pitch, freqs[n.Pitch])
	node, ok := cn.notes[n.Pitch]
	if !ok {
		node = cn.ctx.NoteNodes(freqs[n.Pitch], cn.inst.adsr)
		cn.notes[n.Pitch] = node
	}
	node.TriggerEnvelope(n.Time)
}
