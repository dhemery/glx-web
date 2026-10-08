package entity

import (
	"log/slog"

	"github.com/genealogix/glx/go-glx"
)

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
	return sortValues(l)
}

func newParticipants(a archive, inner []glx.Participant, propertyDefs map[string]*glx.PropertyDefinition, l *slog.Logger) ParticipantList {
	roleDefs := a.g.ParticipantRoles
	var participants ParticipantList

	for _, p := range inner {
		participants = append(participants, newParticipant(a, &p, roleDefs, propertyDefs, l))
	}

	return participants
}

func newParticipant(a archive, inner *glx.Participant, roleDefs map[string]*glx.VocabularyEntry, propertyDefs map[string]*glx.PropertyDefinition, l *slog.Logger) *Participant {
	if inner == nil {
		return nil
	}
	l = l.With("participant", inner.Person)
	return &Participant{
		Participant: inner,
		Person:      a.c.PersonsByID[inner.Person],
		Properties:  parseProperties(a, inner.Properties, propertyDefs, l),
		Role:        newVocabularyValue(inner.Role, roleDefs, l),
	}
}
