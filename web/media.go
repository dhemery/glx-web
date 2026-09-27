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
}

func newMedia(id string, gm *glx.Media) *Media {
	entityType := glx.EntityTypeMedia
	return &Media{
		Media:      gm,
		EntityType: entityType,
		ID:         id,
		Date:       newDate(gm.Date.String()),
	}
}

func (m *Media) compile(archive *Archive) {
	inner := m.Media
	m.Type = VocabularyValue{Value: inner.Type, Definition: archive.g.MediaTypes[inner.Type]}
	m.Properties = newProperties(inner.Properties, archive.g.MediaProperties, archive)
	m.Source = archive.Sources[inner.Source]
}

func (m *Media) String() string {
	return m.Title
}
