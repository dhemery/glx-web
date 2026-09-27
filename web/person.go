package web

import (
	"github.com/genealogix/glx/go-glx"
)

// Person represents a GLX person entity.
type Person struct {
	entity
	*glx.Person
	Properties map[string]Property
}

func newPerson(id string, inner *glx.Person) *Person {
	return &Person{
		Person:     inner,
		EntityType: glx.EntityTypePersons,
		ID:         id,
	}
}

func (p *Person) compile(a *Archive, g *glx.GLXFile) {
	inner := p.Person
	p.Properties = newProperties(inner.Properties, g.PersonProperties, a, g)
}

// DisplayName returns a display name extracted from p's properties.
func (p *Person) DisplayName() string {
	return glx.PersonDisplayName(p.Person)
}

// String returns the display name of p.
func (p *Person) String() string {
	return p.DisplayName()
}
