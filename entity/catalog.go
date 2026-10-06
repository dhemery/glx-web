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
	a := &Catalog{
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
		a.AssertionsByID[id] = newAssertion(id, ga)
	}

	for id, gc := range g.Citations {
		a.CitationsByID[id] = newCitation(id, gc)
	}

	for id, ge := range g.Events {
		a.EventsByID[id] = newEvent(id, ge)
	}

	for id, gm := range g.Media {
		a.MediaByID[id] = newMedia(id, gm)
	}

	for id, gp := range g.Persons {
		a.PersonsByID[id] = newPerson(id, gp)
	}

	for id, gp := range g.Places {
		a.PlacesByID[id] = newPlace(id, gp)
	}

	for id, gr := range g.Relationships {
		a.RelationshipsByID[id] = newRelationship(id, gr)
	}

	for id, gr := range g.Repositories {
		a.RepositoriesByID[id] = newRepository(id, gr)
	}

	for id, gs := range g.Sources {
		a.SourcesByID[id] = newSource(id, gs)
	}

	a.compile(a.AssertionsByID, g)
	a.compile(a.CitationsByID, g)
	a.compile(a.EventsByID, g)
	a.compile(a.MediaByID, g)
	a.compile(a.PersonsByID, g)
	a.compile(a.PlacesByID, g)
	a.compile(a.RelationshipsByID, g)
	a.compile(a.RepositoriesByID, g)
	a.compile(a.SourcesByID, g)

	return a
}

type compiler interface {
	// Compile initializes fields and properties that refer to entities and
	// vocabularies in the archive.
	compile(*Catalog, *glx.GLXFile)
}

func (c *Catalog) compile[C compiler](compilers map[string]C, g *glx.GLXFile) {
	for _, compiler := range compilers {
		compiler.compile(c, g)
	}
}

type entityReference interface {
	PagePath() string
	String() string
}

func (c *Catalog) entity(id string, entityType string) entityReference {
	switch entityType {
	case "citations":
		return c.CitationsByID[id]
	case "persons":
		return c.PersonsByID[id]
	case "places":
		return c.PlacesByID[id]
	default:
		return nil
	}
}

func (c *Catalog) citationsWithIDs(ids []string) []*Citation {
	var citations []*Citation
	for _, id := range ids {
		citations = append(citations, c.CitationsByID[id])
	}
	return citations
}

func (c *Catalog) mediaWithIDs(ids []string) []*Media {
	var media []*Media
	for _, mediaID := range ids {
		media = append(media, c.MediaByID[mediaID])
	}
	return media
}

func (c *Catalog) sourcesWithIDs(ids []string) []*Source {
	var sources []*Source
	for _, id := range ids {
		sources = append(sources, c.SourcesByID[id])
	}
	return sources
}

func vocabulary(glxFile *glx.GLXFile, name string) map[string]*glx.VocabularyEntry {
	switch name {
	case glx.VocabRelationshipTypes:
		return glxFile.RelationshipTypes
	case glx.VocabEventTypes:
		return glxFile.EventTypes
	case glx.VocabPlaceTypes:
		return glxFile.PlaceTypes
	case glx.VocabRepositoryTypes:
		return glxFile.RepositoryTypes
	case glx.VocabParticipantRoles:
		return glxFile.ParticipantRoles
	case glx.VocabMediaTypes:
		return glxFile.MediaTypes
	case glx.VocabConfidenceLevels:
		return glxFile.ConfidenceLevels
	case glx.VocabSourceTypes:
		return glxFile.SourceTypes
	case glx.VocabSexTypes:
		return glxFile.SexTypes
	case glx.VocabGenderTypes:
		return glxFile.GenderTypes
	case glx.VocabSearchResultTypes:
		return glxFile.SearchResultTypes
	case glx.VocabResearchLogStatusTypes:
		return glxFile.ResearchLogStatusTypes
	case glx.VocabStudyTypes:
		return glxFile.StudyTypes
	case glx.VocabStudyStatuses:
		return glxFile.StudyStatuses
	case glx.VocabLegalStatuses:
		return glxFile.LegalStatuses
	case glx.VocabSourceNatures:
		return glxFile.SourceNatures
	case glx.VocabInformationTypes:
		return glxFile.InformationTypes
	default:
		return make(map[string]*glx.VocabularyEntry)
	}
}
