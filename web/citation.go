package web

import "github.com/genealogix/glx/go-glx"

// Citation represents a GLX citation entity.
type Citation struct {
	a          *Archive
	g          *glx.Citation
	ID         string
	properties map[string]Property
}

// String returns a string representation of c formed by concatenating its
// locator (if it has one) onto its source's Title.
func (c *Citation) String() string {
	return "Citation " + c.ID
}

func newCitation(id string, gc *glx.Citation, a *Archive) *Citation {
	return &Citation{
		a:  a,
		g:  gc,
		ID: id,
	}
}

// Properties returns a map of c's properties indexed by name.
func (c *Citation) Properties() map[string]Property {
	if c.properties == nil {
		c.properties = newProperties(c.g.Properties, c.a.g.CitationProperties, c.a)
	}
	return c.properties
}
