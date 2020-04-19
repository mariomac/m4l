package song

import "github.com/mariomac/msxmml/solfa"

type Song struct {
	Channels map[int]Channel
}

type Channel struct {
	Number int
	Notes []solfa.Note
}
