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

	// Citations that directly reference this repository. Note that this
	// does not include citations that indirectly reference this repository
	// via sources.
	Citations []*Citation
}

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

func (r *Repository) compile(a *Archive, g *glx.GLXFile) {
	inner := r.Repository
	r.Type = VocabularyValue{Value: inner.Type, Definition: g.RepositoryTypes[inner.Type]}
	r.Properties = newProperties(inner.Properties, g.RepositoryProperties, a, g)
}

func (r *Repository) addCitation(c *Citation) {
	r.Citations = append(r.Citations, c)
}
