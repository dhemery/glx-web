package web

import (
	"github.com/genealogix/glx/go-glx"
)

type Relationship struct {
	entity
	*glx.Relationship
	EndEvent   *Event
	Properties map[string]Property
	StartEvent *Event
	Type       VocabularyValue
	// TODO: Participants
}

func newRelationship(id string, gr *glx.Relationship) *Relationship {
	entityType := glx.EntityTypeRelationships
	return &Relationship{
		Relationship: gr,
		EntityType:   entityType,
		ID:           id,
	}
}

func (r *Relationship) compile(archive *Archive) {
	inner := r.Relationship
	r.EndEvent = archive.Events[inner.EndEvent]
	r.Properties = newProperties(inner.Properties, archive.g.RelationshipProperties, archive)
	r.StartEvent = archive.Events[inner.StartEvent]
	r.Type = newVocabularyValue(inner.Type, archive.g.RelationshipTypes)

}

func (r *Relationship) String() string {
	// TODO: Better String()
	return "Relationship " + r.ID
}
