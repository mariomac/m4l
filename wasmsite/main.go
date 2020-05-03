// +build wasm

package main

import (
	"bytes"
	"syscall/js"

	"github.com/mariomac/msxmml/export/wasm"
	"github.com/mariomac/msxmml/lang"
)

func main() {
	js.Global().Get("document").
		Call("getElementById", "play").
		Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			tocalaSam()
			return nil
		}))

	<-make(chan struct{})
}

func tocalaSam() {
	// https://musescore.com/user/20360426/scores/4880846
	str := `

@voice1 { wave: square }
@voice2 { wave: sine }
@voice3 { wave: triangle }
@voice4 { wave: sawtooth }
@voice5 {
	wave: square
	adsr: 5->100, 20->60, 25, 30
}

@voice1 <- o4 e8e8 r8 e8 r8 c8 e
@voice2 <- r1 | g r < g r
@voice3 <- r1 r1 | c r8 <g r8 e
@voice4 <- r1 r1 r1 | r8 a b b-8 a
@voice5 <- r1 r1 r1 r1 | {g>eg}3 a f8 g8 | r8 e c8d8 <b

`

	s, err := lang.Parse(lang.NewTokenizer(bytes.NewReader([]byte(str))))
	panicIfErr(err)
	wasm.ExportWasm(s)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
