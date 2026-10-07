package entity

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Assertion struct {
	entity
	*glx.Assertion

	Citations   CitationList
	Date        glxdate.Date
	Media       MediaList
	Participant *Participant
	Sources     SourceList
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

func (a *Assertion) compile(arch archive) {
	inner := a.Assertion
	a.Citations = arch.citationsWithIDs(inner.Citations)
	a.Media = arch.mediaWithIDs(inner.Media)
	a.Sources = arch.sourcesWithIDs(inner.Sources)

	s := inner.Subject
	switch {
	case s.Person != "":
		a.Subject = arch.c.PersonsByID[s.Person]
	case s.Event != "":
		a.Subject = arch.c.EventsByID[s.Event]
		roles := arch.g.ParticipantRoles
		props := arch.g.EventProperties
		a.Participant = newParticipant(inner.Participant, roles, props, arch)
	case s.Relationship != "":
		a.Subject = arch.c.RelationshipsByID[s.Relationship]
		roles := arch.g.ParticipantRoles
		props := arch.g.RelationshipProperties
		a.Participant = newParticipant(inner.Participant, roles, props, arch)
	case s.Place != "":
		a.Subject = arch.c.PlacesByID[s.Place]
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
