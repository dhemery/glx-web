package entity

import (
	"iter"
	"slices"

	"github.com/genealogix/glx/go-glx"
)

// Person represents a GLX person entity.
type Person struct {
	*glx.Person
	id string

	Properties map[string]*Property

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

// EntityType returns [glx.EntityTypePersons].
func (p *Person) EntityType() glx.EntityType {
	return glx.EntityTypePersons
}

// ID returns the person's ID.
func (p *Person) ID() string {
	return p.id
}

// String returns the person's DisplayName.
func (p *Person) String() string {
	return p.DisplayName()
}

type PersonList []*Person

func (l PersonList) All() iter.Seq[*Person] {
	return slices.Values(l)
}

func (l PersonList) EntityType() glx.EntityType {
	return glx.EntityTypePersons
}

// Sort returns the persons sorted by String value.
func (l PersonList) Sort() PersonList {
	return sortedByString(l)
}

func extractSurname(props map[string]*Property) string {
	nameProp, ok := props["name"]
	if !ok {
		return ""
	}
	surnameField, ok := nameProp.Entry().Fields["surname"]
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
		Person: inner,
		id:     id,
	}
}

func (p *Person) compile(ctx *context) {
	pctx := ctx.ForEntity(p)

	inner := p.Person
	p.Properties = parseProperties(pctx, inner.Properties, pctx.GLX.PersonProperties)
}

func (p *Person) addAssertion(a *Assertion) resolver {
	p.Assertions = append(p.Assertions, a)

	return resolverForProperty(glx.PropPersonProperties, a.Assertion.Property)
}

func (p *Person) addEvent(e *Event) {
	if p == nil {
		return
	}
	p.Events = append(p.Events, e)
}

func (p *Person) addRelationship(e *Relationship) {
	p.Relationships = append(p.Relationships, e)
}
