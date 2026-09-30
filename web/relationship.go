package web

import (
	"github.com/genealogix/glx/go-glx"
)

type Relationship struct {
	entity
	*glx.Relationship
	EndEvent     *Event
	Participants []*Participant
	Properties   map[string]Property
	StartEvent   *Event
	Type         VocabularyValue

	// Assertions about this relationship.
	Assertions []*Assertion
}

func (r *Relationship) String() string {
	// TODO: Better String()
	return "Relationship " + r.ID
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
	r.Participants = newParticipants(inner.Participants, g.RelationshipProperties, a, g)
	r.Properties = newProperties(inner.Properties, g.RelationshipProperties, a, g)
	r.StartEvent = a.Events[inner.StartEvent]
	r.Type = newVocabularyValue(inner.Type, g.RelationshipTypes)

	for _, p := range r.Participants {
		p.Person.addRelationship(r)
	}
	if r.StartEvent != nil {
		r.StartEvent.addRelationship(r)
	}
	if r.EndEvent != nil {
		r.EndEvent.addRelationship(r)
	}
}

func (r *Relationship) addAssertion(a *Assertion) {
	r.Assertions = append(r.Assertions, a)
}
