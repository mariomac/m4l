package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var input, output string
	var help bool
	flag.StringVar(&input, "in", "", "input file")
	flag.StringVar(&output, "out", "", "output binary file for PSG")
	flag.BoolVar(&help, "h", false, "show help")
	fmt.Printf("in = %#v out = %#v help %#v\n", input, output, help)

	if input == "" || output == "" || help {
		flag.PrintDefaults()
		os.Exit(0)
	}

}
