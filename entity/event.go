package entity

import (
	"log/slog"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Event struct {
	entity
	*glx.Event

	Date         glxdate.Date
	Participants ParticipantList
	Place        *Place
	Properties   map[string]Property
	Type         VocabularyValue

	// Assertions about this event.
	Assertions AssertionList
	// Relationships started or ended by this event.
	Relationships RelationshipList
}

// String returns e's Title.
func (e *Event) String() string {
	return e.Title
}

type EventList []*Event

// Sort returns the events sorted by String value.
func (l EventList) Sort() EventList {
	return sortValues(l)
}

func newEvent(id string, inner *glx.Event) *Event {
	return &Event{
		Event:      inner,
		EntityType: glx.EntityTypeEvents,
		id:         id,
		Date:       newDate(inner.Date.String(), slog.Default().With("new", id)),
	}
}

func (e *Event) resolve(a archive, l *slog.Logger) {
	l = l.With("entity", e.PagePath())

	inner := e.Event
	e.Participants = newParticipants(a, inner.Participants, a.g.EventProperties, l)
	e.Place = a.c.PlacesByID[e.PlaceID]
	e.Properties = newProperties(a, inner.Properties, a.g.EventProperties, l)
	e.Type = newVocabularyValue(inner.Type, a.g.EventTypes, l)

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
