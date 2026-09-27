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

func newSource(id string, gs *glx.Source) *Source {
	entityType := glx.EntityTypeSources
	return &Source{
		Source:     gs,
		EntityType: entityType,
		ID:         id,
		Date:       newDate(gs.Date.String()),
	}
}

func (s *Source) compile(archive *Archive) {
	inner := s.Source
	s.Media = archive.mediaWithIDs(inner.Media)
	s.Properties = newProperties(inner.Properties, archive.g.SourceProperties, archive)
	s.Repository = archive.Repositories[s.RepositoryID]
	s.Type = VocabularyValue{Value: inner.Type, Definition: archive.g.SourceTypes[inner.Type]}
}

func (s *Source) String() string {
	return s.Title
}
