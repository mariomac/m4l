package wasm

import (
	"log"
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
	log.Printf("rampValue %f at %v", v, t)
	on.val.Get("gain").Call("linearRampToValueAtTime", v, t.Seconds())
}
