package entity

import (
	"strconv"

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

func newParticipants(ctx *context, inner []glx.Participant, propertyDefs map[string]*glx.PropertyDefinition) ParticipantList {
	roleDefs := ctx.GLX.ParticipantRoles
	var participants ParticipantList

	for i, p := range inner {
		pctx := ctx.Sub(strconv.Itoa(i))
		participants = append(participants, newParticipant(pctx, &p, roleDefs, propertyDefs))
	}

	return participants
}

func newParticipant(ctx *context, inner *glx.Participant, roleDefs map[string]*glx.VocabularyEntry, propertyDefs map[string]*glx.PropertyDefinition) *Participant {
	if inner == nil {
		return nil
	}
	return &Participant{
		Participant: inner,
		Person:      ctx.Catalog.PersonsByID[inner.Person],
		Properties:  parseProperties(ctx, inner.Properties, propertyDefs),
		Role:        newVocabularyValue(ctx, inner.Role, roleDefs),
	}
}
