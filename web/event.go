package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Event struct {
	*glx.Event
	archive    *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	Date       glxdate.Date
	Type       VocabularyValue
	properties map[string]Property
}

func newEvent(id string, ge *glx.Event, archive *Archive) *Event {
	entityType := glx.EntityTypeEvents
	return &Event{
		archive:    archive,
		Event:      ge,
		EntityType: entityType,
		ID:         id,
		Slug:       path.Join(entityType.Plural(), id),
		Date:       newDate(ge.Date.String()),
		Type:       newVocabularyValue(ge.Type, archive.g.EventTypes),
	}
}

func (e *Event) Place() *Place {
	return e.archive.Places[e.PlaceID]
}

func (e *Event) Properties() map[string]Property {
	if e.properties == nil {
		e.properties = newProperties(e.Event.Properties, e.archive.g.EventProperties, e.archive)
	}
	return e.properties
}

func (e *Event) String() string {
	return e.Title
}
