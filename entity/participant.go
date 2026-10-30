package entity

import (
	"errors"
	"slices"

	"github.com/genealogix/glx/go-glx"
)

type Participant struct {
	*glx.Participant
	Person     *Person
	Properties map[string]*Property
	Role       *VocabularyValue
}

// HasRole reports whether p's Role is present in roles.
func (p *Participant) HasRole(roles ...string) bool {
	return slices.Contains(roles, p.Role.Value)
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

// Diff returns the participants in l that do not appear in participants.
func (l ParticipantList) Diff(participants ParticipantList) ParticipantList {
	return slices.DeleteFunc(slices.Clone(l), func(p *Participant) bool {
		return slices.Contains(participants, p)
	})
}

// Persons returns the persons of the participants in l.
func (l ParticipantList) Persons() PersonList {
	var persons PersonList
	for _, p := range l {
		persons = append(persons, p.Person)
	}
	return persons
}

// Sort returns the participants sorted by String value.
func (l ParticipantList) Sort() ParticipantList {
	return sortedByString(l)
}

// WithRoles returns the participants in l whose Role is present in roles.
func (l ParticipantList) WithRoles(roles ...string) ParticipantList {
	var matches ParticipantList
	for _, p := range l {
		if p.HasRole(roles...) {
			matches = append(matches, p)
		}
	}
	return matches
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
