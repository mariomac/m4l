package main

import (
	"github.com/mariomac/msxmml/export"
	"github.com/mariomac/msxmml/lang"
	"os"
)

func main() {
	// https://musescore.com/user/20360426/scores/4880846
	str := `
@uno <- abcde
@dos <- O2a2c2d2
`


	exp := export.TypeScript{}
	out, err := os.OpenFile("/Users/mmacias/code/tonejs-experiments/src/song.ts",
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	panicIfErr(err)
	defer out.Close()

	s, err := lang.Root(lang.NewTokenizer([]byte(str)))
	panicIfErr(err)
	panicIfErr(exp.Export(s, out))
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
