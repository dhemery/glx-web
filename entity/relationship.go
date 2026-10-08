package entity

import (
	"log/slog"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

type Relationship struct {
	entity
	*glx.Relationship

	EndEvent     *Event
	Participants ParticipantList
	Properties   map[string]Property
	StartEvent   *Event
	Type         VocabularyValue

	// Assertions about this relationship.
	Assertions AssertionList
}

// String returns a description of r
// composed by concatenating
// the string values of all participants,
// separated by commas and conjunctions as appropriate.
func (r *Relationship) String() string {
	var parts []string

	// TODO(feature): Include only principals?
	for _, p := range r.Participants {
		parts = append(parts, p.String())
	}
	switch len(parts) {
	case 0:
		return r.Type.Definition.Label + r.ID()
	case 1:
		return r.Type.Definition.Label + " of " + parts[0]
	case 2:
		return strings.Join(parts, " and ")
	default:
		parts[len(parts)-1] = "and " + parts[len(parts)-1]
		return strings.Join(parts, ", ")
	}
}

type RelationshipList []*Relationship

// Sort returns the relationships sorted by String value.
func (l RelationshipList) Sort() RelationshipList {
	return sortValues(l)
}

func newRelationship(id string, inner *glx.Relationship) *Relationship {
	return &Relationship{
		Relationship: inner,
		EntityType:   glx.EntityTypeRelationships,
		id:           id,
	}
}

func (r *Relationship) resolve(a archive, l *slog.Logger) {
	l = l.With("entity_type", glx.EntityTypeRelationships, "id", r.id)

	inner := r.Relationship
	r.EndEvent = a.c.EventsByID[inner.EndEvent]
	r.Participants = newParticipants(a, inner.Participants, a.g.RelationshipProperties, l)
	r.Properties = parseProperties(a, inner.Properties, a.g.RelationshipProperties, l)
	r.StartEvent = a.c.EventsByID[inner.StartEvent]
	r.Type = newVocabularyValue(inner.Type, a.g.RelationshipTypes, l)

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
