package entity

import "github.com/genealogix/glx/go-glx"

type Participant struct {
	*glx.Participant
	Person     *Person
	Properties map[string]Property
	Role       VocabularyValue
}

// String returns a description of p
// composed by concatenating
// the string values
// of its role and person,
// separated by a space.
func (p *Participant) String() string {
	return p.Role.Definition.Label + " " + p.Person.String()
}

type ParticipantList []*Participant

// Sort returns the participants sorted by String value.
func (l ParticipantList) Sort() ParticipantList {
	return sortStringers(l)
}

func newParticipants(inner []glx.Participant, propertyDefs map[string]*glx.PropertyDefinition,
	catalog *Catalog, glxfile *glx.GLXFile) ParticipantList {
	roleDefs := glxfile.ParticipantRoles
	var participants ParticipantList

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
		Person:      catalog.PersonsByID[inner.Person],
		Properties:  newProperties(inner.Properties, propertyDefs, catalog, glxfile),
		Role:        newVocabularyValue(inner.Role, roleDefs),
	}
}
