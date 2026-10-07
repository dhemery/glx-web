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

func newProperties(a archive, rawProperties map[string]any, defs map[string]*glx.PropertyDefinition) map[string]Property {
	properties := make(map[string]Property)

	for name, rawProperty := range rawProperties {
		properties[name] = newProperty(a, rawProperty, defs[name])
	}

	return properties
}

func newProperty(a archive, rawProperty any, def *glx.PropertyDefinition) Property {
	property := Property{Definition: def}

	switch v := rawProperty.(type) {
	case []any: // In is multi-valued.
		for _, rawPropertyValue := range v {
			propertyValue := newPropertyValue(a, asPropertyMap(rawPropertyValue), def)
			property.Values = append(property.Values, propertyValue)
		}
	case map[string]any: // In is an object.
		property.Values = append(property.Values, newPropertyValue(a, v, def))
	default: // In is a single non-object value.
		rawValueMap := map[string]any{"value": v}
		property.Values = append(property.Values, newPropertyValue(a, rawValueMap, def))
	}

	return property
}

func asPropertyMap(raw any) map[string]any {
	if rawMap, ok := raw.(map[string]any); ok {
		return rawMap
	}

	return map[string]any{"value": raw}
}

func newPropertyValue(a archive, rawPropertyValue map[string]any, def *glx.PropertyDefinition) PropertyValue {
	propertyValue := PropertyValue{
		Date:   newDate(fmt.Sprint(rawPropertyValue["date"])),
		Fields: newPropertyFields(rawPropertyValue["fields"], def),
	}

	rawValue := rawPropertyValue["value"]

	// TODO(dale): Does GLX guarantee that exactly one type field names a type?
	switch {
	case def.ValueType != "":
		propertyValue.Value = newPrimitiveValue(rawValue, def.ValueType)
	case def.ReferenceType != "":
		// TODO(dale): Does GLX guarantee a string?
		id := rawValue.(string)
		propertyValue.Value = a.entity(id, def.ReferenceType)
	case def.VocabularyType != "":
		key := rawValue.(string)
		propertyValue.Value = newVocabularyValue(key, a.vocabulary(def.VocabularyType))
	}

	return propertyValue
}

func newPropertyFields(rawFields any, def *glx.PropertyDefinition) map[string]PropertyField {
	propertyFields := make(map[string]PropertyField)

	if rawFields == nil {
		return propertyFields
	}

	rawFieldsMap, ok := rawFields.(map[string]any)
	if !ok {
		// TODO(dale): Warn and return empty map
		propertyFields["ERROR"] = PropertyField{
			Value: StringValue(fmt.Sprintf("property fields has unexpected type %T", rawFields)),
		}
	}

	for name, rawFieldValue := range rawFieldsMap {
		propertyFields[name] = PropertyField{
			// TODO(dale): If no definition, warn and store string value.
			Definition: def.Fields[name],
			Value:      newPrimitiveValue(rawFieldValue, def.Fields[name].ValueType),
		}
	}

	return propertyFields
}
