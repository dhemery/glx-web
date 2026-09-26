package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

type Relationship struct {
	a          *Archive
	g          *glx.Relationship
	Type       VocabularyValue
	ID         string
	properties map[string]Property
}

func newRelationship(id string, gr *glx.Relationship, a *Archive) *Relationship {
	return &Relationship{
		a:    a,
		g:    gr,
		ID:   id,
		Type: VocabularyValue{Value: gr.Type, Definition: a.g.RelationshipTypes[gr.Type]},
	}
}

func (r *Relationship) EntityType() glx.EntityType {
	return glx.EntityTypeRelationships
}

func (r *Relationship) Slug() string {
	return path.Join(r.EntityType().Plural(), r.ID)
}
