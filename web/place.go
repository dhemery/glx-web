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

func newPlace(id string, gp *glx.Place) *Place {
	entityType := glx.EntityTypePlaces
	return &Place{
		Place:      gp,
		EntityType: entityType,
		ID:         id,
	}
}

func (p *Place) compile(a *Archive) {
	inner := p.Place
	p.Parent = a.Places[p.ParentID]
	p.Properties = newProperties(inner.Properties, a.g.PlaceProperties, a)
	p.Type = VocabularyValue{Value: inner.Type, Definition: a.g.PlaceTypes[inner.Type]}
}

func (p *Place) FullName() string {
	return p.Name
}

func (p *Place) String() string {
	return p.FullName()
}
