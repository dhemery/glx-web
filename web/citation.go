package web

import (
	"path"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

// Citation represents a GLX citation entity.
type Citation struct {
	*glx.Citation
	archive    *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	properties map[string]Property
}

func newCitation(id string, gc *glx.Citation, archive *Archive) *Citation {
	entityType := glx.EntityTypeCitations
	return &Citation{
		Citation:   gc,
		archive:    archive,
		EntityType: entityType,
		ID:         id,
		Slug:       path.Join(entityType.Plural(), id),
	}
}

func (c *Citation) Media() []*Media {
	return c.archive.mediaWithIDs(c.Citation.Media)
}

// Properties returns a map of c's properties indexed by name.
func (c *Citation) Properties() map[string]Property {
	if c.properties == nil {
		c.properties = newProperties(c.Citation.Properties, c.archive.g.CitationProperties, c.archive)
	}
	return c.properties
}

func (c *Citation) Repository() *Repository {
	return c.archive.Repositories[c.RepositoryID]
}

func (c *Citation) Source() *Source {
	return c.archive.Sources[c.SourceID]
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
