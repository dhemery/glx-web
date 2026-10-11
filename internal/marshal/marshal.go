// Package marshal implements custom marshaling functions for [entity.Catalog]
// elements.
package marshal

import (
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

func Date(e *jsontext.Encoder, d glxdate.Date) error {
	return e.WriteToken(jsontext.String(d.Raw()))
}

func PropertyDefinition(e *jsontext.Encoder, d *glx.PropertyDefinition) error {
	v := struct {
		Label          string `json:"label"`
		ValueType      string `json:"value_type"`
		ReferenceType  string `json:"reference_type"`
		VocabularyType string `json:"vocabulary_type"`
	}{
		Label:          d.Label,
		ValueType:      d.ValueType,
		ReferenceType:  d.ReferenceType,
		VocabularyType: d.VocabularyType,
	}
	return json.MarshalEncode(e, &v)
}

func FieldDefinition(e *jsontext.Encoder, d *glx.FieldDefinition) error {
	v := struct {
		Label     string `json:"label"`
		ValueType string `json:"value_type"`
	}{
		Label:     d.Label,
		ValueType: d.ValueType,
	}
	return json.MarshalEncode(e, &v)
}
