package web

import (
	"strings"

	"github.com/genealogix/glx/go-glx"
)

// Citation represents a GLX citation entity.
type Citation struct {
	a          *Archive
	g          *glx.Citation
	ID         string
	properties map[string]Property
}

func newCitation(id string, gc *glx.Citation, a *Archive) *Citation {
	return &Citation{
		a:  a,
		g:  gc,
		ID: id,
	}
}

func (c *Citation) Media() []*Media {
	return c.a.mediaWithIDs(c.g.Media)
}

func (c *Citation) Notes() glx.NoteList {
	return c.g.Notes
}

// Properties returns a map of c's properties indexed by name.
func (c *Citation) Properties() map[string]Property {
	if c.properties == nil {
		c.properties = newProperties(c.g.Properties, c.a.g.CitationProperties, c.a)
	}
	return c.properties
}

func (c *Citation) Repository() *Repository {
	return c.a.Repositories[c.g.RepositoryID]
}

func (c *Citation) Source() *Source {
	return c.a.Sources[c.g.SourceID]
}

// String returns a string representation of c formed by concatenating its
// locator (if it has one) onto its source's Title.
func (c *Citation) String() string {
	parts := []string{
		c.Source().Title(),
	}

	if locator, ok := c.Properties()["locator"]; ok {
		parts = append(parts, locator.String())
	}

	return strings.Join(parts, ", ")
}
