package entity

import (
	"maps"
	"slices"

	"github.com/genealogix/glx/go-glx"
)

// Catalog is a collection of entities compiled from a GLX archive.
type Catalog struct {
	AssertionsByID    map[string]*Assertion
	CitationsByID     map[string]*Citation
	EventsByID        map[string]*Event
	MediaByID         map[string]*Media
	PersonsByID       map[string]*Person
	PlacesByID        map[string]*Place
	RelationshipsByID map[string]*Relationship
	RepositoriesByID  map[string]*Repository
	SourcesByID       map[string]*Source
}

func (c *Catalog) Assertions() AssertionList {
	return slices.Collect(maps.Values(c.AssertionsByID))
}

func (c *Catalog) Citations() CitationList {
	return slices.Collect(maps.Values(c.CitationsByID))
}

func (c *Catalog) Events() EventList {
	return slices.Collect(maps.Values(c.EventsByID))
}

func (c *Catalog) Media() MediaList {
	return slices.Collect(maps.Values(c.MediaByID))
}

func (c *Catalog) Persons() PersonList {
	return slices.Collect(maps.Values(c.PersonsByID))
}

func (c *Catalog) Places() PlaceList {
	return slices.Collect(maps.Values(c.PlacesByID))
}

func (c *Catalog) Relationships() RelationshipList {
	return slices.Collect(maps.Values(c.RelationshipsByID))
}

func (c *Catalog) Repositories() RepositoryList {
	return slices.Collect(maps.Values(c.RepositoriesByID))
}

func (c *Catalog) Sources() SourceList {
	return slices.Collect(maps.Values(c.SourcesByID))
}

// NewCatalog compiles a catalog of entities from a GLX archive.
func NewCatalog(g *glx.GLXFile) *Catalog {
	c := &Catalog{
		AssertionsByID:    make(map[string]*Assertion),
		CitationsByID:     make(map[string]*Citation),
		EventsByID:        make(map[string]*Event),
		MediaByID:         make(map[string]*Media),
		PersonsByID:       make(map[string]*Person),
		PlacesByID:        make(map[string]*Place),
		RelationshipsByID: make(map[string]*Relationship),
		RepositoriesByID:  make(map[string]*Repository),
		SourcesByID:       make(map[string]*Source),
	}

	for id, ga := range g.Assertions {
		c.AssertionsByID[id] = newAssertion(id, ga)
	}

	for id, gc := range g.Citations {
		c.CitationsByID[id] = newCitation(id, gc)
	}

	for id, ge := range g.Events {
		c.EventsByID[id] = newEvent(id, ge)
	}

	for id, gm := range g.Media {
		c.MediaByID[id] = newMedia(id, gm)
	}

	for id, gp := range g.Persons {
		c.PersonsByID[id] = newPerson(id, gp)
	}

	for id, gp := range g.Places {
		c.PlacesByID[id] = newPlace(id, gp)
	}

	for id, gr := range g.Relationships {
		c.RelationshipsByID[id] = newRelationship(id, gr)
	}

	for id, gr := range g.Repositories {
		c.RepositoriesByID[id] = newRepository(id, gr)
	}

	for id, gs := range g.Sources {
		c.SourcesByID[id] = newSource(id, gs)
	}

	a := archive{c: c, g: g}
	a.resolve(c.AssertionsByID)
	a.resolve(c.CitationsByID)
	a.resolve(c.EventsByID)
	a.resolve(c.MediaByID)
	a.resolve(c.PersonsByID)
	a.resolve(c.PlacesByID)
	a.resolve(c.RelationshipsByID)
	a.resolve(c.RepositoriesByID)
	a.resolve(c.SourcesByID)

	return c
}
