package web

import "github.com/genealogix/glx/go-glx"

type Person struct {
	a          *Archive
	g          *glx.Person
	ID         string
	properties map[string]Property
}

func newPerson(id string, gp *glx.Person, a *Archive) *Person {
	return &Person{
		a:  a,
		g:  gp,
		ID: id,
	}
}

func (p *Person) DisplayName() string {
	return glx.PersonDisplayName(p.g)
}

func (p *Person) Notes() glx.NoteList {
	return p.g.Notes
}

func (p *Person) Properties() map[string]Property {
	if p.properties == nil {
		p.properties = newProperties(p.g.Properties, p.a.g.PersonProperties, p.a)
	}
	return p.properties
}

func (p *Person) String() string {
	return p.DisplayName()
}
