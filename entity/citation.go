package entity

import (
	"strings"

	"github.com/genealogix/glx/go-glx"
)

// Citation represents a GLX citation entity.
type Citation struct {
	*glx.Citation
	id string

	Media      MediaList
	Properties map[string]*Property
	Repository *Repository
	Source     *Source

	// Assertions that cite this citation as evidence.
	Assertions AssertionList
}

// EntityType returns [glx.EntityTypeCitations].
func (c *Citation) EntityType() glx.EntityType {
	return glx.EntityTypeCitations
}

// ID returns the citation's ID.
func (c *Citation) ID() string {
	return c.id
}

// String returns a string representation of the citation.
// The string is formed by concatenating
// the string value of c's locator property
// (if it has one)
// onto the string value of C's source,
// separated by a comma.
func (c *Citation) String() string {
	parts := []string{c.Source.Title}

	if locator, ok := c.Properties["locator"]; ok {
		parts = append(parts, locator.String())
	}

	return strings.Join(parts, ", ")
}

func newCitation(id string, inner *glx.Citation) *Citation {
	return &Citation{
		Citation: inner,
		id:       id,
	}
}

func (c *Citation) compile(ctx *context) {
	cctx := ctx.ForEntity(c)

	inner := c.Citation
	c.Media = cctx.Catalog.media(inner.Media)
	c.Properties = parseProperties(cctx, inner.Properties, cctx.GLX.CitationProperties)
	c.Repository = cctx.Catalog.RepositoriesByID[c.RepositoryID]
	c.Source = cctx.Catalog.SourcesByID[c.SourceID]

	for _, m := range c.Media {
		m.addCitation(c)
	}
	if c.Repository != nil {
		c.Repository.addCitation(c)
	}
	if c.Source != nil {
		c.Source.addCitation(c)
	}
}

func (c *Citation) addAssertion(a *Assertion) {
	c.Assertions = append(c.Assertions, a)
}
