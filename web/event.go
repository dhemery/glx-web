package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Event struct {
	a          *Archive
	g          *glx.Event
	ID         string
	Type       VocabularyValue
	properties map[string]Property
}

func newEvent(id string, ge *glx.Event, a *Archive) *Event {
	return &Event{
		a:    a,
		g:    ge,
		ID:   id,
		Type: newVocabularyValue(ge.Type, a.g.EventTypes),
	}
}

func (e *Event) Date() glxdate.Date {
	return newDate(e.g.Date.String())
}

func (e *Event) Notes() glx.NoteList {
	return e.g.Notes
}

func (e *Event) Path() string {
	return path.Join(glx.EntityTypeEvents.Plural(), e.ID)
}

func (e *Event) Place() *Place {
	return e.a.Places[e.g.PlaceID]
}

func (e *Event) Properties() map[string]Property {
	if e.properties == nil {
		e.properties = newProperties(e.g.Properties, e.a.g.EventProperties, e.a)
	}
	return e.properties
}

func (e *Event) Title() string {
	return e.g.Title
}
