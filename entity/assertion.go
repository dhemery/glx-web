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

	// TODO(feature): Resolve value if the property has a reference type or vocabulary type
}

// String returns a string representation of a.
// The value returned by the current implementation
// is useful only to identify which assertion produced it.
func (a *Assertion) String() string {
	// TODO(feature): Better String()
	return "Assertion " + a.ID()
}

type AssertionList []*Assertion

// Sort returns the assertions sorted by String value.
func (l AssertionList) Sort() AssertionList {
	return sortValues(l)
}

func newAssertion(id string, inner *glx.Assertion) *Assertion {
	return &Assertion{
		Assertion:  inner,
		EntityType: glx.EntityTypeAssertions,
		id:         id,
		Date:       parseDateString(inner.Date.String()),
	}
}

func (a *Assertion) compile(ctx *context) {
	actx := ctx.ForEntity(a.entity)

	inner := a.Assertion
	a.Citations = actx.Catalog.citations(inner.Citations)
	a.Media = actx.Catalog.media(inner.Media)
	a.Sources = actx.Catalog.sources(inner.Sources)

	s := inner.Subject
	switch {
	case s.Person != "":
		a.Subject = actx.Catalog.PersonsByID[s.Person]
	case s.Event != "":
		a.Subject = actx.Catalog.EventsByID[s.Event]
		roles := actx.GLX.ParticipantRoles
		props := actx.GLX.EventProperties
		a.Participant = newParticipant(actx.Sub("Participant"), inner.Participant, roles, props)
	case s.Relationship != "":
		a.Subject = actx.Catalog.RelationshipsByID[s.Relationship]
		roles := actx.GLX.ParticipantRoles
		props := actx.GLX.RelationshipProperties
		a.Participant = newParticipant(actx.Sub("Participant"), inner.Participant, roles, props)
	case s.Place != "":
		a.Subject = actx.Catalog.PlacesByID[s.Place]
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

type entityReference interface {
	PagePath() string
	String() string
}
