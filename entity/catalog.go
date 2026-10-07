package entity

import (
	"log/slog"
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

	l := slog.Default().With("phase", "create entities")

	for id, ga := range g.Assertions {
		c.AssertionsByID[id] = newAssertion(id, ga, l)
	}

	for id, gc := range g.Citations {
		c.CitationsByID[id] = newCitation(id, gc)
	}

	for id, ge := range g.Events {
		c.EventsByID[id] = newEvent(id, ge, l)
	}

	for id, gm := range g.Media {
		c.MediaByID[id] = newMedia(id, gm, l)
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
		c.SourcesByID[id] = newSource(id, gs, l)
	}

	l = slog.Default().With("phase", "resolve references")

	a := archive{c: c, g: g}
	a.resolve(c.AssertionsByID, l)
	a.resolve(c.CitationsByID, l)
	a.resolve(c.EventsByID, l)
	a.resolve(c.MediaByID, l)
	a.resolve(c.PersonsByID, l)
	a.resolve(c.PlacesByID, l)
	a.resolve(c.RelationshipsByID, l)
	a.resolve(c.RepositoriesByID, l)
	a.resolve(c.SourcesByID, l)

	return c
}
