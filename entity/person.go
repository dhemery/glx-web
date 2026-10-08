package entity

import (
	"log/slog"

	"github.com/genealogix/glx/go-glx"
)

// Person represents a GLX person entity.
type Person struct {
	entity
	*glx.Person

	Properties map[string]Property

	// Assertions about this person.
	Assertions AssertionList
	// Events this person participated in.
	Events EventList
	// Relationships this person participated in.
	Relationships RelationshipList
}

// DisplayName returns a display name extracted from p's properties.
func (p *Person) DisplayName() string {
	return glx.PersonDisplayName(p.Person)
}

// String returns p's display name.
func (p *Person) String() string {
	return p.DisplayName()
}

type PersonList []*Person

// Sort returns the persons sorted by String value.
func (l PersonList) Sort() PersonList {
	return sortValues(l)
}

func extractSurname(props map[string]Property) string {
	nameProp, ok := props["name"]
	if !ok {
		return ""
	}
	surnameField, ok := nameProp.Value().Fields["surname"]
	if !ok {
		return ""
	}
	return surnameField.String()
}

func (l PersonList) GroupBySurname() map[string]PersonList {
	grouped := make(map[string]PersonList)

	for _, p := range l {
		surname := extractSurname(p.Properties)
		grouped[surname] = append(grouped[surname], p)
	}
	return grouped
}

func newPerson(id string, inner *glx.Person) *Person {
	return &Person{
		Person:     inner,
		EntityType: glx.EntityTypePersons,
		id:         id,
	}
}

func (p *Person) resolve(a archive, l *slog.Logger) {
	l = l.With("entity_type", glx.EntityTypePersons, "id", p.id)

	inner := p.Person
	p.Properties = parseProperties(a, inner.Properties, a.g.PersonProperties, l)
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
