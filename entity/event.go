package entity

import (
	"slices"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Event struct {
	*glx.Event
	id string

	Date         glxdate.Date
	Participants ParticipantList
	Place        *Place
	Properties   map[string]*Property
	Type         *VocabularyValue

	// Assertions about this event.
	Assertions AssertionList
	// Relationships started or ended by this event.
	Relationships RelationshipList
}

// EntityType returns [glx.EntityTypeEvents].
func (e *Event) EntityType() glx.EntityType {
	return glx.EntityTypeEvents
}

// HasType reports whether the event's Type is present in types.
func (e *Event) HasType(types ...string) bool {
	return slices.Contains(types, e.Type.Value)
}

// ID returns the event's ID.
func (e *Event) ID() string {
	return e.id
}

// String returns the event's Title.
func (e *Event) String() string {
	return e.Title
}

type EventList []*Event

// Sort returns the events sorted by String value.
func (l EventList) Sort() EventList {
	return sortedByString(l)
}

// WithType returns the events in l whose Type is present in types.
func (l EventList) WithType(types ...string) EventList {
	var matches EventList
	for _, e := range l {
		if e.HasType(types...) {
			matches = append(matches, e)
		}
	}
	return matches
}

func cmpEventDate(x, y *Event) int {
	return CompareDates(x.Date, y.Date)
}

func (l EventList) ByDate() EventList {
	return slices.SortedStableFunc(slices.Values(l), cmpEventDate)
}

func newEvent(id string, inner *glx.Event) *Event {
	return &Event{
		Event: inner,
		id:    id,
		Date:  parseDate(inner.Date.String()),
	}
}

func (e *Event) compile(ctx *context) {
	ectx := ctx.ForEntity(e)

	inner := e.Event
	// TODO(validation): Warn if no participants. Spec requires at least one.
	e.Participants = parseParticipants(ectx, inner.Participants, ectx.GLX.EventProperties)
	e.Place = ectx.Catalog.PlacesByID[e.PlaceID]
	e.Properties = parseProperties(ectx, inner.Properties, ectx.GLX.EventProperties)
	// Ignoring error. GLX validation guarantees that the type is present and valid.
	e.Type, _ = newVocabularyValue(inner.Type, ectx.GLX.EventTypes)

	for _, p := range e.Participants {
		p.Person.addEvent(e)
	}

	if e.Place != nil {
		e.Place.addEvent(e)
	}
}

func (e *Event) addAssertion(a *Assertion) {
	e.Assertions = append(e.Assertions, a)
}

func (e *Event) addRelationship(r *Relationship) {
	e.Relationships = append(e.Relationships, r)
}
