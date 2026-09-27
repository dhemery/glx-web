package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

// Person represents a GLX person entity.
type Person struct {
	*glx.Person
	archive    *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	properties map[string]Property
}

func newPerson(id string, gp *glx.Person, archive *Archive) *Person {
	entityType := glx.EntityTypePersons
	return &Person{
		archive:    archive,
		Person:     gp,
		EntityType: entityType,
		ID:         id,
		Slug:       path.Join(entityType.Plural(), id),
	}
}

// DisplayName returns a display name extracted from the properties of p.
func (p *Person) DisplayName() string {
	return glx.PersonDisplayName(p.Person)
}

// Properties returns p's properties indexed by name.
func (p *Person) Properties() map[string]Property {
	if p.properties == nil {
		p.properties = newProperties(p.Person.Properties, p.archive.g.PersonProperties, p.archive)
	}
	return p.properties
}

// String returns the display name of p.
func (p *Person) String() string {
	return p.DisplayName()
}
