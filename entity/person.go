package entity

import (
	"github.com/genealogix/glx/go-glx"
)

// Person represents a GLX person entity.
type Person struct {
	entity
	*glx.Person

	Properties map[string]Property

	// Assertions about this person.
	Assertions []*Assertion
	// Events this person participated in.
	Events []*Event
	// Relationships this person participated in.
	Relationships []*Relationship
}

// DisplayName returns a display name extracted from p's properties.
func (p *Person) DisplayName() string {
	return glx.PersonDisplayName(p.Person)
}

// String returns p's display name.
func (p *Person) String() string {
	return p.DisplayName()
}

func newPerson(id string, inner *glx.Person) *Person {
	return &Person{
		Person:     inner,
		EntityType: glx.EntityTypePersons,
		ID:         id,
	}
}

func (p *Person) compile(catalog *Catalog, glxfile *glx.GLXFile) {
	inner := p.Person
	p.Properties = newProperties(inner.Properties, glxfile.PersonProperties, catalog, glxfile)
}

func (p *Person) addAssertion(a *Assertion) {
	p.Assertions = append(p.Assertions, a)
}

func (p *Person) addEvent(e *Event) {
	p.Events = append(p.Events, e)
}

func (p *Person) addRelationship(e *Relationship) {
	p.Relationships = append(p.Relationships, e)
}
