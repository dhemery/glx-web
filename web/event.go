package web

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Event struct {
	entity
	*glx.Event
	Assertions    []*Assertion // Assertions about this event.
	Date          glxdate.Date
	Participants  []*Participant
	Place         *Place
	Properties    map[string]Property
	Relationships []*Relationship // Relationships started or ended by this event.
	Type          VocabularyValue
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

func (e *Event) compile(a *Archive, g *glx.GLXFile) {
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
