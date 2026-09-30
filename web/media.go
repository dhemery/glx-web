package web

import (
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Media struct {
	entity
	*glx.Media

	Date       glxdate.Date
	Properties map[string]Property
	Source     *Source
	Type       VocabularyValue

	// Assertions directly citing this media entity as evidence.
	Assertions []*Assertion
}

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

func (m *Media) compile(a *Archive, g *glx.GLXFile) {
	inner := m.Media
	m.Type = VocabularyValue{Value: inner.Type, Definition: g.MediaTypes[inner.Type]}
	m.Properties = newProperties(inner.Properties, g.MediaProperties, a, g)
	m.Source = a.Sources[inner.Source]
}

func (m *Media) addAssertion(a *Assertion) {
	m.Assertions = append(m.Assertions, a)
}
