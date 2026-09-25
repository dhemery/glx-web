package web

import (
	"github.com/genealogix/glx/go-glx"
)

type Assertion struct {
	a  *Archive
	g  *glx.Assertion
	ID string
}

func newAssertion(id string, ga *glx.Assertion, a *Archive) *Assertion {
	return &Assertion{
		a:  a,
		g:  ga,
		ID: id,
	}
}
