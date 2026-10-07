package entity

import (
	"strings"

	"github.com/genealogix/glx/go-glx"
)

// Citation represents a GLX citation entity.
type Citation struct {
	entity
	*glx.Citation

	Media      MediaList
	Properties map[string]Property
	Repository *Repository
	Source     *Source

	// Assertions that cite this citation as evidence.
	Assertions AssertionList
}

// String returns a string formed by concatenating
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

type CitationList []*Citation

// Sort returns the citations sorted by String value.
func (l CitationList) Sort() CitationList {
	return sortStringers(l)
}

func newCitation(id string, inner *glx.Citation) *Citation {
	return &Citation{
		Citation:   inner,
		EntityType: glx.EntityTypeCitations,
		id:         id,
	}
}

func (c *Citation) compile(a archive) {
	inner := c.Citation
	c.Media = a.mediaWithIDs(inner.Media)
	c.Properties = newProperties(inner.Properties, a.g.CitationProperties, a)
	c.Repository = a.c.RepositoriesByID[c.RepositoryID]
	c.Source = a.c.SourcesByID[c.SourceID]

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
