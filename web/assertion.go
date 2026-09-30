package web

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Assertion struct {
	*glx.Assertion
	entity
	Citations   []*Citation
	Date        glxdate.Date
	Media       []*Media
	Participant *Participant
	Sources     []*Source
	Subject     entityReference
	// TODO: Resolve value if the property has a reference type or vocabulary type
}

func newAssertion(id string, inner *glx.Assertion) *Assertion {
	return &Assertion{
		Assertion:  inner,
		EntityType: glx.EntityTypeAssertions,
		ID:         id,
		Date:       newDate(inner.Date.String()),
	}
}

func (a *Assertion) compile(archive *Archive, g *glx.GLXFile) {
	inner := a.Assertion
	a.Citations = archive.citationsWithIDs(inner.Citations)
	a.Media = archive.mediaWithIDs(inner.Media)
	a.Sources = archive.sourcesWithIDs(inner.Sources)

	s := inner.Subject
	switch {
	case s.Person != "":
		a.Subject = archive.Persons[s.Person]
	case s.Event != "":
		a.Subject = archive.Events[s.Event]
		a.Participant = newParticipant(inner.Participant, g.ParticipantRoles, g.EventProperties, archive, g)
	case s.Relationship != "":
		a.Subject = archive.Relationships[s.Relationship]
		a.Participant = newParticipant(inner.Participant, g.ParticipantRoles, g.RelationshipProperties, archive, g)
	case s.Place != "":
		a.Subject = archive.Places[s.Place]
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

func (a *Assertion) String() string {
	// TODO: Better String()
	return "Assertion " + a.ID
}
