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
	Citations CitationList
	// Sources in this repository.
	Sources SourceList
}

// String returns r's Name.
func (r *Repository) String() string {
	return r.Name
}

type RepositoryList []*Repository

// Sort returns the repositories sorted by String value.
func (l RepositoryList) Sort() RepositoryList {
	return sortValues(l)
}

func newRepository(id string, inner *glx.Repository) *Repository {
	return &Repository{
		Repository: inner,
		EntityType: glx.EntityTypeRepositories,
		id:         id,
	}
}

func (r *Repository) compile(ctx *context) {

	inner := r.Repository
	r.Type = newVocabularyValue(ctx, inner.Type, ctx.G.RepositoryTypes)
	r.Properties = parseProperties(ctx, inner.Properties, ctx.G.RepositoryProperties)
}

func (r *Repository) addCitation(c *Citation) {
	r.Citations = append(r.Citations, c)
}

func (r *Repository) addSource(s *Source) {
	r.Sources = append(r.Sources, s)
}
