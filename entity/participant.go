package entity

import "github.com/genealogix/glx/go-glx"

type Participant struct {
	*glx.Participant
	Person     *Person
	Properties map[string]Property
	Role       VocabularyValue
}

func newParticipants(inner []glx.Participant, propertyDefs map[string]*glx.PropertyDefinition,
	catalog *Catalog, glxfile *glx.GLXFile) []*Participant {
	roleDefs := glxfile.ParticipantRoles
	var participants []*Participant

	for _, p := range inner {
		participants = append(participants, newParticipant(&p, roleDefs, propertyDefs, catalog, glxfile))
	}

	return participants
}

func newParticipant(inner *glx.Participant, roleDefs map[string]*glx.VocabularyEntry,
	propertyDefs map[string]*glx.PropertyDefinition, catalog *Catalog, glxfile *glx.GLXFile) *Participant {
	if inner == nil {
		return nil

	}
	return &Participant{
		Participant: inner,
		Person:      catalog.Persons[inner.Person],
		Properties:  newProperties(inner.Properties, propertyDefs, catalog, glxfile),
		Role:        newVocabularyValue(inner.Role, roleDefs),
	}
}
