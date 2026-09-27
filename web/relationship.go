package web

import (
	"github.com/genealogix/glx/go-glx"
)

type Relationship struct {
	*glx.Relationship
	entity
	EndEvent   *Event
	Properties map[string]Property
	StartEvent *Event
	Type       VocabularyValue
	// TODO: Participantas
}

func newRelationship(id string, inner *glx.Relationship) *Relationship {
	return &Relationship{
		Relationship: inner,
		EntityType:   glx.EntityTypeRelationships,
		ID:           id,
	}
}

func (r *Relationship) compile(a *Archive, g *glx.GLXFile) {
	inner := r.Relationship
	r.EndEvent = a.Events[inner.EndEvent]
	r.Properties = newProperties(inner.Properties, g.RelationshipProperties, a, g)
	r.StartEvent = a.Events[inner.StartEvent]
	r.Type = newVocabularyValue(inner.Type, g.RelationshipTypes)

}

func (r *Relationship) String() string {
	// TODO: Better String()
	return "Relationship " + r.ID
}
