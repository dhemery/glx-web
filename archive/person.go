package archive

import "github.com/genealogix/glx/go-glx"

type Person struct {
	g *glx.Person
}

func newPerson(g *glx.Person) *Person {
	return new(Person{g: g})
}

func (p Person) Name() string {
	return glx.PersonDisplayName(p.g)
}

func (p Person) Notes() glx.NoteList {
	return p.g.Notes
}
