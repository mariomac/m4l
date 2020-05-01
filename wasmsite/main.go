package main

import (
	"syscall/js"
	"time"

	"github.com/mariomac/msxmml/export/wasm"
)

func main() {
	js.Global().Get("document").
		Call("getElementById", "play").
		Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ctx := wasm.WindowAudioContext()
			node := ctx.NoteNodes(440, wasm.ADSR{
				{1, 50 * time.Millisecond},
				{0.7, 100 * time.Millisecond},
				{0.7, 200 * time.Millisecond},
				{0, 250 * time.Millisecond},
			})
			node.TriggerEnvelope(1 * time.Second)
			node.TriggerEnvelope(2 * time.Second)
			node.TriggerEnvelope(3 * time.Second)
			node.TriggerEnvelope(5 * time.Second)
			return nil
		}))

	<-make(chan struct{})
}
