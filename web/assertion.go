package web

import (
	"fmt"
	"path"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Assertion struct {
	*glx.Assertion
	archive    *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	Date       glxdate.Date
}

func newAssertion(id string, ga *glx.Assertion, archive *Archive) *Assertion {
	entityType := glx.EntityTypeAssertions
	return &Assertion{
		Assertion:  ga,
		archive:    archive,
		Date:       newDate(ga.Date.String()),
		EntityType: entityType,
		ID:         id,
		Slug:       path.Join(entityType.Plural(), id),
	}
}

func (a *Assertion) Citations() []*Citation {
	return a.archive.citationsWithIDs(a.Assertion.Citations)
}

func (a *Assertion) Media() []*Media {
	return a.archive.mediaWithIDs(a.Assertion.Media)
}

func (a *Assertion) Sources() []*Source {
	return a.archive.sourcesWithIDs(a.Assertion.Sources)
}

func (a *Assertion) Subject() fmt.Stringer {
	s := a.Assertion.Subject
	switch {
	case s.Person != "":
		return a.archive.Persons[s.Person]
	case s.Event != "":
		return a.archive.Events[s.Event]
	case s.Relationship != "":
		return a.archive.Relationships[s.Relationship]
	case s.Place != "":
		return a.archive.Places[s.Place]
	default:
		return nil
	}
}

func (a *Assertion) String() string {
	// TODO: Better String()
	return "Assertion " + a.ID
}
