package web

import (
	"path"

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

func (a *Assertion) EntityType() glx.EntityType {
	return glx.EntityTypeAssertions
}

func (a *Assertion) Slug() string {
	return path.Join(a.EntityType().Plural(), a.ID)
}

func (a *Assertion) String() string {
	// TODO: subject + details asserted
	return "Assertion " + a.ID

}
