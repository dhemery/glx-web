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

func newPerson(id string, gp *glx.Person) *Person {
	entityType := glx.EntityTypePersons
	return &Person{
		Person:     gp,
		EntityType: entityType,
		ID:         id,
	}
}

func (p *Person) compile(archive *Archive) {
	inner := p.Person
	p.Properties = newProperties(inner.Properties, archive.g.PersonProperties, archive)
}

// DisplayName returns a display name extracted from p's properties.
func (p *Person) DisplayName() string {
	return glx.PersonDisplayName(p.Person)
}

// String returns the display name of p.
func (p *Person) String() string {
	return p.DisplayName()
}
