package web

import "github.com/genealogix/glx/go-glx"

type Source struct {
	a  *Archive
	g  *glx.Source
	ID string
}

func newSource(id string, gs *glx.Source, a *Archive) *Source {
	return &Source{
		a:  a,
		g:  gs,
		ID: id,
	}
}
