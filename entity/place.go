package entity

import (
	"iter"
	"slices"

	"github.com/genealogix/glx/go-glx"
)

type Place struct {
	*glx.Place
	id string

	Parent     *Place
	Properties map[string]*Property
	Type       PlaceType

	// Assertions about this place.
	Assertions AssertionList
	// Places that identify this place as their parent.
	ChildPlaces PlaceList
	// Events identified as happening specifically in this place. This does
	// not include events in the "child" places of this place.
	Events EventList
}

// EntityType returns [glx.EntityTypePlaces].
func (p *Place) EntityType() glx.EntityType {
	return glx.EntityTypePlaces
}

// FullName returns the place's Name
// followed by the FullName of its parent Place, if any,
// separated by ", ".
func (p *Place) FullName() string {
	if p.Parent == nil {
		return p.Name
	}
	return p.Name + ", " + p.Parent.FullName()
}

// ID returns the place's ID.
func (p *Place) ID() string {
	return p.id
}

// String returns the place's name.
func (p *Place) String() string {
	return p.Name
}

type PlaceList []*Place

func (l PlaceList) All() iter.Seq[*Place] {
	return slices.Values(l)
}

func (l PlaceList) EntityType() glx.EntityType {
	return glx.EntityTypePlaces
}

// Sort returns the places sorted by String value.
func (l PlaceList) Sort() PlaceList {
	return sortedByString(l)
}

func newPlace(id string, inner *glx.Place) *Place {
	return &Place{
		Place: inner,
		id:    id,
	}
}

func (p *Place) compile(ctx *context) {
	pctx := ctx.ForEntity(p)

	inner := p.Place
	p.Parent = pctx.Catalog.PlacesByID[p.ParentID]
	p.Properties = parseProperties(pctx, inner.Properties, pctx.GLX.PlaceProperties)
	p.Type = PlaceType(inner.Type)

	if p.Parent != nil {
		p.Parent.addChild(p)
	}
}

func (p *Place) addAssertion(a *Assertion) resolver {
	p.Assertions = append(p.Assertions, a)

	propName := a.Assertion.Property
	switch propName {
	case "name":
		return stringResolver("Name")
	case "parent":
		return entityResolver{"Place", glx.EntityTypePlaces.Plural()}
	default:
		return resolverForProperty(glx.PropPlaceProperties, propName)
	}
}

func (p *Place) addChild(child *Place) {
	p.ChildPlaces = append(p.ChildPlaces, child)
}

func (p *Place) addEvent(e *Event) {
	p.Events = append(p.Events, e)
}
