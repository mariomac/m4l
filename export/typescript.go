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
	sixteenths := uint(0) // todo: consider higher?
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
			if _, err := fmt.Fprintf(out, `%d","time":"%d:%d:%d","velocity":1},%c`,
				note.Octave,
				sixteenths/16,
				(sixteenths/4)%4,
				sixteenths%4, '\n'); err != nil {
				return err
			}
		}
		sixteenths += 16 / note.Length
	}
	_, err = fmt.Fprintln(out, "];")
	return err
}
