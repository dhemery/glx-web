package entity

import (
	"strings"

	"github.com/genealogix/glx/go-glx"
)

type Assertion struct {
	*glx.Assertion
	id string

	// The Subject of the assertion, resolved to an Entity.
	Subject Entity

	// The label for the embedded assertion's property. Property is derived
	// from the field name or property definition if applicable. Otherwise
	// Property is a copy of the embedded assertion's Property.
	Property string

	// The value of the asserted property. When the embedded assertion's
	// Value corresponds to an entity or vocabulary term, Value is resolved
	// to the appropriate entity or VocabularyValue. Otherwise Value is
	// represented as a StringValue. The is_entity and is_vocabulary_value
	// template functions are helpful for determining how to interact with
	// this Value.
	Value Stringer

	Participant *Participant

	Citations CitationList
	Media     MediaList
	Sources   SourceList

	Confidence ConfidenceLevel
}

// EntityType returns [glx.EntityTypeAssertions].
func (a *Assertion) EntityType() glx.EntityType {
	return glx.EntityTypeAssertions
}

// ID returns the assertion's ID.
func (a *Assertion) ID() string {
	return a.id
}

// String returns a string representation of the assertion.
// The string is composed
// by concatenating the string values of
// the subject, property, participant, and assertion date,
// in that order,
// separated by spaces.
// If the participant is present,
// the role appears before the person's name.
// If the assertion date is present
// its string value is enclosed in parentheses.
func (a *Assertion) String() string {
	var parts []string
	parts = append(parts, a.Subject.String())
	if a.Property != "" {
		parts = append(parts, a.Property)
	}
	switch {
	case a.Participant != nil:
		// Valid participant.
		parts = append(parts, a.Participant.Role.Label(), a.Participant.Person.String())
	case a.Assertion.Participant != nil:
		// Discarded participant. Use details from underlying glx.Assertion as is.
		parts = append(parts, a.Assertion.Participant.Role, a.Assertion.Participant.Person)
	}
	if a.Date.String() != "" {
		parts = append(parts, "("+a.Date.String()+")")

	}
	return strings.Join(parts, " ")
}

type AssertionList []*Assertion

func (l AssertionList) EntityType() glx.EntityType {
	return glx.EntityTypeAssertions
}

// Sort returns the assertions sorted by String value.
func (l AssertionList) Sort() AssertionList {
	return sortedByString(l)
}

func newAssertion(id string, inner *glx.Assertion) *Assertion {
	return &Assertion{
		Assertion: inner,
		id:        id,
	}
}

func (a *Assertion) compile(ctx *context) {
	actx := ctx.ForEntity(a)

	inner := a.Assertion
	a.Confidence = ConfidenceLevel(inner.Confidence)

	a.Citations = actx.Catalog.citations(inner.Citations)
	a.Media = actx.Catalog.media(inner.Media)
	a.Sources = actx.Catalog.sources(inner.Sources)

	s := inner.Subject
	switch {
	case s.Person != "":
		a.Subject = actx.Catalog.PersonsByID[s.Person]
	case s.Event != "":
		a.Subject = actx.Catalog.EventsByID[s.Event]
		props := actx.GLX.EventProperties
		participant, err := parseParticipant(ctx, inner.Participant, props)
		if err != nil {
			actx.Sub("Participant").Warnf("discarded participant: %s", err)
		}
		a.Participant = participant
	case s.Relationship != "":
		a.Subject = actx.Catalog.RelationshipsByID[s.Relationship]
		props := actx.GLX.RelationshipProperties
		participant, err := parseParticipant(ctx, inner.Participant, props)
		if err != nil {
			actx.Sub("Participant").Warnf("discarded participant: %s", err)
		}
		a.Participant = participant
	case s.Place != "":
		a.Subject = actx.Catalog.PlacesByID[s.Place]
	}

	if s, ok := a.Subject.(subject); ok {
		resolver := s.addAssertion(a)
		a.Property = resolver.Label()
		a.Value = resolver.Value(ctx.Catalog, a.Assertion.Value)
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
	addAssertion(*Assertion) resolver
}

type resolver interface {
	Label() string
	Value(*Catalog, string) Stringer
}

type entityResolver struct {
	label    string
	typeName string
}

func (er entityResolver) Label() string {
	return er.label
}

func (er entityResolver) Value(catalog *Catalog, id string) Stringer {
	e := catalog.entity(id, er.typeName)
	if e == nil {
		return StringValue(id)
	}
	return e
}

type stringResolver string

func (sr stringResolver) Label() string {
	return string(sr)
}

func (sr stringResolver) Value(_ *Catalog, s string) Stringer {
	return StringValue(s)
}

type vocabularyResolver struct {
	label     string
	vocabName string
}

func (vr vocabularyResolver) Label() string {
	return vr.label
}

func (vr vocabularyResolver) Value(_ *Catalog, term string) Stringer {
	value, err := newVocabularyValue(vr.vocabName, term)
	if err != nil {
		return StringValue(term)
	}
	return value
}

var (
	_ subject = (*Event)(nil)
	_ subject = (*Person)(nil)
	_ subject = (*Place)(nil)
	_ subject = (*Relationship)(nil)
)
