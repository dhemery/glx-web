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
		properties[name] = parseProperty(a, rawProperty, defs[name], l.With("property", name))
	}

	return properties
}

// TODO(errors): If parsePropertyValueValue can fail, parseProperty must handle
// the possibility of a Property with no values.
func parseProperty(a archive, rawProperty any, def *glx.PropertyDefinition, l *slog.Logger) Property {
	if def == nil {
		// TODO(errors): Handle unknown property.
		// Maybe assign a synthetic definition and coerce the value to string.
		l.Error("unknown property")
		return Property{}
	}

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
		l.Warn("cannot parse date",
			"element", "property date", "raw_type", fmt.Sprintf("%T", raw))

		return glxdate.Date{}
	}

	return parseDate(s)
}

func parsePropertyValueFields(rawFields any, defs map[string]*glx.FieldDefinition, l *slog.Logger) map[string]PropertyField {
	l = l.With("element", "property fields")

	propertyFields := make(map[string]PropertyField)

	if rawFields == nil {
		return propertyFields
	}

	rawFieldsMap, ok := rawFields.(map[string]any)
	if !ok {
		l.Warn("cannot parse fields", "raw_type", fmt.Sprintf("%T", rawFieldsMap))
		return propertyFields
	}

	for name, rawFieldValue := range rawFieldsMap {
		fieldLogger := l.With("field", name)
		def, ok := defs[name]
		if !ok {
			fieldLogger.Warn("field discarded: no field definition")
			continue
		}
		value := parsePrimitiveValue(rawFieldValue, def.ValueType, fieldLogger)
		if value == nil {
			continue
		}
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
			l.Error("cannot parse reference", "raw_type", raw)
			return nil
		}
		return a.entity(id, def.ReferenceType)
	case def.VocabularyType != "":
		key := raw.(string)
		return newVocabularyValue(key, a.vocabulary(def.VocabularyType), l)
	}

	// GLX validation guarantees exactly one type field has a value.
	l.Error("discarded property value: property definition has no type")
	return nil
}
