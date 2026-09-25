package web

import "github.com/genealogix/glx/go-glx"

type Relationship struct {
	a  *Archive
	g  *glx.Relationship
	ID string
}

func newRelationship(id string, gr *glx.Relationship, a *Archive) *Relationship {
	return &Relationship{
		a:  a,
		g:  gr,
		ID: id,
	}
}
