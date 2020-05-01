package main

import (
	"bytes"
	"github.com/mariomac/msxmml/export/wasm"
	"github.com/mariomac/msxmml/lang"
	"syscall/js"
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
@1 <- o4 e8e8 r8 e8 r8 c8 e | g r < g r
      >c r8 <g r8 e | r8 a b b-8 a
      {g>eg}3 a f8 g8 | r8 e c8d8 <b

@2 <- o3 b8b8 r8 b8 r8 b8 b | >e r <e r
      e r8 c r8 < a | r8 >c d d-8 c
      {ca>c}3 d <b8>c8 | r8 <a f8g8 e
`

	s, err := lang.Root(lang.NewTokenizer(bytes.NewReader([]byte(str))))
	panicIfErr(err)
	wasm.ExportWasm(s)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
