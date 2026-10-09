package entity

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
	Assertions AssertionList
	// Citations of this source.
	Citations CitationList
	// Media that reference this source.
	ReferencingMedia MediaList
}

// String returns s's Title.
func (s *Source) String() string {
	return s.Title
}

type SourceList []*Source

// Sort returns the sources sorted by String value.
func (l SourceList) Sort() SourceList {
	return sortValues(l)
}

func newSource(id string, inner *glx.Source) *Source {
	return &Source{
		Source:     inner,
		EntityType: glx.EntityTypeSources,
		id:         id,
		Date:       parseDate(inner.Date.String()),
	}
}

func (s *Source) compile(ctx *context) {
	inner := s.Source
	s.Media = ctx.C.media(inner.Media)
	s.Properties = parseProperties(ctx, inner.Properties, ctx.G.SourceProperties)
	s.Repository = ctx.C.RepositoriesByID[s.RepositoryID]
	s.Type = newVocabularyValue(ctx, inner.Type, ctx.G.SourceTypes)

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
