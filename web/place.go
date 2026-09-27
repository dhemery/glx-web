package web

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

type Place struct {
	*glx.Place
	a          *Archive
	EntityType glx.EntityType
	ID         string
	Slug       string
	Type       VocabularyValue
	properties map[string]Property
}

func newPlace(id string, gp *glx.Place, a *Archive) *Place {
	entityType := glx.EntityTypePlaces
	return &Place{
		Place:      gp,
		a:          a,
		EntityType: entityType,
		ID:         id,
		Slug:       path.Join(entityType.Plural(), id),
		Type:       VocabularyValue{Value: gp.Type, Definition: a.g.PlaceTypes[gp.Type]},
	}
}

func (p *Place) FullName() string {
	return p.Name
}

func (p *Place) Parent() *Place {
	return p.a.Places[p.ParentID]
}

func (p *Place) Path() string {
	return path.Join(glx.EntityTypePlaces.Plural(), p.ID)
}

func (p *Place) Properties() map[string]Property {
	if p.properties == nil {
		p.properties = newProperties(p.Place.Properties, p.a.g.PlaceProperties, p.a)
	}
	return p.properties
}

func (p *Place) String() string {
	return p.FullName()
}
