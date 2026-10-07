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

func newProperties(a archive, rawProperties map[string]any, defs map[string]*glx.PropertyDefinition, l *slog.Logger) map[string]Property {
	properties := make(map[string]Property)

	for name, rawProperty := range rawProperties {
		properties[name] = newProperty(a, rawProperty, defs[name], l.With("property", name))
	}

	return properties
}

// TODO(dale): Property can lack definition. Property can be incomplete.
func newProperty(a archive, rawProperty any, def *glx.PropertyDefinition, l *slog.Logger) Property {
	if def == nil {
		l.Warn("property incomplete: no definition")
		return Property{}
	}

	property := Property{Definition: def}

	switch v := rawProperty.(type) {
	case []any:
		// Convert each element to a property value map and parse it.
		for _, rawPropertyValue := range v {
			propertyValue := newPropertyValue(a, asPropertyMap(rawPropertyValue), def, l)
			property.Values = append(property.Values, propertyValue)
		}
	case map[string]any:
		// Assume the map is a property value map and parse it.
		property.Values = append(property.Values, newPropertyValue(a, v, def, l))
	default:
		// Assume raw property is a scalar property value. Wrap it in a
		// property value map and parse it.
		rawValueMap := map[string]any{"value": v}
		property.Values = append(property.Values, newPropertyValue(a, rawValueMap, def, l))
	}

	return property
}

func asPropertyMap(raw any) map[string]any {
	if rawMap, ok := raw.(map[string]any); ok {
		return rawMap
	}

	return map[string]any{"value": raw}
}

func newPropertyValue(a archive, rawPropertyValue map[string]any, def *glx.PropertyDefinition, l *slog.Logger) PropertyValue {
	propertyValue := PropertyValue{
		Fields: newPropertyFields(rawPropertyValue["fields"], def.Fields, l.With("element", "property fields")),
	}

	if rawDate, ok := rawPropertyValue["date"]; ok {
		propertyValue.Date = newDate(rawDate, l.With("element", "property date"))
	}

	rawValue := rawPropertyValue["value"]

	l = l.With("element", "property value")

	// TODO(dale): Does GLX guarantee that exactly one type field names a type?
	switch {
	case def.ValueType != "":
		propertyValue.Value = newPrimitiveValue(rawValue, def.ValueType, l)
	case def.ReferenceType != "":
		// TODO(dale): Does GLX guarantee a string?
		id := rawValue.(string)
		propertyValue.Value = a.entity(id, def.ReferenceType)
	case def.VocabularyType != "":
		key := rawValue.(string)
		propertyValue.Value = newVocabularyValue(key, a.vocabulary(def.VocabularyType, l), l)
	}

	return propertyValue
}

func newPropertyFields(rawFields any, defs map[string]*glx.FieldDefinition, l *slog.Logger) map[string]PropertyField {
	propertyFields := make(map[string]PropertyField)

	if rawFields == nil {
		return propertyFields
	}

	rawFieldsMap, ok := rawFields.(map[string]any)
	if !ok {
		l.Warn("all fields dropped: cannot parse raw type", "type", fmt.Sprintf("%T", rawFieldsMap))
		return propertyFields
	}

	for name, rawFieldValue := range rawFieldsMap {
		fieldLogger := l.With("field", name)
		def, ok := defs[name]
		if !ok {
			fieldLogger.Warn("field dropped: no field definition")
			continue
		}
		value := newPrimitiveValue(rawFieldValue, def.ValueType, fieldLogger)
		if value == nil {
			continue
		}
		propertyFields[name] = PropertyField{Definition: def, Value: value}
	}

	return propertyFields
}
