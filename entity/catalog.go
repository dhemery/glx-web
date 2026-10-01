// Package entity wraps the entities of a GLX archive to make them easier to
// render via templates. Each entity has a Slug method that returns the path to
// the directory where it is rendered relative to the output directory. Entity
// references are resolved to pointers to the referenced entities. Properties
// and vocabulary values are presented along with their definitions. Most
// entities, properties, and vocabulary values have String methods for
// convenience. The String methods of some entities (Assertion, Citation, and
// Relationship) return strings that, though useful for debugging templates,
// may not be especially useful for display.
package entity

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

type entity struct {
	EntityType glx.EntityType
	ID         string
}

// Slug returns the path to the directory where glx-web renders the entity,
// relative to the site's output directory.
func (e entity) Slug() string {
	return path.Join(e.EntityType.Plural(), e.ID)
}

// Catalog collects all entities.
type Catalog struct {
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

func NewArchive(g *glx.GLXFile) *Catalog {
	a := &Catalog{
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
		a.Assertions[id] = newAssertion(id, ga)
	}

	for id, gc := range g.Citations {
		a.Citations[id] = newCitation(id, gc)
	}

	for id, ge := range g.Events {
		a.Events[id] = newEvent(id, ge)
	}

	for id, gm := range g.Media {
		a.Media[id] = newMedia(id, gm)
	}

	for id, gp := range g.Persons {
		a.Persons[id] = newPerson(id, gp)
	}

	for id, gp := range g.Places {
		a.Places[id] = newPlace(id, gp)
	}

	for id, gr := range g.Relationships {
		a.Relationships[id] = newRelationship(id, gr)
	}

	for id, gr := range g.Repositories {
		a.Repositories[id] = newRepository(id, gr)
	}

	for id, gs := range g.Sources {
		a.Sources[id] = newSource(id, gs)
	}

	a.compile(a.Assertions, g)
	a.compile(a.Citations, g)
	a.compile(a.Events, g)
	a.compile(a.Media, g)
	a.compile(a.Persons, g)
	a.compile(a.Places, g)
	a.compile(a.Relationships, g)
	a.compile(a.Repositories, g)
	a.compile(a.Sources, g)

	return a
}

type compiler interface {
	// Compile initializes fields and properties that refer to entities and
	// vocabularies in the archive.
	compile(*Catalog, *glx.GLXFile)
}

func (a *Catalog) compile[C compiler](compilers map[string]C, g *glx.GLXFile) {
	for _, compiler := range compilers {
		compiler.compile(a, g)
	}
}

type entityReference interface {
	Slug() string
	String() string
}

func (a *Catalog) entity(id string, entityType string) entityReference {
	switch entityType {
	case "citations":
		return a.Citations[id]
	case "persons":
		return a.Persons[id]
	case "places":
		return a.Places[id]
	default:
		return nil
	}
}

func (a *Catalog) citationsWithIDs(ids []string) []*Citation {
	var citations []*Citation
	for _, id := range ids {
		citations = append(citations, a.Citations[id])
	}
	return citations
}

func (a *Catalog) mediaWithIDs(ids []string) []*Media {
	var media []*Media
	for _, mediaID := range ids {
		media = append(media, a.Media[mediaID])
	}
	return media
}

func (a *Catalog) sourcesWithIDs(ids []string) []*Source {
	var sources []*Source
	for _, id := range ids {
		sources = append(sources, a.Sources[id])
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
