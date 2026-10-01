package entity

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

	// Assertions that cite this citation as evidence.
	Assertions []*Assertion
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

func newCitation(id string, inner *glx.Citation) *Citation {
	return &Citation{
		Citation:   inner,
		EntityType: glx.EntityTypeCitations,
		ID:         id,
	}
}

func (c *Citation) compile(a *Catalog, g *glx.GLXFile) {
	inner := c.Citation
	c.Media = a.mediaWithIDs(inner.Media)
	c.Properties = newProperties(inner.Properties, g.CitationProperties, a, g)
	c.Repository = a.Repositories[c.RepositoryID]
	c.Source = a.Sources[c.SourceID]

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
