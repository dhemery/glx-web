package entity

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

// Property holds the values of a property of an entity. Each Property holds a
// slice of values even if the GLX property is defined as single-valued.
type Property struct {
	Definition *glx.PropertyDefinition `json:"definition,omitzero"` // The GLX definition of the property.
	Values     []PropertyValue         `json:"values,omitzero,omitempty"`
}

// String returns the first value of p as a string.
func (p Property) String() string {
	return p.Value().String()
}

// Value returns the first value of p.
func (p Property) Value() PropertyValue {
	return p.Values[0]
}

// PropertyValue holds a single value of a property of an entity. Each
// PropertyValue has a date and fields even if the property's definition does
// not specify them.
type PropertyValue struct {
	Value  Stringer                 `json:"value,omitzero,omitempty"`
	Date   glxdate.Date             `json:"date,omitzero"` // The date when the value applied.
	Fields map[string]PropertyField `json:"fields,omitzero,omitempty"`
}

// Field returns the named field of v.
func (v PropertyValue) Field(name string) PropertyField {
	return v.Fields[name]
}

// String returns the string representation of the value of v.
func (v PropertyValue) String() string {
	return v.Value.String()
}

// PropertyField represents a field in a property of an entity.
type PropertyField struct {
	Definition *glx.FieldDefinition `json:"definition"` // The GLX definition of the field.
	Value      Stringer             `json:"value"`
}

// String returns the string representation of the value of f.
func (f PropertyField) String() string {
	return f.Value.String()
}

func parseProperties(ctx *context, rawProperties map[string]any, defs map[string]*glx.PropertyDefinition) map[string]Property {
	ppctx := ctx.Sub("Properties")

	properties := make(map[string]Property)
	for name, rawProperty := range rawProperties {
		def := defs[name]
		if def == nil {
			ppctx.Warnf("property discarded: %s", errUnknown{"property", name})
			continue
		}

		property := parseProperty(ppctx.Sub(name), rawProperty, def)
		if len(property.Values) > 0 {
			properties[name] = property
		}
	}

	return properties
}

func parseProperty(ctx *context, rawProperty any, def *glx.PropertyDefinition) Property {
	var property Property

	switch v := rawProperty.(type) {
	case []any:
		// Convert each element to a property value map and parse it.
		for i, rawPropertyValue := range v {
			pvictx := ctx.Subi(i)
			propertyValue := parsePropertyValue(pvictx, asPropertyMap(rawPropertyValue), def)
			if propertyValue.Value == nil {
				continue
			}
			property.Values = append(property.Values, propertyValue)
		}
	case map[string]any:
		// Assume the map is a property value map and parse it.
		pv := parsePropertyValue(ctx, v, def)
		if pv.Value == nil {
			return property
		}
		property.Values = append(property.Values, pv)
	default:
		// Assume raw property is a scalar property value. Wrap it in a
		// property value map and parse it.
		rawValueMap := map[string]any{"value": v}
		pv := parsePropertyValue(ctx, rawValueMap, def)
		if pv.Value == nil {
			return property
		}
		property.Values = append(property.Values, pv)
	}

	property.Definition = def
	return property
}

func asPropertyMap(raw any) map[string]any {
	if rawMap, ok := raw.(map[string]any); ok {
		return rawMap
	}

	return map[string]any{"value": raw}
}

func parsePropertyValue(ctx *context, rawPropertyValue map[string]any, def *glx.PropertyDefinition) PropertyValue {
	var propertyValue PropertyValue

	value, err := parsePropertyValueValue(ctx, rawPropertyValue["value"], def)
	if err != nil {
		ctx.Warnf("property value discarded: %s", err)
		return propertyValue
	}
	propertyValue.Value = value
	propertyValue.Date = parsePropertyValueDate(ctx, rawPropertyValue["date"])
	propertyValue.Fields = parsePropertyValueFields(ctx, rawPropertyValue["fields"], def.Fields)

	return propertyValue
}

func parsePropertyValueDate(ctx *context, raw any) glxdate.Date {
	if raw == nil {
		return glxdate.Date{}
	}

	d, err := parseDateRaw(raw)
	if err != nil {
		ctx.Warnf("temporal date discarded: %s", err)
	}
	return d
}

func parsePropertyValueFields(ctx *context, rawFields any, defs map[string]*glx.FieldDefinition) map[string]PropertyField {
	propertyFields := make(map[string]PropertyField)

	if rawFields == nil {
		return propertyFields
	}

	rawFieldsMap, ok := rawFields.(map[string]any)
	if !ok {
		ctx.Warnf("all fields discarded: %s", errParse{"object", rawFields})
		return propertyFields
	}

	for name, rawFieldValue := range rawFieldsMap {
		def, ok := defs[name]
		if !ok {
			ctx.Warnf("field discarded: %s", errUnknown{"field", name})
			continue
		}
		value, err := parsePrimitiveValue(rawFieldValue, def.ValueType)
		if err != nil {
			ctx.Warnf("field discarded: %s", err)
			continue
		}
		propertyFields[name] = PropertyField{Definition: def, Value: value}
	}

	return propertyFields
}

func parsePropertyValueValue(ctx *context, raw any, def *glx.PropertyDefinition) (Stringer, error) {
	// TODO(dale): Can raw be nil? GLX validation prevents nil "value" for
	// temporal properties. What about single-value or other multi-value?
	switch {
	case def.ValueType != "":
		return parsePrimitiveValue(raw, def.ValueType)
	case def.ReferenceType != "":
		id, ok := raw.(string)
		if !ok {
			return nil, errParse{fmt.Sprintf("string for %q reference lookup", def.ReferenceType), raw}
		}
		// TODO: Warn if no such entity.
		return ctx.Catalog.entity(id, def.ReferenceType), nil
	case def.VocabularyType != "":
		term, ok := raw.(string)
		if !ok {
			return nil, errParse{fmt.Sprintf("string for %q vocabulary lookup", def.VocabularyType), raw}
		}
		return newVocabularyValue(term, vocabulary(ctx.GLX, def.VocabularyType))
	default:
		// GLX validation guarantees exactly one type field has a value.
		return nil, nil
	}
}
