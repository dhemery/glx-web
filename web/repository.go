package web

import (
	"github.com/genealogix/glx/go-glx"
)

// Repository represents a GLX repository entity.
type Repository struct {
	entity
	*glx.Repository
	Properties map[string]Property
	Type       VocabularyValue
}

func newRepository(id string, gr *glx.Repository) *Repository {
	entityType := glx.EntityTypeRepositories
	return &Repository{
		Repository: gr,
		EntityType: entityType,
		ID:         id,
	}
}

func (r *Repository) compile(archive *Archive) {
	inner := r.Repository
	r.Type = VocabularyValue{Value: inner.Type, Definition: archive.g.RepositoryTypes[inner.Type]}
	r.Properties = newProperties(inner.Properties, archive.g.RepositoryProperties, archive)
}

func (r *Repository) String() string {
	return r.Name
}
