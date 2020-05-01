// +build wasm

package webaudio

import (
	"syscall/js"
	"time"
)

type AudioContext struct {
	val js.Value
}

func WindowAudioContext() *AudioContext {
	ac := js.Global().Get("AudioContext").New()
	return &AudioContext{val: ac}
}

func (ac *AudioContext) Time() time.Duration {
	seconds := ac.val.Get("currentTime").Float()
	return time.Duration(seconds * float64(time.Second))
}

// oscillator --> gain --> ctx.destination
func (ac *AudioContext) NoteNodes(freq float64, adsr ADSR) *OscillatorNode {
	gain := ac.val.Call("createGain")
	gain.Call("connect", ac.val.Get("destination"))
	gainObj := &Gain{val: gain}
	gainObj.SetValueAtTime(0, ac.Time())
	osc := ac.val.Call("createOscillator")
	osc.Set("type", "square")
	osc.Call("connect", gain)
	osc.Get("frequency").Set("value", freq)
	osc.Call("start", ac.Time().Seconds())
	return &OscillatorNode{
		val:  osc,
		ctx:  ac,
		adsr: adsr,
		gain: gainObj,
	}
}
