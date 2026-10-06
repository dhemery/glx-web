package entity

import (
	"github.com/genealogix/glx/go-glx"
)

type Place struct {
	entity
	*glx.Place

	Parent     *Place
	Properties map[string]Property
	Type       VocabularyValue

	// Assertions about this place.
	Assertions AssertionList
	// Places that identify this place as their parent.
	ChildPlaces PlaceList
	// Events identified as happening specifically in this place. This does
	// not include events in the "child" places of this place.
	Events EventList
}

// FullName returns p's Name followed by the FullName of its parent Place, if
// any, separated by ", ".
func (p *Place) FullName() string {
	if p.Parent == nil {
		return p.Name
	}
	return p.Name + ", " + p.Parent.FullName()
}

// String returns p's full name.
func (p *Place) String() string {
	return p.FullName()
}

type PlaceList []*Place

// Sort returns the places sorted by String value.
func (l PlaceList) Sort() PlaceList {
	return sortStringers(l)
}

func newPlace(id string, inner *glx.Place) *Place {
	return &Place{
		Place:      inner,
		EntityType: glx.EntityTypePlaces,
		id:         id,
	}
}

func (p *Place) compile(catalog *Catalog, glxfile *glx.GLXFile) {
	inner := p.Place
	p.Parent = catalog.PlacesByID[p.ParentID]
	p.Properties = newProperties(inner.Properties, glxfile.PlaceProperties, catalog, glxfile)
	p.Type = VocabularyValue{Value: inner.Type, Definition: glxfile.PlaceTypes[inner.Type]}

	if p.Parent != nil {
		p.Parent.addChild(p)
	}

}

func (p *Place) addAssertion(a *Assertion) {
	p.Assertions = append(p.Assertions, a)
}

func (p *Place) addChild(child *Place) {
	p.ChildPlaces = append(p.ChildPlaces, child)
}

func (p *Place) addEvent(e *Event) {
	p.Events = append(p.Events, e)
}
