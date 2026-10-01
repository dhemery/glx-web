package entity

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Event struct {
	entity
	*glx.Event

	Date         glxdate.Date
	Participants []*Participant
	Place        *Place
	Properties   map[string]Property
	Type         VocabularyValue

	// Assertions about this event.
	Assertions []*Assertion
	// Relationships started or ended by this event.
	Relationships []*Relationship
}

func (e *Event) String() string {
	return e.Title
}

func newEvent(id string, inner *glx.Event) *Event {
	return &Event{
		Event:      inner,
		EntityType: glx.EntityTypeEvents,
		ID:         id,
		Date:       newDate(inner.Date.String()),
	}
}

func (e *Event) compile(a *Catalog, g *glx.GLXFile) {
	inner := e.Event
	e.Participants = newParticipants(inner.Participants, g.EventProperties, a, g)
	e.Place = a.Places[e.PlaceID]
	e.Properties = newProperties(inner.Properties, g.EventProperties, a, g)
	e.Type = newVocabularyValue(inner.Type, g.EventTypes)

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
