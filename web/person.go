package web

import "github.com/genealogix/glx/go-glx"

// Person represents a GLX person entity.
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

// DisplayName returns a display name extracted from the properties of p.
func (p *Person) DisplayName() string {
	return glx.PersonDisplayName(p.g)
}

// Notes returns the notes of p.
func (p *Person) Notes() glx.NoteList {
	return p.g.Notes
}

// Properties returns the p's properties indexed by name.
func (p *Person) Properties() map[string]Property {
	if p.properties == nil {
		p.properties = newProperties(p.g.Properties, p.a.g.PersonProperties, p.a)
	}
	return p.properties
}

// String returns the display name of p.
func (p *Person) String() string {
	return p.DisplayName()
}
