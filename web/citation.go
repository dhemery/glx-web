package web

import (
	"strings"

	"github.com/genealogix/glx/go-glx"
)

// Citation represents a GLX citation entity.
type Citation struct {
	entity
	*glx.Citation
	Media      []*Media
	Properties map[string]Property
	Repository *Repository
	Source     *Source
}

func newCitation(id string, gc *glx.Citation) *Citation {
	entityType := glx.EntityTypeCitations
	return &Citation{
		EntityType: entityType,
		ID:         id,
		Citation:   gc,
	}
}

func (c *Citation) compile(archive *Archive) {
	inner := c.Citation
	c.Media = archive.mediaWithIDs(inner.Media)
	c.Properties = newProperties(inner.Properties, archive.g.CitationProperties, archive)
	c.Repository = archive.Repositories[c.RepositoryID]
	c.Source = archive.Sources[c.SourceID]
}

// String returns a string representation of c formed by concatenating its
// locator (if it has one) onto its source's Title (if it has one).
func (c *Citation) String() string {
	parts := []string{c.Source.Title}

	if locator, ok := c.Properties["locator"]; ok {
		parts = append(parts, locator.String())
	}

	return strings.Join(parts, ", ")
}
