// +build wasm

package webaudio

import (
	"syscall/js"
	"time"
)

type Gain struct {
	val js.Value
}

func (on *Gain) SetValueAtTime(v float64, t time.Duration) {
	on.val.Get("gain").Call("setValueAtTime", v, t.Seconds())
}

func (on *Gain) LinearRampToValueAtTime(v float64, t time.Duration) {
	on.val.Get("gain").Call("linearRampToValueAtTime", v, t.Seconds())
}
