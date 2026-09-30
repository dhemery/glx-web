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

	// Assertions that directly cite this source as evidence. Assertions
	// that cite this source only indirectly via citations or media are not
	// included.
	Assertions []*Assertion
	// Citations of this source.
	Citations []*Citation
	// Media that reference this source.
	ReferencingMedia []*Media
}

func (s *Source) String() string {
	return s.Title
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

	for _, m := range s.Media {
		m.addSource(s)
	}
	if s.Repository != nil {
		s.Repository.addSource(s)
	}
}

func (s *Source) addAssertion(a *Assertion) {
	s.Assertions = append(s.Assertions, a)
}

func (s *Source) addCitation(c *Citation) {
	s.Citations = append(s.Citations, c)
}

func (s *Source) addMedia(m *Media) {
	s.ReferencingMedia = append(s.ReferencingMedia, m)
}
