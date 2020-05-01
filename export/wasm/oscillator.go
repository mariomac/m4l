// +build wasm

package wasm

import (
	"syscall/js"
	"time"
)

type ADSR [4]struct { // attack decay sustain release
	Val  float64
	Time time.Duration
}

type OscillatorNode struct {
	val  js.Value
	ctx  *AudioContext
	adsr ADSR
	gain *Gain
}

func (on *OscillatorNode) Envelope(adsr ADSR) {
	on.adsr = adsr
}

func (on *OscillatorNode) Frequency(f float64) {
	on.val.Get("frequency").Set("value", f)
}

func (on *OscillatorNode) TriggerEnvelope(t time.Duration) {
	ct := on.ctx.Time()
	on.gain.SetValueAtTime(0, ct+t)
	for _, a := range on.adsr {
		on.gain.LinearRampToValueAtTime(a.Val, ct+t+a.Time)
	}
}
