package entity

import (
	"fmt"
	"log/slog"

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

func parseProperties(a archive, rawProperties map[string]any, defs map[string]*glx.PropertyDefinition, l *slog.Logger) map[string]Property {
	properties := make(map[string]Property)

	for name, rawProperty := range rawProperties {
		propertyLogger := l.With("property", name)
		def := defs[name]
		if def == nil {
			def = synthesizePropertyDefinition(name)
			propertyLogger.Warn("unknown property: using synthetic property definition",
				"definition", def)
		}
		properties[name] = parseProperty(a, rawProperty, def, propertyLogger)
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

func parseProperty(a archive, rawProperty any, def *glx.PropertyDefinition, l *slog.Logger) Property {

	property := Property{Definition: def}

	switch v := rawProperty.(type) {
	case []any:
		// Convert each element to a property value map and parse it.
		for _, rawPropertyValue := range v {
			propertyValue := parsePropertyValue(a, asPropertyMap(rawPropertyValue), def, l)
			property.Values = append(property.Values, propertyValue)
		}
	case map[string]any:
		// Assume the map is a property value map and parse it.
		property.Values = append(property.Values, parsePropertyValue(a, v, def, l))
	default:
		// Assume raw property is a scalar property value. Wrap it in a
		// property value map and parse it.
		rawValueMap := map[string]any{"value": v}
		property.Values = append(property.Values, parsePropertyValue(a, rawValueMap, def, l))
	}

	return property
}

func asPropertyMap(raw any) map[string]any {
	if rawMap, ok := raw.(map[string]any); ok {
		return rawMap
	}

	return map[string]any{"value": raw}
}

func parsePropertyValue(a archive, rawPropertyValue map[string]any, def *glx.PropertyDefinition, l *slog.Logger) PropertyValue {
	var propertyValue PropertyValue

	propertyValue.Date = parsePropertyValueDate(rawPropertyValue["date"], l)
	propertyValue.Fields = parsePropertyValueFields(rawPropertyValue["fields"], def.Fields, l)
	propertyValue.Value = parsePropertyValueValue(a, rawPropertyValue["value"], def, l)

	return propertyValue
}

func parsePropertyValueDate(raw any, l *slog.Logger) glxdate.Date {
	if raw == nil {
		return glxdate.Date{}
	}

	s, ok := raw.(string)
	if !ok {
		coerced := fmt.Sprint(raw)
		l.Warn("cannot parse date: coercing to string",
			"element", "property date", "coerced", coerced)

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

func parsePropertyValueFields(rawFields any, defs map[string]*glx.FieldDefinition, l *slog.Logger) map[string]PropertyField {
	l = l.With("element", "property fields")

	propertyFields := make(map[string]PropertyField)

	if rawFields == nil {
		return propertyFields
	}

	rawFieldsMap, ok := rawFields.(map[string]any)
	if !ok {
		l.Warn("cannot parse fields", "unparseable", fmt.Sprint(rawFields))
		return propertyFields
	}

	for name, rawFieldValue := range rawFieldsMap {
		fieldLogger := l.With("field", name)
		def, ok := defs[name]
		if !ok {
			def = synthesizeFieldDefinition(name)
			fieldLogger.Warn("unknown field: using synthetic field definition", "definition", def)
			rawFieldValue = fmt.Sprint(rawFieldValue)
		}
		value := parsePrimitiveValue(rawFieldValue, def.ValueType, fieldLogger)
		propertyFields[name] = PropertyField{Definition: def, Value: value}
	}

	return propertyFields
}

func parsePropertyValueValue(a archive, raw any, def *glx.PropertyDefinition, l *slog.Logger) Stringer {
	l = l.With("element", "property value")

	switch {
	case def.ValueType != "":
		return parsePrimitiveValue(raw, def.ValueType, l)
	case def.ReferenceType != "":
		id, ok := raw.(string)
		if !ok {
			// GLX validation guarantees raw is a string.
			// TODO(errors): Fails when raw == nil.
			l.Warn("cannot parse reference: using empty string", "raw_type", raw)
			return StringValue("")
		}
		return a.entity(id, def.ReferenceType)
	case def.VocabularyType != "":
		key := raw.(string)
		return newVocabularyValue(key, a.vocabulary(def.VocabularyType), l)
	}

	// GLX validation guarantees exactly one type field has a value.
	l.Warn("discarded property value: property definition has no type")
	return nil
}
