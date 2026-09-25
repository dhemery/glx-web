package web

import (
	"fmt"
	"strconv"

	"github.com/genealogix/glx/go-glx"
)

type Property struct {
	Definition *glx.PropertyDefinition
	Values     []PropertyValue
}

func (p Property) String() string {
	return p.Value().String()
}

func (p Property) Value() PropertyValue {
	return p.Values[0]
}

type PropertyValue struct {
	Value  fmt.Stringer
	Date   Date
	Fields map[string]PropertyField
}

func (v PropertyValue) Field(name string) PropertyField {
	return v.Fields[name]

}

func (v PropertyValue) String() string {
	return v.Value.String()
}

type PropertyField struct {
	Definition *glx.FieldDefinition
	Value      fmt.Stringer
}

type Date string

func (d Date) String() string {
	return string(d)
}

type IntValue int

func (i IntValue) String() string {
	return strconv.Itoa((int(i)))
}

type BoolValue bool

func (b BoolValue) String() string {
	return strconv.FormatBool(bool(b))
}

type StringValue string

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
	out.Fields = newPropertyFields(in["fields"], def, a)

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

func newPropertyFields(in any, def *glx.PropertyDefinition, a *Archive) map[string]PropertyField {
	return nil
}

func newReferenceValue(id string, entityType string, a *Archive) fmt.Stringer {
	switch entityType {
	case "citations":
		return a.Citations[id]
	case "persons":
		return a.Persons[id]
	case "places":
		return a.Places[id]
	}
	panic("newReferenceValue unimplemented entity type: " + entityType)
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
	}
	return nil
}

func newDate(in any) Date {
	return Date(fmt.Sprint(in))
}
