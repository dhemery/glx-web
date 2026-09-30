package web

import (
	"github.com/genealogix/glx/go-glx"
)

type Place struct {
	*glx.Place
	entity

	// Whatever is this

	Parent     *Place
	Properties map[string]Property
	Type       VocabularyValue

	// The assertions about this place.
	Assertions []*Assertion
	// The places that identify this place as their parent.
	ChildPlaces []*Place
	// Tne events specified as happening in this place. Note that this does
	// not include Events that in "child" places of this place.
	Events []*Event
}

// FullName returns p's Name followed by the FullName of its parent Place, if
// any, separated by ", ".
func (p *Place) FullName() string {
	if p.Parent == nil {
		return p.Name
	}
	return p.Name + ", " + p.Parent.FullName()
}

// String returns p's FullName.
func (p *Place) String() string {
	return p.FullName()
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
