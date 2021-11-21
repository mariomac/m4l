package main

import (
	"flag"
	"fmt"
	"github.com/mariomac/msxmml/pkg/lang"
	"github.com/mariomac/msxmml/pkg/psg"
	"os"
)

func main() {
	var input, output string
	var help bool
	flag.StringVar(&input, "in", "", "input file")
	flag.StringVar(&output, "out", "", "output binary file for PSG")
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()
	if input == "" || output == "" || help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	in, err := os.Open(input)
	if err != nil {
		fmt.Printf("ERROR opening input file %q: %v\n", input, err)
		os.Exit(-1)
	}
	defer in.Close()
	song, err := lang.Parse(in)
	if err != nil {
		fmt.Printf("ERROR parsing file %q: %v\n", input, err)
		os.Exit(-1)
	}
	songBytes, err := psg.Export(song)
	if err != nil {
		fmt.Printf("ERROR exporting song: %v\n", err)
		os.Exit(-1)
	}
	if err :=os.WriteFile(output, songBytes, 0644); err != nil {
		fmt.Printf("ERROR writing file %q: %v\n", output, err)
		os.Exit(-1)
	}
}
