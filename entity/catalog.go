package entity

import (
	"maps"
	"os"
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

	ctx := &context{C: c, G: g, ErrOut: os.Stderr}
	compile(ctx, c.AssertionsByID)
	compile(ctx, c.CitationsByID)
	compile(ctx, c.EventsByID)
	compile(ctx, c.MediaByID)
	compile(ctx, c.PersonsByID)
	compile(ctx, c.PlacesByID)
	compile(ctx, c.RelationshipsByID)
	compile(ctx, c.RepositoriesByID)
	compile(ctx, c.SourcesByID)

	return c
}

type compiler interface {
	// Compile initializes the receiver's references to entities and
	// vocabularies in the context and notifies any referenced entities
	// that the receiver refers to them.
	compile(*context)
}

func compile[C compiler](ctx *context, compilers map[string]C) {
	for _, r := range compilers {
		r.compile(ctx)
	}
}

func (c *Catalog) entity(id string, entityType string) entityReference {
	// TODO(feature): Handle the rest of the entity types.
	switch entityType {
	case glx.EntityTypeCitations.Plural():
		return c.CitationsByID[id]
	case glx.EntityTypePersons.Plural():
		return c.PersonsByID[id]
	case glx.EntityTypeMedia.Plural():
		return c.PlacesByID[id]
	default:
		return nil
	}
}

func (c *Catalog) citations(ids []string) []*Citation {
	var citations []*Citation
	for _, id := range ids {
		citations = append(citations, c.CitationsByID[id])
	}
	return citations
}

func (c *Catalog) media(ids []string) []*Media {
	var media []*Media
	for _, mediaID := range ids {
		media = append(media, c.MediaByID[mediaID])
	}
	return media
}

func (c *Catalog) sources(ids []string) []*Source {
	var sources []*Source
	for _, id := range ids {
		s := c.SourcesByID[id]
		sources = append(sources, s)
	}
	return sources
}

func vocabulary(g *glx.GLXFile, name string) map[string]*glx.VocabularyEntry {
	switch name {
	case glx.VocabRelationshipTypes:
		return g.RelationshipTypes
	case glx.VocabEventTypes:
		return g.EventTypes
	case glx.VocabPlaceTypes:
		return g.PlaceTypes
	case glx.VocabRepositoryTypes:
		return g.RepositoryTypes
	case glx.VocabParticipantRoles:
		return g.ParticipantRoles
	case glx.VocabMediaTypes:
		return g.MediaTypes
	case glx.VocabConfidenceLevels:
		return g.ConfidenceLevels
	case glx.VocabSourceTypes:
		return g.SourceTypes
	case glx.VocabSexTypes:
		return g.SexTypes
	case glx.VocabGenderTypes:
		return g.GenderTypes
	case glx.VocabSearchResultTypes:
		return g.SearchResultTypes
	case glx.VocabResearchLogStatusTypes:
		return g.ResearchLogStatusTypes
	case glx.VocabStudyTypes:
		return g.StudyTypes
	case glx.VocabStudyStatuses:
		return g.StudyStatuses
	case glx.VocabLegalStatuses:
		return g.LegalStatuses
	case glx.VocabSourceNatures:
		return g.SourceNatures
	case glx.VocabInformationTypes:
		return g.InformationTypes
	default:
		return nil
	}
}
