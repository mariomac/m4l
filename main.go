package main

import (
	"github.com/mariomac/msxmml/export"
	"os"
)

func main() {
	// https://musescore.com/user/20360426/scores/4880846
	str := `
a2{r4a8}3{a8a8a8}3 | {ag8}3 a {r4a8}3 {a8a8a8}3
{ag8}3 a {r4a8}3 {a8a8a8}3 | a8e16e16e8e16e16e8e16e16e8e8
aer8a8a16b16>c#16d16 | e2 r8e8{e8f8g9}3 | a2 r8a8{a8g8f8}3
g8r16f16e2e | d8d16e16f2 e8d8 | c8c16d16e2d8c8
<b8b16>c#16d#2f# | e8<e16e16 e8e16e16 e8e16e16 e8 e8
`


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
