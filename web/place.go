package web

import (
	"github.com/genealogix/glx/go-glx"
)

type Place struct {
	entity
	*glx.Place
	EntityType glx.EntityType
	ID         string
	Parent     *Place
	Properties map[string]Property
	Type       VocabularyValue
}

func newPlace(id string, inner *glx.Place) *Place {
	return &Place{
		Place:      inner,
		EntityType: glx.EntityTypePlaces,
		ID:         id,
	}
}

func (p *Place) compile(a *Archive, g *glx.GLXFile) {
	inner := p.Place
	p.Parent = a.Places[p.ParentID]
	p.Properties = newProperties(inner.Properties, g.PlaceProperties, a, g)
	p.Type = VocabularyValue{Value: inner.Type, Definition: g.PlaceTypes[inner.Type]}
}

func (p *Place) FullName() string {
	return p.Name
}

func (p *Place) String() string {
	return p.FullName()
}
