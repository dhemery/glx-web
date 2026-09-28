package web

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Event struct {
	*glx.Event
	entity
	Date         glxdate.Date
	Participants []*Participant
	Place        *Place
	Properties   map[string]Property
	Type         VocabularyValue
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
}

func (e *Event) String() string {
	return e.Title
}
