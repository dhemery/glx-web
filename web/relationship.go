package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

type Relationship struct {
	*glx.Relationship
	archive    *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	Type       VocabularyValue
	properties map[string]Property
	// TODO: Participants
}

func newRelationship(id string, gr *glx.Relationship, archive *Archive) *Relationship {
	entityType := glx.EntityTypeRelationships
	return &Relationship{
		Relationship: gr,
		archive:      archive,
		Type:         newVocabularyValue(gr.Type, archive.g.RelationshipTypes),
		EntityType:   entityType,
		ID:           id,
		Slug:         path.Join(entityType.Plural(), id),
	}
}

func (r *Relationship) EndEvent() *Event {
	return r.archive.Events[r.Relationship.EndEvent]
}

func (r *Relationship) StartEvent() *Event {
	return r.archive.Events[r.Relationship.StartEvent]
}

// Properties returns r's properties indexed by name.
func (r *Relationship) Properties() map[string]Property {
	if r.properties == nil {
		r.properties = newProperties(r.Relationship.Properties, r.archive.g.RelationshipProperties, r.archive)
	}
	return r.properties
}

func (r *Relationship) String() string {
	// TODO: Better String()
	return "Relationship " + r.ID
}
