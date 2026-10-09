package entity

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

// Property holds the values of a property of an entity. Each Property holds a
// slice of values even if the GLX property is defined as single-valued.
type Property struct {
	Definition *glx.PropertyDefinition // The GLX definition of the property.
	Values     []PropertyValue
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
	Value  Stringer
	Date   glxdate.Date // The date when the value applied.
	Fields map[string]PropertyField
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
	Definition *glx.FieldDefinition // The GLX definition of the field.
	Value      Stringer
}

// String returns the string representation of the value of f.
func (f PropertyField) String() string {
	return f.Value.String()
}

func parseProperties(ctx *context, rawProperties map[string]any, defs map[string]*glx.PropertyDefinition) map[string]Property {
	properties := make(map[string]Property)

	for name, rawProperty := range rawProperties {
		def := defs[name]
		if def == nil {
			def = synthesizePropertyDefinition(name)
		}
		properties[name] = parseProperty(ctx, rawProperty, def)
	}

	return properties
}

func synthesizePropertyDefinition(name string) *glx.PropertyDefinition {
	return &glx.PropertyDefinition{
		Label:       name,
		Description: "Property definition synthesized by glx-web",
		ValueType:   "string",
	}
}

func parseProperty(ctx *context, rawProperty any, def *glx.PropertyDefinition) Property {

	property := Property{Definition: def}

	switch v := rawProperty.(type) {
	case []any:
		// Convert each element to a property value map and parse it.
		for _, rawPropertyValue := range v {
			propertyValue := parsePropertyValue(ctx, asPropertyMap(rawPropertyValue), def)
			property.Values = append(property.Values, propertyValue)
		}
	case map[string]any:
		// Assume the map is a property value map and parse it.
		property.Values = append(property.Values, parsePropertyValue(ctx, v, def))
	default:
		// Assume raw property is a scalar property value. Wrap it in a
		// property value map and parse it.
		rawValueMap := map[string]any{"value": v}
		property.Values = append(property.Values, parsePropertyValue(ctx, rawValueMap, def))
	}

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

	propertyValue.Date = parsePropertyValueDate(ctx, rawPropertyValue["date"])
	propertyValue.Fields = parsePropertyValueFields(ctx, rawPropertyValue["fields"], def.Fields)
	propertyValue.Value = parsePropertyValueValue(ctx, rawPropertyValue["value"], def)

	return propertyValue
}

func parsePropertyValueDate(_ *context, raw any) glxdate.Date {
	if raw == nil {
		return glxdate.Date{}
	}

	s, ok := raw.(string)
	if !ok {
		coerced := fmt.Sprint(raw)
		return parseDate(coerced)
	}

	return parseDate(s)
}

func synthesizeFieldDefinition(name string) *glx.FieldDefinition {
	return &glx.FieldDefinition{
		Label:       name,
		Description: "Field definition synthesized by glx-web",
		ValueType:   "string",
	}
}

func parsePropertyValueFields(ctx *context, rawFields any, defs map[string]*glx.FieldDefinition) map[string]PropertyField {
	propertyFields := make(map[string]PropertyField)

	if rawFields == nil {
		return propertyFields
	}

	rawFieldsMap, ok := rawFields.(map[string]any)
	if !ok {
		return propertyFields
	}

	for name, rawFieldValue := range rawFieldsMap {
		def, ok := defs[name]
		if !ok {
			def = synthesizeFieldDefinition(name)
			rawFieldValue = fmt.Sprint(rawFieldValue)
		}
		value := parsePrimitiveValue(ctx, rawFieldValue, def.ValueType)
		propertyFields[name] = PropertyField{Definition: def, Value: value}
	}

	return propertyFields
}

func parsePropertyValueValue(ctx *context, raw any, def *glx.PropertyDefinition) Stringer {
	switch {
	case def.ValueType != "":
		return parsePrimitiveValue(ctx, raw, def.ValueType)
	case def.ReferenceType != "":
		id, ok := raw.(string)
		if !ok {
			// GLX validation guarantees raw is a string.
			// TODO(errors): Fails when raw == nil.
			return StringValue("")
		}
		return ctx.C.entity(id, def.ReferenceType)
	case def.VocabularyType != "":
		key := raw.(string)
		return newVocabularyValue(ctx, key, vocabulary(ctx.G, def.VocabularyType))
	}

	// GLX validation guarantees exactly one type field has a value.
	return nil
}
