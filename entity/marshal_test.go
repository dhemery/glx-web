package entity_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

// JSON marshal options to reduce the JSON representation of the catalog to
// only the details we want to compare. The resulting JSON omits each entity's
// back-references and inner GLX entity. It writes each primitive value and
// vocabulary value as a string that shows its type and value.
var catalogMarshalOptions = json.JoinOptions(
	json.OmitZeroStructFields(true),
	json.WithMarshalers(
		json.JoinMarshalers(
			// Entities
			json.MarshalToFunc(marshalAssertion),
			json.MarshalToFunc(marshalCitation),
			json.MarshalToFunc(marshalEvent),
			json.MarshalToFunc(marshalPerson),
			json.MarshalToFunc(marshalPlace),

			json.MarshalToFunc(marshalParticipant),
			json.MarshalToFunc(marshalPropertyEntry),

			// Vocabulary Values
			json.MarshalToFunc(func(out *jsontext.Encoder, in entity.ConfidenceLevel) error {
				return json.MarshalEncode(out, fmt.Sprintf("ConfidenceLevel{%q}", in.Value()))
			}),
			json.MarshalToFunc(func(out *jsontext.Encoder, in entity.EventType) error {
				return json.MarshalEncode(out, fmt.Sprintf("EventType{%q}", in.Value()))
			}),
			json.MarshalToFunc(func(out *jsontext.Encoder, in entity.ParticipantRole) error {
				return json.MarshalEncode(out, fmt.Sprintf("ParticipantRole{%q}", in.Value()))
			}),
			// GenderType
			// InformationType
			// LegalStatus
			// MediaType
			// PlaceType
			// RelationshipType
			// RepositoryType
			// SexType
			// SourceNature

			// Primitive Values
			json.MarshalToFunc(func(out *jsontext.Encoder, in entity.BoolValue) error {
				return json.MarshalEncode(out, fmt.Sprintf("BoolValue{%t}", in))
			}),
			json.MarshalToFunc(func(out *jsontext.Encoder, in entity.IntValue) error {
				return json.MarshalEncode(out, fmt.Sprintf("IntValue{%d}", in))
			}),
			json.MarshalToFunc(func(out *jsontext.Encoder, in entity.StringValue) error {
				return json.MarshalEncode(out, fmt.Sprintf("StringValue{%q}", in))
			}),

			// GLX types
			json.MarshalToFunc(marshalFieldDefinition),
			json.MarshalToFunc(marshalPropertyDefinition),
			json.MarshalToFunc(marshalVocabularyEntry),
		),
	),
)

func marshalAssertion(out *jsontext.Encoder, in *entity.Assertion) error {
	v := struct {
		Date        glx.DateString      `json:"date"`
		Participant *entity.Participant `json:"participant"`
		Property    string              `json:"property"`
		Subject     entity.Entity       `json:"subject"`
	}{
		Date:        in.Date,
		Participant: in.Participant,
		Property:    in.Property,
		Subject:     in.Subject,
	}
	return json.MarshalEncode(out, &v)
}

func marshalCitation(out *jsontext.Encoder, in *entity.Citation) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.ID()))
	}
	v := struct {
		Properties map[string]*entity.Property `json:"properties"`
	}{
		Properties: in.Properties,
	}
	return json.MarshalEncode(out, &v)
}

func marshalEvent(out *jsontext.Encoder, in *entity.Event) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.ID()))
	}
	v := struct {
		Date         glx.DateString              `json:"date"`
		Participants []*entity.Participant       `json:"participants"`
		Place        *entity.Place               `json:"place"`
		Properties   map[string]*entity.Property `json:"properties"`
		Type         entity.EventType            `json:"type"`
	}{
		Date:         in.Date,
		Participants: in.Participants,
		Place:        in.Place,
		Properties:   in.Properties,
		Type:         in.Type,
	}
	return json.MarshalEncode(out, &v)
}

func marshalFieldDefinition(out *jsontext.Encoder, in *glx.FieldDefinition) error {
	v := struct {
		Label     string `json:"label"`
		ValueType string `json:"value_type"`
	}{
		Label:     in.Label,
		ValueType: in.ValueType,
	}
	return json.MarshalEncode(out, &v)
}

func marshalParticipant(out *jsontext.Encoder, in *entity.Participant) error {
	v := struct {
		Person     *entity.Person              `json:"person"`
		Properties map[string]*entity.Property `json:"properties,omitempty"`
		Role       entity.ParticipantRole      `json:"role"`
	}{
		Person:     in.Person,
		Properties: in.Properties,
		Role:       in.Role,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPerson(out *jsontext.Encoder, in *entity.Person) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.ID()))
	}
	v := struct {
		Properties map[string]*entity.Property `json:"properties"`
	}{
		Properties: in.Properties,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPlace(out *jsontext.Encoder, in *entity.Place) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.ID()))
	}
	v := struct {
		Name       string                      `json:"name"`
		Parent     *entity.Place               `json:"parent"`
		Properties map[string]*entity.Property `json:"properties"`
		Type       entity.PlaceType            `json:"type"`
	}{
		Name:       in.Name,
		Parent:     in.Parent,
		Properties: in.Properties,
		Type:       in.Type,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPropertyDefinition(out *jsontext.Encoder, in *glx.PropertyDefinition) error {
	v := struct {
		Label          string `json:"label"`
		ValueType      string `json:"value_type"`
		ReferenceType  string `json:"reference_type"`
		VocabularyType string `json:"vocabulary_type"`
	}{
		Label:          in.Label,
		ValueType:      in.ValueType,
		ReferenceType:  in.ReferenceType,
		VocabularyType: in.VocabularyType,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPropertyEntry(out *jsontext.Encoder, in entity.PropertyEntry) error {
	v := struct {
		Date   glx.DateString                  `json:"date"`
		Fields map[string]entity.PropertyField `json:"fields,omitempty"`
		Value  entity.Stringer                 `json:"value"`
	}{
		Date:   in.Date,
		Fields: in.Fields,
		Value:  in.Value,
	}
	return json.MarshalEncode(out, &v)

}

func marshalVocabularyEntry(out *jsontext.Encoder, in *glx.VocabularyEntry) error {
	v := struct {
		Label string `json:"label"`
	}{
		Label: in.Label,
	}

	return json.MarshalEncode(out, &v)
}
