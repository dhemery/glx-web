package entity

import "github.com/genealogix/glx/go-glx"

type Participant struct {
	*glx.Participant
	Person     *Person
	Properties map[string]Property
	Role       VocabularyValue
}

func newParticipants(inner []glx.Participant, propertyDefs map[string]*glx.PropertyDefinition, a *Catalog, g *glx.GLXFile) []*Participant {
	roleDefs := g.ParticipantRoles
	var participants []*Participant

	for _, p := range inner {
		participants = append(participants, newParticipant(&p, roleDefs, propertyDefs, a, g))
	}

	return participants
}

func newParticipant(inner *glx.Participant, roleDefs map[string]*glx.VocabularyEntry, propertyDefs map[string]*glx.PropertyDefinition, a *Catalog, g *glx.GLXFile) *Participant {
	if inner == nil {
		return nil

	}
	return &Participant{
		Participant: inner,
		Person:      a.Persons[inner.Person],
		Properties:  newProperties(inner.Properties, propertyDefs, a, g),
		Role:        newVocabularyValue(inner.Role, roleDefs),
	}
}
