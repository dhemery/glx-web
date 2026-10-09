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
		Date:       parseDate(inner.Date.String()),
	}
}

func (a *Assertion) compile(ctx *context) {
	inner := a.Assertion
	a.Citations = ctx.C.citations(inner.Citations)
	a.Media = ctx.C.media(inner.Media)
	a.Sources = ctx.C.sources(inner.Sources)

	s := inner.Subject
	switch {
	case s.Person != "":
		a.Subject = ctx.C.PersonsByID[s.Person]
	case s.Event != "":
		a.Subject = ctx.C.EventsByID[s.Event]
		roles := ctx.G.ParticipantRoles
		props := ctx.G.EventProperties
		a.Participant = newParticipant(ctx, inner.Participant, roles, props)
	case s.Relationship != "":
		a.Subject = ctx.C.RelationshipsByID[s.Relationship]
		roles := ctx.G.ParticipantRoles
		props := ctx.G.RelationshipProperties
		a.Participant = newParticipant(ctx, inner.Participant, roles, props)
	case s.Place != "":
		a.Subject = ctx.C.PlacesByID[s.Place]
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
