package entity

import (
	"github.com/genealogix/glx/go-glx"
)

type Media struct {
	*glx.Media
	id string

	Properties map[string]*Property
	// The source documented by this media.
	Source *Source
	Type   MediaType

	// Assertions that directly cite this media as evidence. Assertions
	// that reference this media only indirectly through citations or
	// sources are not included.
	Assertions []*Assertion
	// Citations of this media.
	Citations []*Citation
	// Sources that reference this media.
	ReferencingSources []*Source
}

// EntityType returns [glx.EntityTypeMedia].
func (m *Media) EntityType() glx.EntityType {
	return glx.EntityTypeMedia
}

// ID returns the media's ID.
func (m *Media) ID() string {
	return m.id
}

// String returns the media's Title.
func (m *Media) String() string {
	return m.Title
}

type MediaList []*Media

func (l MediaList) EntityType() glx.EntityType {
	return glx.EntityTypeMedia
}

// Sort returns the media sorted by String value.
func (l MediaList) Sort() MediaList {
	return sortedByString(l)
}

func newMedia(id string, inner *glx.Media) *Media {
	return &Media{
		Media: inner,
		id:    id,
	}
}

func (m *Media) compile(ctx *context) {
	mctx := ctx.ForEntity(m)

	inner := m.Media
	// Ignore error. GLX validation guarantees that the type, if any, is valid.
	m.Type = MediaType(inner.Type)
	m.Properties = parseProperties(mctx, inner.Properties, mctx.GLX.MediaProperties)
	m.Source = mctx.Catalog.SourcesByID[inner.Source]

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
