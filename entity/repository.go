package entity

import (
	"github.com/genealogix/glx/go-glx"
)

// Repository represents a GLX repository entity.
type Repository struct {
	entity
	*glx.Repository

	Properties map[string]Property
	Type       VocabularyValue

	// Citations that directly reference this repository. Citations that
	// reference this repository only indirectly through sources are not
	// included.
	Citations []*Citation
	// Sources in this repository.
	Sources []*Source
}

// String returns r's Name.
func (r *Repository) String() string {
	return r.Name
}

func newRepository(id string, inner *glx.Repository) *Repository {
	return &Repository{
		Repository: inner,
		EntityType: glx.EntityTypeRepositories,
		ID:         id,
	}
}

func (r *Repository) compile(catalog *Catalog, glxfile *glx.GLXFile) {
	inner := r.Repository
	r.Type = VocabularyValue{Value: inner.Type, Definition: glxfile.RepositoryTypes[inner.Type]}
	r.Properties = newProperties(inner.Properties, glxfile.RepositoryProperties, catalog, glxfile)
}

func (r *Repository) addCitation(c *Citation) {
	r.Citations = append(r.Citations, c)
}

func (r *Repository) addSource(s *Source) {
	r.Sources = append(r.Sources, s)
}
