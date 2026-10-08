package entity

import (
	"log/slog"

	"github.com/genealogix/glx/go-glx"
)

type archive struct {
	c *Catalog
	g *glx.GLXFile
}

type resolver interface {
	// Resolve initializes the receiver's references to entities and
	// vocabularies in the archive and notifies any referenced entities
	// that the receiver refers to them.
	resolve(archive, *slog.Logger)
}

func (a archive) resolve[R resolver](resolvers map[string]R, l *slog.Logger) {
	for _, r := range resolvers {
		r.resolve(a, l)
	}
}

func (a archive) entity(id string, entityType string) entityReference {
	// TODO(dale): Handle the rest of the entity types.
	switch entityType {
	case glx.EntityTypeCitations.Plural():
		return a.c.CitationsByID[id]
	case glx.EntityTypePersons.Plural():
		return a.c.PersonsByID[id]
	case glx.EntityTypeMedia.Plural():
		return a.c.PlacesByID[id]
	default:
		return nil
	}
}

func (a archive) citations(ids []string) []*Citation {
	var citations []*Citation
	for _, id := range ids {
		citations = append(citations, a.c.CitationsByID[id])
	}
	return citations
}

func (a archive) media(ids []string) []*Media {
	var media []*Media
	for _, mediaID := range ids {
		media = append(media, a.c.MediaByID[mediaID])
	}
	return media
}

func (a archive) sources(ids []string) []*Source {
	var sources []*Source
	for _, id := range ids {
		s := a.c.SourcesByID[id]
		sources = append(sources, s)
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
		return nil
	}
}
