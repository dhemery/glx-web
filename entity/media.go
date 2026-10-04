package entity

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Media struct {
	entity
	*glx.Media

	Date       glxdate.Date
	Properties map[string]Property
	// The source documented by this media.
	Source *Source
	Type   VocabularyValue

	// Assertions that directly cite this media as evidence. Assertions
	// that reference this media only indirectly through citations or
	// sources are not included.
	Assertions []*Assertion
	// Citations of this media.
	Citations []*Citation
	// Sources that reference this media.
	ReferencingSources []*Source
}

// String returns m's Title.
func (m *Media) String() string {
	return m.Title
}

func newMedia(id string, inner *glx.Media) *Media {
	return &Media{
		Media:      inner,
		EntityType: glx.EntityTypeMedia,
		ID:         id,
		Date:       newDate(inner.Date.String()),
	}
}

func (m *Media) compile(catalog *Catalog, glxfile *glx.GLXFile) {
	inner := m.Media
	m.Type = VocabularyValue{Value: inner.Type, Definition: glxfile.MediaTypes[inner.Type]}
	m.Properties = newProperties(inner.Properties, glxfile.MediaProperties, catalog, glxfile)
	m.Source = catalog.Sources[inner.Source]

	if m.Source != nil {
		m.Source.addMedia(m)
	}
}

func (m *Media) addAssertion(a *Assertion) {
	m.Assertions = append(m.Assertions, a)
}

func (m *Media) addCitation(c *Citation) {
	m.Citations = append(m.Citations, c)
}

func (m *Media) addSource(s *Source) {
	m.ReferencingSources = append(m.ReferencingSources, s)
}
