// Package web represents a GLX archive in a form suitable for Go templates to render as HTML.
//
// Every exported type in this package has a String method to make it easy to
// render each value in a template.
package web

import (
	"fmt"
	"path"

	"github.com/genealogix/glx/go-glx"
)

type entity struct {
	EntityType glx.EntityType
	ID         string
}

func (e entity) Slug() string {
	return path.Join(e.EntityType.Plural(), e.ID)
}

type Archive struct {
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
	compile(*Archive, *glx.GLXFile)
}

func (a *Archive) compile[C compiler](compilers map[string]C, g *glx.GLXFile) {
	for _, compiler := range compilers {
		compiler.compile(a, g)
	}
}

func (a *Archive) entity(id string, entityType string) fmt.Stringer {
	switch entityType {
	case "citations":
		return a.Citations[id]
	case "persons":
		return a.Persons[id]
	case "places":
		return a.Places[id]
	default:
		panic("newReferenceValue unimplemented entity type: " + entityType)
	}
}

func (a *Archive) citationsWithIDs(ids []string) []*Citation {
	var citations []*Citation
	for _, id := range ids {
		citations = append(citations, a.Citations[id])
	}
	return citations
}

func (a *Archive) mediaWithIDs(ids []string) []*Media {
	var media []*Media
	for _, mediaID := range ids {
		media = append(media, a.Media[mediaID])
	}
	return media
}

func (a *Archive) sourcesWithIDs(ids []string) []*Source {
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
