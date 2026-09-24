package archive

import "github.com/genealogix/glx/go-glx"

type Person struct {
	g  *glx.Person
	ID string
}

func newPerson(id string, gp *glx.Person) *Person {
	return &Person{
		ID: id,
		g:  gp,
	}
}

func (p Person) DisplayName() string {
	return glx.PersonDisplayName(p.g)
}

func (p Person) Notes() glx.NoteList {
	return p.g.Notes
}
