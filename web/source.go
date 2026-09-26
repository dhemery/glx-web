package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Source struct {
	a          *Archive
	g          *glx.Source
	ID         string
	Type       VocabularyValue
	properties map[string]Property
}

func newSource(id string, gs *glx.Source, a *Archive) *Source {
	return &Source{
		a:    a,
		g:    gs,
		ID:   id,
		Type: VocabularyValue{Value: gs.Type, Definition: a.g.SourceTypes[gs.Type]},
	}
}

func (s *Source) Authors() []string {
	return s.g.Authors
}

func (s *Source) Date() glxdate.Date {
	return newDate(s.g.Date.String())
}

func (s *Source) Language() string {
	return s.g.Language
}

func (s *Source) Media() []*Media {
	return s.a.mediaWithIDs(s.g.Media)
}

func (s *Source) Notes() glx.NoteList {
	return s.g.Notes
}

func (s *Source) Properties() map[string]Property {
	if s.properties == nil {
		s.properties = newProperties(s.g.Properties, s.a.g.SourceProperties, s.a)
	}
	return s.properties
}

func (s *Source) Repository() *Repository {
	return s.a.Repositories[s.g.RepositoryID]
}

func (s *Source) Title() string {
	return s.g.Title
}

func (s *Source) EntityType() glx.EntityType {
	return glx.EntityTypeSources
}

func (s *Source) Slug() string {
	return path.Join(s.EntityType().Plural(), s.ID)
}

func (s *Source) String() string {
	return s.Title()
}
