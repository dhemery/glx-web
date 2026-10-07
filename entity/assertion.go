package entity

import (
	"log/slog"

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
	return sortValues(l)
}

func newAssertion(id string, inner *glx.Assertion) *Assertion {
	return &Assertion{
		Assertion:  inner,
		EntityType: glx.EntityTypeAssertions,
		id:         id,
		Date:       newDate(inner.Date.String(), slog.Default().With("new", id)),
	}
}

func (a *Assertion) resolve(arch archive, l *slog.Logger) {
	l = l.With("entity", a.PagePath())
	inner := a.Assertion
	a.Citations = arch.citations(inner.Citations)
	a.Media = arch.media(inner.Media)
	a.Sources = arch.sources(inner.Sources)

	s := inner.Subject
	switch {
	case s.Person != "":
		a.Subject = arch.c.PersonsByID[s.Person]
	case s.Event != "":
		a.Subject = arch.c.EventsByID[s.Event]
		roles := arch.g.ParticipantRoles
		props := arch.g.EventProperties
		a.Participant = newParticipant(arch, inner.Participant, roles, props, l)
	case s.Relationship != "":
		a.Subject = arch.c.RelationshipsByID[s.Relationship]
		roles := arch.g.ParticipantRoles
		props := arch.g.RelationshipProperties
		a.Participant = newParticipant(arch, inner.Participant, roles, props, l)
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

type entityReference interface {
	PagePath() string
	String() string
}
