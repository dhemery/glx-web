package entity

import (
	"github.com/genealogix/glx/go-glx"
)

// Repository represents a GLX repository entity.
type Repository struct {
	*glx.Repository
	id string

	Properties map[string]*Property
	Type       RepositoryType

	// Citations that directly reference this repository. Citations that
	// reference this repository only indirectly through sources are not
	// included.
	Citations CitationList
	// Sources in this repository.
	Sources SourceList
}

// EntityType returns [glx.EntityTypeRepositories].
func (r *Repository) EntityType() glx.EntityType {
	return glx.EntityTypeRepositories
}

// ID returns the repository's ID.
func (r *Repository) ID() string {
	return r.id
}

// String returns the repository's Name.
func (r *Repository) String() string {
	return r.Name
}

type RepositoryList []*Repository

// Sort returns the repositories sorted by String value.
func (l RepositoryList) Sort() RepositoryList {
	return sortedByString(l)
}

func newRepository(id string, inner *glx.Repository) *Repository {
	return &Repository{
		Repository: inner,
		id:         id,
	}
}

func (r *Repository) compile(ctx *context) {
	rctx := ctx.ForEntity(r)

	inner := r.Repository

	r.Properties = parseProperties(rctx, inner.Properties, rctx.GLX.RepositoryProperties)
	r.Type = RepositoryType(inner.Type)
}

func (r *Repository) addCitation(c *Citation) {
	r.Citations = append(r.Citations, c)
}

func (r *Repository) addSource(s *Source) {
	r.Sources = append(r.Sources, s)
}
