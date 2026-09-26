package web

import (
	"fmt"
	"strconv"

	"github.com/genealogix/glx/go-glx"
)

// Property holds the values of a property of an entity. Each Property holds a
// list of values even if the GLX entity has only a single value.
type Property struct {
	// The GLX definition of the property.
	Definition *glx.PropertyDefinition
	// The values of the property.
	Values []PropertyValue
}

// String returns the string representation of the first value of p.
func (p Property) String() string {
	return p.Value().String()
}

// Value returns the first value of p.
func (p Property) Value() PropertyValue {
	return p.Values[0]
}

// PropertyValue holds a single value of a property of an entity. Each
// PropertyValue is represented with fields and a date even if the GLX property
// does not have them.
type PropertyValue struct {
	// Value is the value of the property.
	Value fmt.Stringer
	// Date is the date of the property.
	Date Date
	// Fields holds the property's fields.
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
	Definition *glx.FieldDefinition
	Value      fmt.Stringer
}

// String returns the string representation of the value of f.
func (f PropertyField) String() string {
	return f.Value.String()
}

// Date represents a date value of a property or the date of a value of a
// temporal property.
type Date string

// String returns the string representation of d.
func (d Date) String() string {
	return string(d)
}

// IntValue reprensents an integer value of a property or field.
type IntValue int

// String returns the string representation of i.
func (i IntValue) String() string {
	return strconv.Itoa((int(i)))
}

// BoolValue represents a boolean value of a property or field.
type BoolValue bool

// String returns the string representation of b.
func (b BoolValue) String() string {
	return strconv.FormatBool(bool(b))
}

// StringValue represents a string value of a property or field.
type StringValue string

// String returns the string value of s.
func (s StringValue) String() string {
	return string(s)
}

func newProperties(in map[string]any, defs map[string]*glx.PropertyDefinition, a *Archive) map[string]Property {
	out := make(map[string]Property)

	for name, value := range in {
		out[name] = newProperty(value, defs[name], a)
	}
	return out
}

func newProperty(in any, def *glx.PropertyDefinition, a *Archive) Property {
	prop := Property{Definition: def}

	switch typedIn := in.(type) {
	case []any: // In is multi-valued.
		for _, v := range typedIn {
			prop.Values = append(prop.Values, newPropertyValue(asPropertyMap(v), def, a))
		}
	case map[string]any: // In is an object.
		prop.Values = append(prop.Values, newPropertyValue(typedIn, def, a))
	default: // In is a single non-object value.
		inMap := map[string]any{"value": typedIn}
		prop.Values = append(prop.Values, newPropertyValue(inMap, def, a))
	}

	return prop
}

func asPropertyMap(in any) map[string]any {
	if v, ok := in.(map[string]any); ok {
		// In is already a property map.
		return v
	}

	return map[string]any{"value": in}
}

func newPropertyValue(in map[string]any, def *glx.PropertyDefinition, a *Archive) PropertyValue {
	var out PropertyValue
	out.Date = newDate(in["date"])
	out.Fields = newPropertyFields(in["fields"], def)

	inValue := in["value"]
	switch {
	case def.ValueType != "":
		out.Value = newPrimitiveValue(inValue, def.ValueType)
	case def.ReferenceType != "":
		id := inValue.(string)
		out.Value = newReferenceValue(id, def.ReferenceType, a)
	case def.VocabularyType != "":
		key := inValue.(string)
		out.Value = newVocabularyValue(key, def.VocabularyType, a)
	}

	return out
}

func newPropertyFields(in any, def *glx.PropertyDefinition) map[string]PropertyField {
	var out map[string]PropertyField
	if in == nil {
		return out
	}

	fieldMap, ok := in.(map[string]any)
	if !ok {
		panic(fmt.Sprintf("Fields has unexpected type %T", in))
	}

	for name, value := range fieldMap {
		fieldMap[name] = newPrimitiveValue(value, def.Fields[name].ValueType)
	}
	return out
}

func newReferenceValue(id string, entityType string, a *Archive) fmt.Stringer {
	switch entityType {
	case "citations":
		return a.Citations[id]
	case "persons":
		return a.Persons[id]
	case "places":
		return a.Places[id]
	default:
		panic("newReferenceValue unimplemented entity type: " + entityType)
	}
}

func newVocabularyValue(key string, vocabularyType string, a *Archive) fmt.Stringer {
	return StringValue(key)
}

func newPrimitiveValue(in any, valueType string) fmt.Stringer {
	switch typedIn := in.(type) {
	case string:
		if valueType == "date" {
			return Date(typedIn)
		}
		return StringValue(typedIn)
	case int:
		return IntValue(typedIn)
	case bool:
		return BoolValue(typedIn)
	default:
		panic(fmt.Sprintf("newPrimitiveValue unknown type %T", in))
	}
}

func newDate(in any) Date {
	return Date(fmt.Sprint(in))
}
