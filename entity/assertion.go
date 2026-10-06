package entity

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Assertion struct {
	entity
	*glx.Assertion

	Citations   []*Citation
	Date        glxdate.Date
	Media       []*Media
	Participant *Participant
	Sources     []*Source
	Subject     entityReference

	// TODO(dale): Resolve value if the property has a reference type or vocabulary type
}

// String returns a string representation of a.
// The value returned by the current implementation
// is useful only to identify which assertion produced it.
func (a *Assertion) String() string {
	// TODO: Better String()
	return "Assertion " + a.ID()
}

type AssertionList []*Assertion

// Sort returns the assertions sorted by String value.
func (l AssertionList) Sort() AssertionList {
	return sortStringers(l)
}

func newAssertion(id string, inner *glx.Assertion) *Assertion {
	return &Assertion{
		Assertion:  inner,
		EntityType: glx.EntityTypeAssertions,
		id:         id,
		Date:       newDate(inner.Date.String()),
	}
}

func (a *Assertion) compile(catalog *Catalog, glxfile *glx.GLXFile) {
	inner := a.Assertion
	a.Citations = catalog.citationsWithIDs(inner.Citations)
	a.Media = catalog.mediaWithIDs(inner.Media)
	a.Sources = catalog.sourcesWithIDs(inner.Sources)

	s := inner.Subject
	switch {
	case s.Person != "":
		a.Subject = catalog.PersonsByID[s.Person]
	case s.Event != "":
		a.Subject = catalog.EventsByID[s.Event]
		roles := glxfile.ParticipantRoles
		props := glxfile.EventProperties
		a.Participant = newParticipant(inner.Participant, roles, props, catalog, glxfile)
	case s.Relationship != "":
		a.Subject = catalog.RelationshipsByID[s.Relationship]
		roles := glxfile.ParticipantRoles
		props := glxfile.RelationshipProperties
		a.Participant = newParticipant(inner.Participant, roles, props, catalog, glxfile)
	case s.Place != "":
		a.Subject = catalog.PlacesByID[s.Place]
	}

	if s, ok := a.Subject.(subject); ok {
		s.addAssertion(a)
	}

	for _, c := range a.Citations {
		c.addAssertion(a)
	}

	for _, m := range a.Media {
		m.addAssertion(a)
	}

	for _, s := range a.Sources {
		s.addAssertion(a)
	}
}

type subject interface {
	addAssertion(*Assertion)
}
