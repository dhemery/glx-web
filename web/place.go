package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

func newPlace(id string, gp *glx.Place, a *Archive) *Place {
	return &Place{
		a:  a,
		g:  gp,
		ID: id,
	}
}

type Place struct {
	a  *Archive
	g  *glx.Place
	ID string
}

func (p *Place) FullName() string {
	return p.Name()
}

func (p *Place) Latitude() *float64 {
	return p.g.Latitude
}

func (p *Place) Longitude() *float64 {
	return p.g.Longitude
}

func (p *Place) Name() string {
	return p.g.Name
}

func (p *Place) Notes() glx.NoteList {
	return p.g.Notes
}

func (p *Place) Parent() *Place {
	return p.a.Places[p.g.ParentID]
}

func (p *Place) Path() string {
	return path.Join(glx.EntityTypePlaces.Plural(), p.ID)
}

func (p *Place) String() string {
	return p.Name()
}

func (p *Place) Type() string {
	return p.g.Type
}
