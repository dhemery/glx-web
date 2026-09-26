package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

func newRepository(id string, gr *glx.Repository, a *Archive) *Repository {
	return &Repository{
		a:    a,
		g:    gr,
		ID:   id,
		Type: VocabularyValue{Value: gr.Type, Definition: a.g.RepositoryTypes[gr.Type]},
	}
}

// Repository represents a GLX repository entity.
type Repository struct {
	a          *Archive
	g          *glx.Repository
	ID         string
	Type       VocabularyValue
	properties map[string]Property
}

func (r Repository) Address() string {
	return r.g.Address
}

func (r Repository) City() string {
	return r.g.City
}

func (r Repository) Country() string {
	return r.g.Country
}

func (r Repository) Name() string {
	return r.g.Name
}

func (r Repository) Notes() glx.NoteList {
	return r.g.Notes
}

func (r Repository) PostalCode() string {
	return r.g.PostalCode
}

func (r Repository) Properties() map[string]Property {
	if r.properties == nil {
		r.properties = newProperties(r.g.Properties, r.a.g.RepositoryProperties, r.a)
	}
	return r.properties
}

func (r Repository) State() string {
	return r.g.State
}

func (r Repository) Website() string {
	return r.g.Website
}

func (r *Repository) EntityType() glx.EntityType {
	return glx.EntityTypeRepositories
}

func (r *Repository) Slug() string {
	return path.Join(r.EntityType().Plural(), r.ID)
}

func (r *Repository) String() string {
	return r.Name()
}
