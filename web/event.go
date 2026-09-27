package web

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Event struct {
	entity
	*glx.Event
	Date       glxdate.Date
	Place      *Place
	Properties map[string]Property
	Type       VocabularyValue
}

func newEvent(id string, event *glx.Event) *Event {
	return &Event{
		Event:      event,
		ID:         id,
		EntityType: glx.EntityTypeEvents,
		Date:       newDate(event.Date.String()),
	}
}

func (e *Event) compile(archive *Archive) {
	e.Place = archive.Places[e.PlaceID]
	e.Properties = newProperties(e.Event.Properties, archive.g.EventProperties, archive)
	e.Type = newVocabularyValue(e.Event.Type, archive.g.EventTypes)
}

func (e *Event) String() string {
	return e.Title
}
