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
		a.Assertions[id] = newAssertion(id, ga)
	}

	for id, gc := range g.Citations {
		a.Citations[id] = newCitation(id, gc)
	}

	for id, ge := range g.Events {
		a.Events[id] = newEvent(id, ge)
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

	a.compile(a.Assertions)
	a.compile(a.Citations)
	a.compile(a.Events)

	return a
}

type compiler interface {
	// Compile initializes fields and properties that refer to entities and
	// vocabularies in the archive.
	compile(a *Archive)
}

func (a *Archive) compile[C compiler](compilers map[string]C) {
	for _, compiler := range compilers {
		compiler.compile(a)
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

func (a *Archive) vocabulary(name string) map[string]*glx.VocabularyEntry {
	g := a.g
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
		return make(map[string]*glx.VocabularyEntry)
	}
}
