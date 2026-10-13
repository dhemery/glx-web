package entity

import (
	"strconv"

	"github.com/genealogix/glx/go-glx"
)

type Participant struct {
	*glx.Participant
	Person     *Person
	Properties map[string]Property
	Role       *VocabularyValue
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

func parseParticipants(ctx *context, inner []glx.Participant, propertyDefs map[string]*glx.PropertyDefinition) ParticipantList {
	ppctx := ctx.Sub("Participants")

	roleDefs := ppctx.GLX.ParticipantRoles
	var participants ParticipantList

	for i, p := range inner {
		pctx := ppctx.Sub(strconv.Itoa(i))
		participants = append(participants, parseParticipant(pctx, &p, roleDefs, propertyDefs))
	}

	return participants
}

func parseParticipant(ctx *context, inner *glx.Participant, roleDefs map[string]*glx.VocabularyEntry, propertyDefs map[string]*glx.PropertyDefinition) *Participant {
	if inner == nil {
		return nil
	}
	// TODO(validation): Does GLX validation guarantee a valid value?
	role, _ := newVocabularyValue(inner.Role, roleDefs)
	return &Participant{
		Participant: inner,
		Person:      ctx.Catalog.PersonsByID[inner.Person],
		Properties:  parseProperties(ctx, inner.Properties, propertyDefs),
		Role:        role,
	}
}
