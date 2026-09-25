package archive

import "github.com/genealogix/glx/go-glx"

type Citation struct {
	a          *Archive
	g          *glx.Citation
	ID         string
	properties map[string]Property
}

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

func (c *Citation) Properties() map[string]Property {
	if c.properties == nil {
		c.properties = newProperties(c.g.Properties, c.a.g.CitationProperties, c.a)
	}
	return c.properties
}
