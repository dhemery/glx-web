package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type Media struct {
	*glx.Media
	archive    *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	Date       glxdate.Date
	Type       VocabularyValue
	properties map[string]Property
}

func newMedia(id string, gm *glx.Media, archive *Archive) *Media {
	entityType := glx.EntityTypeMedia
	return &Media{
		archive:    archive,
		Media:      gm,
		EntityType: entityType,
		ID:         id,
		Slug:       path.Join(entityType.Plural(), id),
		Date:       newDate(gm.Date.String()),
		Type:       VocabularyValue{Value: gm.Type, Definition: archive.g.MediaTypes[gm.Type]},
	}
}

func (m *Media) Properties() map[string]Property {
	if m.properties == nil {
		m.properties = newProperties(m.Media.Properties, m.archive.g.MediaProperties, m.archive)
	}
	return m.properties
}

func (m *Media) Source() *Source {
	return m.archive.Sources[m.Media.Source]
}

func (m *Media) String() string {
	return m.Title
}
