package entity

import (
	"maps"
	"path"
	"slices"

	"github.com/genealogix/glx/go-glx"
)

type entity struct {
	glx.EntityType
	id string
}

func (e entity) ID() string {
	return e.id
}

// PagePath returns the path to the directory where glx-web renders the entity,
// relative to the site's output directory.
func (e entity) PagePath() string {
	return path.Join(e.Plural(), e.id)
}

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
	a.compile(c.AssertionsByID, g)
	a.compile(c.CitationsByID, g)
	a.compile(c.EventsByID, g)
	a.compile(c.MediaByID, g)
	a.compile(c.PersonsByID, g)
	a.compile(c.PlacesByID, g)
	a.compile(c.RelationshipsByID, g)
	a.compile(c.RepositoriesByID, g)
	a.compile(c.SourcesByID, g)

	return c
}

type archive struct {
	c *Catalog
	g *glx.GLXFile
}

type compiler interface {
	// Compile initializes fields and properties that refer to entities and
	// vocabularies in the archive.
	compile(archive)
}

func (a archive) compile[C compiler](compilers map[string]C, g *glx.GLXFile) {
	for _, compiler := range compilers {
		compiler.compile(a)
	}
}

type entityReference interface {
	PagePath() string
	String() string
}

func (a archive) entity(id string, entityType string) entityReference {
	switch entityType {
	case "citations":
		return a.c.CitationsByID[id]
	case "persons":
		return a.c.PersonsByID[id]
	case "places":
		return a.c.PlacesByID[id]
	default:
		return nil
	}
}

func (a archive) citationsWithIDs(ids []string) []*Citation {
	var citations []*Citation
	for _, id := range ids {
		citations = append(citations, a.c.CitationsByID[id])
	}
	return citations
}

func (a archive) mediaWithIDs(ids []string) []*Media {
	var media []*Media
	for _, mediaID := range ids {
		media = append(media, a.c.MediaByID[mediaID])
	}
	return media
}

func (a archive) sourcesWithIDs(ids []string) []*Source {
	var sources []*Source
	for _, id := range ids {
		sources = append(sources, a.c.SourcesByID[id])
	}
	return sources
}

func (a archive) vocabulary(name string) map[string]*glx.VocabularyEntry {
	switch name {
	case glx.VocabRelationshipTypes:
		return a.g.RelationshipTypes
	case glx.VocabEventTypes:
		return a.g.EventTypes
	case glx.VocabPlaceTypes:
		return a.g.PlaceTypes
	case glx.VocabRepositoryTypes:
		return a.g.RepositoryTypes
	case glx.VocabParticipantRoles:
		return a.g.ParticipantRoles
	case glx.VocabMediaTypes:
		return a.g.MediaTypes
	case glx.VocabConfidenceLevels:
		return a.g.ConfidenceLevels
	case glx.VocabSourceTypes:
		return a.g.SourceTypes
	case glx.VocabSexTypes:
		return a.g.SexTypes
	case glx.VocabGenderTypes:
		return a.g.GenderTypes
	case glx.VocabSearchResultTypes:
		return a.g.SearchResultTypes
	case glx.VocabResearchLogStatusTypes:
		return a.g.ResearchLogStatusTypes
	case glx.VocabStudyTypes:
		return a.g.StudyTypes
	case glx.VocabStudyStatuses:
		return a.g.StudyStatuses
	case glx.VocabLegalStatuses:
		return a.g.LegalStatuses
	case glx.VocabSourceNatures:
		return a.g.SourceNatures
	case glx.VocabInformationTypes:
		return a.g.InformationTypes
	default:
		return make(map[string]*glx.VocabularyEntry)
	}
}
