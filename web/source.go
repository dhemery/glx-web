package web

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Source struct {
	entity
	*glx.Source
	Date       glxdate.Date
	Media      []*Media
	Properties map[string]Property
	Repository *Repository
	Type       VocabularyValue
}

func newSource(id string, inner *glx.Source) *Source {
	return &Source{
		Source:     inner,
		EntityType: glx.EntityTypeSources,
		ID:         id,
		Date:       newDate(inner.Date.String()),
	}
}

func (s *Source) compile(a *Archive, g *glx.GLXFile) {
	inner := s.Source
	s.Media = a.mediaWithIDs(inner.Media)
	s.Properties = newProperties(inner.Properties, g.SourceProperties, a, g)
	s.Repository = a.Repositories[s.RepositoryID]
	s.Type = VocabularyValue{Value: inner.Type, Definition: g.SourceTypes[inner.Type]}
}

func (s *Source) String() string {
	return s.Title
}
