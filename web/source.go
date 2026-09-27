package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Source struct {
	*glx.Source
	archive    *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	Date       glxdate.Date
	Type       VocabularyValue
	properties map[string]Property
}

func newSource(id string, gs *glx.Source, archive *Archive) *Source {
	entityType := glx.EntityTypeSources
	return &Source{
		archive:    archive,
		Source:     gs,
		EntityType: entityType,
		ID:         id,
		Slug:       path.Join(entityType.Plural(), id),
		Date:       newDate(gs.Date.String()),
		Type:       VocabularyValue{Value: gs.Type, Definition: archive.g.SourceTypes[gs.Type]},
	}
}

func (s *Source) Media() []*Media {
	return s.archive.mediaWithIDs(s.Source.Media)
}

func (s *Source) Properties() map[string]Property {
	if s.properties == nil {
		s.properties = newProperties(s.Source.Properties, s.archive.g.SourceProperties, s.archive)
	}
	return s.properties
}

func (s *Source) Repository() *Repository {
	return s.archive.Repositories[s.RepositoryID]
}

func (s *Source) String() string {
	return s.Title
}
