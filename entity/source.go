package entity

import (
	"iter"
	"slices"

	"github.com/genealogix/glx/go-glx"
)

type Source struct {
	*glx.Source
	id string

	Media      []*Media
	Properties map[string]*Property
	Repository *Repository
	Type       SourceType

	// Assertions that directly cite this source as evidence. Assertions
	// that cite this source only indirectly via citations or media are not
	// included.
	Assertions AssertionList
	// Citations of this source.
	Citations CitationList
	// Media that reference this source.
	ReferencingMedia MediaList
}

// EntityType returns [glx.EntityTypeSources].
func (s *Source) EntityType() glx.EntityType {
	return glx.EntityTypeSources
}

// ID returns the source's ID.
func (s *Source) ID() string {
	return s.id
}

// String returns the source's Title.
func (s *Source) String() string {
	return s.Title
}

type SourceList []*Source

func (l SourceList) All() iter.Seq[*Source] {
	return slices.Values(l)
}

func (l SourceList) EntityType() glx.EntityType {
	return glx.EntityTypeSources
}

// Sort returns the sources sorted by String value.
func (l SourceList) Sort() SourceList {
	return sortedByString(l)
}

func newSource(id string, inner *glx.Source) *Source {
	return &Source{
		Source: inner,
		id:     id,
	}
}

func (s *Source) compile(ctx *context) {
	sctx := ctx.ForEntity(s)

	inner := s.Source
	s.Media = sctx.Catalog.media(inner.Media)
	s.Properties = parseProperties(sctx, inner.Properties, sctx.GLX.SourceProperties)
	s.Repository = sctx.Catalog.RepositoriesByID[s.RepositoryID]
	s.Type = SourceType(inner.Type)

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
