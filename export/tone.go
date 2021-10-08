package export

import (
	"fmt"
	"github.com/mariomac/msxmml/pkg/song"
	"github.com/mariomac/msxmml/pkg/song/note"
	"io"
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

func exportChannel(c *song.Channel, out io.Writer) error {
	sixteenths := float64(0) // todo: consider higher?
	if _, err := fmt.Fprintf(out, "\t\"%s\":[\n", c.Name); err != nil {
		return err
	}
	for _, nt := range c.Notes {
		if nt.Pitch != note.Silence {
			fmt.Fprint(out, "\t\t")
			if _, err := fmt.Fprintf(out, `%c{"duration":"%dn","note":"%c`,
				'\t', nt.Length, nt.Pitch); err != nil {
				return err
			}
			switch nt.Halftone {
			case note.Sharp:
				if _, err := fmt.Fprintf(out, `%c`, '#'); err != nil {
					return err
				}
			case note.Flat:
				if _, err := fmt.Fprintf(out, `%c`, 'b'); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(out, `%d","time":%f,"velocity":1},%c`,
				nt.Octave,
				sixteenths/8, '\n'); err != nil {
				return err
			}
		}
		length := 16.0 / float64(nt.Length)
		if nt.Tuplet >= 3 { // todo consider 5-tuples etc...
			length *= float64(nt.Tuplet-1) / float64(nt.Tuplet)
		}
		sixteenths += length
	}
	_, err := fmt.Fprintf(out, "\t],\n")
	return err
}
