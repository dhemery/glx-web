package entity

import (
	"fmt"

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

// String returns the display name
// composed by [glx.PersonDisplayName].
func (p *Person) String() string {
	return glx.PersonDisplayName(p.Person)
}

type PersonList []*Person

type PersonName struct {
	Prefix        string `json:"prefix"`
	Given         string `json:"given"`
	Nickname      string `json:"nickname"`
	SurnamePrefix string `json:"surname_prefix"`
	Surname       string `json:"surname"`
	Suffix        string `json:"suffix"`
}

// NewPersonName creates a name
// from parts supplied as successive key, value pairs.
// For each pair,
// the value is assigned to the name part identified by the key.
// Name parts are identified
// using the same keys as the "name" property
// in the standard property vocabularies for Person.
func NewPersonName(parts []string) *PersonName {
	return nil
}

// Format implements [fmt.Formatter].
// It is not useful for users to call directly.
func (n *PersonName) Format(f fmt.State, verb rune) {}

// Formatted formets n according to format.
func (n *PersonName) Formatted(format string) string {
	return ""
}

// WithDefaults returns a copy of n
// with its empty fields
// replaced by the corresponding fields of defaults.
func (n *PersonName) WithDefaults(defaults *PersonName) *PersonName {
	return nil
}

// Sort returns the persons sorted by String value.
func (l PersonList) Sort() PersonList {
	return sortStringers(l)
}

type PersonNameList []*PersonName

func (l PersonNameList) SortFormatted(format string) PersonNameList {
	return nil
}

func newPerson(id string, inner *glx.Person) *Person {
	return &Person{
		Person:     inner,
		EntityType: glx.EntityTypePersons,
		id:         id,
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
