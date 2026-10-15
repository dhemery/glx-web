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

	// TODO(feature): Find a way to resolve entity and vocabulary
	// references in Value. Maybe a single struct that captures every
	// possible assertion structure, like how glx.EntityRef has a field for
	// each type of subject, and glx.PropertyDefinition has a field for
	// each kind of value type.
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
		participant, err := parseParticipant(ctx, inner.Participant, roles, props)
		if err != nil {
			actx.Sub("Participant").Warnf("discarded participant: %s", err)
		}
		a.Participant = participant
	case s.Relationship != "":
		a.Subject = actx.Catalog.RelationshipsByID[s.Relationship]
		roles := actx.GLX.ParticipantRoles
		props := actx.GLX.RelationshipProperties
		participant, err := parseParticipant(ctx, inner.Participant, roles, props)
		if err != nil {
			actx.Sub("Participant").Warnf("discarded participant: %s", err)
		}
		a.Participant = participant
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
	ID() string
	PagePath() string
	String() string
}
