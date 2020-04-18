package export

import (
	"fmt"
	"io"

	"github.com/mariomac/msxmml/solfa"
)

type Exporter interface {
	Export(tab []byte, out io.Writer) error
}

type TypeScript struct{}

func (ts *TypeScript) Export(tab []byte, out io.Writer) error {
	notes, err := solfa.ParseChannel(tab)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, "export var Song = ["); err != nil {
		return err
	}
	sixteenths := float64(0) // todo: consider higher?
	for _, note := range notes {
		if note.Pitch != solfa.Silence {
			if _, err := fmt.Fprintf(out, `%c{"duration":"%dn","note":"%c`,
				'\t', note.Length, note.Pitch); err != nil {
				return err
			}
			if note.Halftone != solfa.NoHalftone {
				if _, err := fmt.Fprintf(out, `%c`, note.Halftone); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(out, `%d","time":%f,"velocity":1},%c`,
				note.Octave,
				sixteenths/8, '\n'); err != nil {
				return err
			}
		}
		length := 16.0 / float64(note.Length)
		if note.Tuplet >= 3 { // todo consider 5-tuples etc...
			length *= float64(note.Tuplet-1) / float64(note.Tuplet)
		}
		sixteenths += length
	}
	_, err = fmt.Fprintln(out, "];")
	return err
}
