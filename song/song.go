package song

import "github.com/mariomac/msxmml/solfa"

type Song struct {
	Channels map[string]Channel
}

type Channel struct {
	Name string
	Notes []solfa.Note
}
