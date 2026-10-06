package entity

import (
	"strings"

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
	Assertions AssertionList
}

// String returns a description of r
// composed by concatenating
// the role and name of each participant,
// separated by commas and conjunctions as appropriate.
func (r *Relationship) String() string {
	var parts []string

	// TODO(dale): Maybe include only names of principals.
	for _, p := range r.Participants {
		role := p.Role.Definition.Label
		name := p.Person.DisplayName()
		parts = append(parts, role+" "+name)
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
	return sortStringers(l)
}

func newRelationship(id string, inner *glx.Relationship) *Relationship {
	return &Relationship{
		Relationship: inner,
		EntityType:   glx.EntityTypeRelationships,
		id:           id,
	}
}

func (r *Relationship) compile(catalog *Catalog, glxfile *glx.GLXFile) {
	inner := r.Relationship
	r.EndEvent = catalog.EventsByID[inner.EndEvent]
	r.Participants = newParticipants(inner.Participants, glxfile.RelationshipProperties, catalog, glxfile)
	r.Properties = newProperties(inner.Properties, glxfile.RelationshipProperties, catalog, glxfile)
	r.StartEvent = catalog.EventsByID[inner.StartEvent]
	r.Type = newVocabularyValue(inner.Type, glxfile.RelationshipTypes)

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
