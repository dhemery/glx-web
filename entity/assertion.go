package entity

import (
	"strings"

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
	Subject     Stringer

	// TODO(feature): Find a way to resolve entity and vocabulary
	// references in Value. Maybe a single struct that captures every
	// possible assertion structure, like how glx.EntityRef has a field for
	// each type of subject, and glx.PropertyDefinition has a field for
	// each kind of value type.
}

// String returns a string representation of a
// composed by concatenating the string values of
// the subject, property, participant, and assertion date,
// in that order,
// separated by spaces.
// If the participant is present,
// the role appears before the person's name.
// If the assertion date is present
// its string value is prepended with "as of".
func (a *Assertion) String() string {
	var parts []string
	parts = append(parts, a.Subject.String())
	if a.Property != "" {
		parts = append(parts, a.Property)
	}
	switch {
	case a.Participant != nil:
		// Valid participant.
		parts = append(parts, a.Participant.Role.Value, a.Participant.Person.String())
	case a.Assertion.Participant != nil:
		// Discarded participant. Use details from underlying glx.Assertion as is.
		parts = append(parts, a.Assertion.Participant.Role, a.Assertion.Participant.Person)
	}
	if a.Date.Valid() {
		parts = append(parts, "as of", a.Date.String())

	}
	return strings.Join(parts, " ")
}

type AssertionList []*Assertion

// Sort returns the assertions sorted by String value.
func (l AssertionList) Sort() AssertionList {
	return sortedByString(l)
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
