// Package web represents a GLX archive in a form suitable for Go templates to render as HTML.
//
// Every exported type in this package has a String method to make it easy to
// render each value in a template.
package web

import (
	"github.com/genealogix/glx/go-glx"
)

type Archive struct {
	g             *glx.GLXFile
	Assertions    map[string]*Assertion
	Citations     map[string]*Citation
	Events        map[string]*Event
	Media         map[string]*Media
	Persons       map[string]*Person
	Places        map[string]*Place
	Relationships map[string]*Relationship
	Repositories  map[string]*Repository
	Sources       map[string]*Source
}

func NewArchive(g *glx.GLXFile) *Archive {
	a := &Archive{
		g:             g,
		Assertions:    make(map[string]*Assertion),
		Citations:     make(map[string]*Citation),
		Events:        make(map[string]*Event),
		Media:         make(map[string]*Media),
		Persons:       make(map[string]*Person),
		Places:        make(map[string]*Place),
		Relationships: make(map[string]*Relationship),
		Repositories:  make(map[string]*Repository),
		Sources:       make(map[string]*Source),
	}

	for id, ga := range g.Assertions {
		a.Assertions[id] = newAssertion(id, ga, a)
	}

	for id, gc := range g.Citations {
		a.Citations[id] = newCitation(id, gc, a)
	}

	for id, ge := range g.Events {
		a.Events[id] = newEvent(id, ge, a)
	}

	for id, gm := range g.Media {
		a.Media[id] = newMedia(id, gm, a)
	}

	for id, gp := range g.Persons {
		a.Persons[id] = newPerson(id, gp, a)
	}

	for id, gp := range g.Places {
		a.Places[id] = newPlace(id, gp, a)
	}

	for id, gr := range g.Relationships {
		a.Relationships[id] = newRelationship(id, gr, a)
	}

	for id, gr := range g.Repositories {
		a.Repositories[id] = newRepository(id, gr, a)
	}

	for id, gs := range g.Sources {
		a.Sources[id] = newSource(id, gs, a)
	}

	return a
}
