package web

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Source struct {
	*glx.Source
	entity

	Date       glxdate.Date
	Media      []*Media
	Properties map[string]Property
	Repository *Repository
	Type       VocabularyValue

	// Assertions directly citing this source as evidence. Note that this
	// does not include assertions of this source made via citations.
	Assertions []*Assertion
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
}

func (s *Source) addAssertion(a *Assertion) {
	s.Assertions = append(s.Assertions, a)
}
