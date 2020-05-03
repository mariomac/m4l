package song

import "github.com/mariomac/msxmml/note"

type Song struct {
	Channels map[string]*Channel
}

type Channel struct {
	Status struct {
		Octave int
	}
	Name string
	Notes []note.Note
}
