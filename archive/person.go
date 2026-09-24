package archive

import "github.com/genealogix/glx/go-glx"

type Person struct {
	a  *Archive
	g  *glx.Person
	ID string
}

func (p Person) DisplayName() string {
	return glx.PersonDisplayName(p.g)
}

func (p Person) Notes() glx.NoteList {
	return p.g.Notes
}
