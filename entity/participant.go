package entity

import (
	"errors"

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
	participants := ParticipantList{}

	for i, p := range inner {
		pctx := ppctx.Subi(i)
		participant, err := parseParticipant(pctx, &p, roleDefs, propertyDefs)
		if err != nil {
			pctx.Warnf("discarding participant: %s", err)
			continue
		}
		participants = append(participants, participant)
	}

	return participants
}

func parseParticipant(ctx *context, inner *glx.Participant, roleDefs map[string]*glx.VocabularyEntry, propertyDefs map[string]*glx.PropertyDefinition) (*Participant, error) {
	if inner == nil {
		// Assertion can pass nil *Participant.
		return nil, nil
	}

	if inner.Person == "" {
		// Spec requires Person, but GLX validation allows empty.
		return nil, errors.New("person must not be empty")
	}

	// Role is optional, but GLX validates Role if present.
	role, _ := newVocabularyValue(inner.Role, roleDefs)

	return &Participant{
		Participant: inner,
		Person:      ctx.Catalog.PersonsByID[inner.Person],
		Properties:  parseProperties(ctx, inner.Properties, propertyDefs),
		Role:        role,
	}, nil
}
