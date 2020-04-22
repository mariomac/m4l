package export

import (
	"fmt"
	"io"

	"github.com/mariomac/msxmml/song"

	"github.com/mariomac/msxmml/solfa"
)

type Exporter interface {
	Export(tab []byte, out io.Writer) error
}

type TypeScript struct{}

// TODO: use a go template
func (ts *TypeScript) Export(s *song.Song, out io.Writer) error {
	if _, err := fmt.Fprintln(out, "export var Song = {"); err != nil {
		return err
	}
	for _, c := range s.Channels {
		if err := exportChannel(c, out); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(out, "};")
	return err
}

func exportChannel(c song.Channel, out io.Writer) error {
	sixteenths := float64(0) // todo: consider higher?
	if _, err := fmt.Fprintf(out, "\t\"%s\":[\n", c.Name); err != nil {
		return err
	}
	for _, note := range c.Notes {
		if note.Pitch != solfa.Silence {
			fmt.Fprint(out, "\t\t")
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
	_, err := fmt.Fprintf(out, "\t],\n")
	return err
}
