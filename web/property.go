package web

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
	Value  fmt.Stringer
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
	Value      fmt.Stringer
}

// String returns the string representation of the value of f.
func (f PropertyField) String() string {
	return f.Value.String()
}

func newProperties(in map[string]any, defs map[string]*glx.PropertyDefinition, a *Archive, g *glx.GLXFile) map[string]Property {
	properties := make(map[string]Property)

	for name, value := range in {
		properties[name] = newProperty(value, defs[name], a, g)
	}

	return properties
}

func newProperty(in any, def *glx.PropertyDefinition, a *Archive, g *glx.GLXFile) Property {
	property := Property{Definition: def}

	switch typedIn := in.(type) {
	case []any: // In is multi-valued.
		for _, v := range typedIn {
			property.Values = append(property.Values, newPropertyValue(asPropertyMap(v), def, a, g))
		}
	case map[string]any: // In is an object.
		property.Values = append(property.Values, newPropertyValue(typedIn, def, a, g))
	default: // In is a single non-object value.
		inMap := map[string]any{"value": typedIn}
		property.Values = append(property.Values, newPropertyValue(inMap, def, a, g))
	}

	return property
}

func asPropertyMap(in any) map[string]any {
	if inMap, ok := in.(map[string]any); ok {
		return inMap
	}

	return map[string]any{"value": in}
}

func newPropertyValue(in map[string]any, def *glx.PropertyDefinition, a *Archive, glxFile *glx.GLXFile) PropertyValue {
	propertyValue := PropertyValue{
		Date:   newDate(fmt.Sprint(in["date"])),
		Fields: newPropertyFields(in["fields"], def),
	}

	inValue := in["value"]
	switch {
	case def.ValueType != "":
		propertyValue.Value = newPrimitiveValue(inValue, def.ValueType)
	case def.ReferenceType != "":
		id := inValue.(string)
		propertyValue.Value = a.entity(id, def.ReferenceType)
	case def.VocabularyType != "":
		key := inValue.(string)
		propertyValue.Value = newVocabularyValue(key, vocabulary(glxFile, def.VocabularyType))
	}

	return propertyValue
}

func newPropertyFields(in any, def *glx.PropertyDefinition) map[string]PropertyField {
	propertyFields := make(map[string]PropertyField)

	if in == nil {
		return propertyFields
	}

	inMap, ok := in.(map[string]any)
	if !ok {
		panic(fmt.Sprintf("Fields has unexpected type %T", in))
	}

	for name, value := range inMap {
		propertyFields[name] = PropertyField{
			Definition: def.Fields[name],
			Value:      newPrimitiveValue(value, def.Fields[name].ValueType),
		}
	}

	return propertyFields
}
