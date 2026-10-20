package entity_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

// JSON marshal options to reduce the JSON representation of the catalog to
// only the details we want to compare. The resulting JSON omits each entity's
// back-references and inner GLX entity. It also omits details of GLX
// vocabularies (e.g. Description) that are not useful for validating catalog
// structure.
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

			// Primitive Values
			json.MarshalToFunc(marshalBoolValue),
			json.MarshalToFunc(marshalDate),
			json.MarshalToFunc(marshalIntValue),
			json.MarshalToFunc(marshalStringValue),

			// GLX types
			json.MarshalToFunc(marshalFieldDefinition),
			json.MarshalToFunc(marshalPropertyDefinition),
			json.MarshalToFunc(marshalVocabularyEntry),
		),
	),
)

func marshalAssertion(out *jsontext.Encoder, in *entity.Assertion) error {
	v := struct {
		Date        glxdate.Date        `json:"date"`
		Participant *entity.Participant `json:"participant"`
		Property    string              `json:"property"`
		Subject     entity.Stringer     `json:"subject"`
	}{
		Date:        in.Date,
		Participant: in.Participant,
		Property:    in.Property,
		Subject:     in.Subject,
	}
	return json.MarshalEncode(out, &v)
}

func marshalBoolValue(out *jsontext.Encoder, in entity.BoolValue) error {
	return out.WriteToken(jsontext.String(fmt.Sprintf("BoolValue{%t}", in)))
}

func marshalCitation(out *jsontext.Encoder, in *entity.Citation) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.PagePath()))
	}
	v := struct {
		Properties map[string]entity.Property `json:"properties"`
	}{
		Properties: in.Properties,
	}
	return json.MarshalEncode(out, &v)
}

func marshalDate(out *jsontext.Encoder, in glxdate.Date) error {
	return out.WriteToken(jsontext.String(in.Raw()))
}

func marshalEvent(out *jsontext.Encoder, in *entity.Event) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.PagePath()))
	}
	v := struct {
		Date         glxdate.Date               `json:"date"`
		Participants []*entity.Participant      `json:"participants"`
		Place        *entity.Place              `json:"place"`
		Properties   map[string]entity.Property `json:"properties"`
		Type         *entity.VocabularyValue    `json:"type"`
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

// Marshal IntValue as a formatted string instead of a number. Using string
// prevents json.Unmarshal from unmarshaling the value as float64. The special
// format allows Diff to match the IntValue type as well as the value.
func marshalIntValue(out *jsontext.Encoder, in entity.IntValue) error {
	return out.WriteToken(jsontext.String(fmt.Sprintf("IntValue{%d}", in)))
}

func marshalParticipant(out *jsontext.Encoder, in *entity.Participant) error {
	v := struct {
		Person     *entity.Person             `json:"person"`
		Properties map[string]entity.Property `json:"properties,omitempty"`
		Role       *entity.VocabularyValue    `json:"role"`
	}{
		Person:     in.Person,
		Properties: in.Properties,
		Role:       in.Role,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPerson(out *jsontext.Encoder, in *entity.Person) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.PagePath()))
	}
	v := struct {
		Properties map[string]entity.Property `json:"properties"`
	}{
		Properties: in.Properties,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPlace(out *jsontext.Encoder, in *entity.Place) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.PagePath()))
	}
	v := struct {
		Name       string                     `json:"name"`
		Parent     *entity.Place              `json:"parent"`
		Properties map[string]entity.Property `json:"properties"`
		Type       *entity.VocabularyValue    `json:"type"`
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
		Date   glxdate.Date                    `json:"date"`
		Fields map[string]entity.PropertyField `json:"fields,omitempty"`
		Value  entity.Stringer                 `json:"value"`
	}{
		Date:   in.Date,
		Fields: in.Fields,
		Value:  in.Value,
	}
	return json.MarshalEncode(out, &v)

}

func marshalStringValue(out *jsontext.Encoder, in entity.StringValue) error {
	return out.WriteToken(jsontext.String(fmt.Sprintf("StringValue{%q}", in)))
}

func marshalVocabularyEntry(out *jsontext.Encoder, in *glx.VocabularyEntry) error {
	v := struct {
		Label string `json:"label"`
	}{
		Label: in.Label,
	}

	return json.MarshalEncode(out, &v)
}
