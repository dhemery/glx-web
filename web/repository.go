package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

// Repository represents a GLX repository entity.
type Repository struct {
	*glx.Repository
	archive    *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	Type       VocabularyValue
	properties map[string]Property
}

func newRepository(id string, gr *glx.Repository, archive *Archive) *Repository {
	entityType := glx.EntityTypeRepositories
	return &Repository{
		archive:    archive,
		Repository: gr,
		EntityType: entityType,
		ID:         id,
		Slug:       path.Join(entityType.Plural(), id),
		Type:       VocabularyValue{Value: gr.Type, Definition: archive.g.RepositoryTypes[gr.Type]},
	}
}
func (r Repository) Properties() map[string]Property {
	if r.properties == nil {
		r.properties = newProperties(r.Repository.Properties, r.archive.g.RepositoryProperties, r.archive)
	}
	return r.properties
}

func (r *Repository) String() string {
	return r.Name
}
