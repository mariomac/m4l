package main

import (
	"github.com/mariomac/msxmml/export"
	"os"
)

func main() {
	// https://musescore.com/user/20360426/scores/4880846
	str := "a2 {r4a8}3{a8a8a8}3|{ag8}3 a {r4a8}3 {a8a8a8}3|{ag8}3 a {r4a8}3 {a8a8a8}3" +
		"a8e16e16e8e16e16e8e16e16e8e8|aer8a8a16b16>c#16d16" +
	    "e2"

	exp := export.TypeScript{}
	out, err := os.OpenFile("/Users/mmacias/code/mystuff/tonejs-experiments/src/song.ts",
		os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	if err := exp.Export([]byte(str), out); err != nil {
		panic(err)
	}
}
