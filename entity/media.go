package entity

import (
	"log/slog"

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

type MediaList []*Media

// Sort returns the media sorted by String value.
func (l MediaList) Sort() MediaList {
	return sortValues(l)
}

func newMedia(id string, inner *glx.Media) *Media {
	return &Media{
		Media:      inner,
		EntityType: glx.EntityTypeMedia,
		id:         id,
		Date:       parseDate(inner.Date.String()),
	}
}

func (m *Media) resolve(a archive, l *slog.Logger) {
	l = l.With("entity_type", glx.EntityTypeMedia, "id", m.id)

	inner := m.Media
	m.Type = newVocabularyValue(inner.Type, a.g.MediaTypes, l)
	m.Properties = parseProperties(a, inner.Properties, a.g.MediaProperties, l)
	m.Source = a.c.SourcesByID[inner.Source]

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
