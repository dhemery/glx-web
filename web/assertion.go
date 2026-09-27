package web

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Assertion struct {
	entity
	*glx.Assertion
	Citations []*Citation
	Date      glxdate.Date
	Media     []*Media
	Sources   []*Source
	Subject   fmt.Stringer
	// TODO: Participants
}

func newAssertion(id string, assertion *glx.Assertion) *Assertion {
	return &Assertion{
		ID:         id,
		EntityType: glx.EntityTypeAssertions,
		Assertion:  assertion,
		Date:       newDate(assertion.Date.String()),
	}
}

func (a *Assertion) compile(archive *Archive) {
	a.Citations = archive.citationsWithIDs(a.Assertion.Citations)
	a.Media = archive.mediaWithIDs(a.Assertion.Media)
	a.Sources = archive.sourcesWithIDs(a.Assertion.Sources)

	s := a.Assertion.Subject
	switch {
	case s.Person != "":
		a.Subject = archive.Persons[s.Person]
	case s.Event != "":
		a.Subject = archive.Events[s.Event]
	case s.Relationship != "":
		a.Subject = archive.Relationships[s.Relationship]
	case s.Place != "":
		a.Subject = archive.Places[s.Place]
	}
}

func (a *Assertion) String() string {
	// TODO: Better String()
	return "Assertion " + a.ID
}
