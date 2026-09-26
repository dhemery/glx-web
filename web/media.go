package web

import "github.com/genealogix/glx/go-glx"

type Media struct {
	a          *Archive
	g          *glx.Media
	ID         string
	properties map[string]Property
}

func newMedia(id string, gm *glx.Media, a *Archive) *Media {
	return &Media{
		a:  a,
		g:  gm,
		ID: id,
	}
}

func (m *Media) Date() glx.DateString {
	return m.g.Date
}

func (m *Media) Hash() string {
	return m.g.Hash
}

func (m *Media) MimeType() string {
	return m.g.MimeType
}

func (m *Media) Notes() glx.NoteList {
	return m.g.Notes
}

func (m *Media) Properties() map[string]Property {
	if m.properties == nil {
		m.properties = newProperties(m.g.Properties, m.a.g.MediaProperties, m.a)
	}
	return m.properties
}

func (m *Media) Source() *Source {
	return m.a.Sources[m.g.Source]
}

func (m *Media) String() string {
	return m.Title()
}

func (m *Media) Title() string {
	return m.g.Title
}

func (m *Media) Type() string { // TODO: Type -> VocabularyValue
	return m.g.Type
}

func (m *Media) URI() string {
	return m.g.URI
}
