// +build wasm

package webaudio

import (
	"syscall/js"
	"time"

	"github.com/mariomac/msxmml/song"
)

type OscillatorNode struct {
	val  js.Value
	ctx  *AudioContext
	inst song.Instrument
	gain *Gain
}

func (on *OscillatorNode) FrequencyAtTime(f float64, t time.Duration) {
	on.val.Get("frequency").Call("setValueAtTime", f, t.Seconds())
}

// todo: cancelScheduledValues para evitar que se solapen sonidos
func (on *OscillatorNode) TriggerEnvelope(frequency float64, t time.Duration) {
	ct := on.ctx.Time()
	on.gain.val.Get("gain").Call("cancelScheduledValues", t.Seconds())
	on.FrequencyAtTime(frequency, t)
	on.gain.SetValueAtTime(0, ct+t)
	for _, a := range on.inst.Envelope {
		on.gain.LinearRampToValueAtTime(a.Val, ct+t+a.Time)
	}
}
