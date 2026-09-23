package archive

import "github.com/genealogix/glx/go-glx"

type Place struct {
	g *glx.Place
}

func newPlace(gp *glx.Place) *Place {
	return &Place{
		g: gp,
	}
}

func (p *Place) Name() string {
	return p.g.Name
}

func (p *Place) Notes() glx.NoteList {
	return p.g.Notes
}

func (p *Place) Type() string {
	return p.g.Type
}

func (p *Place) ParentID() string {
	return p.g.ParentID
}

func (p *Place) Latitude() *float64 {
	return p.g.Latitude
}

func (p *Place) Longitude() *float64 {
	return p.g.Longitude
}
