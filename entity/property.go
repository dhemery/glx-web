package entity

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

// Property holds the values of a property of an entity. Each Property holds a
// slice of entries even if the GLX property is defined as single-valued.
type Property struct {
	Definition *glx.PropertyDefinition `json:"definition,omitzero"` // The GLX definition of the property.
	Entries    []PropertyEntry         `json:"values,omitzero,omitempty"`
}

// String returns the string value of the first entry of p.
func (p Property) String() string {
	return p.Value().String()
}

// Value returns the first entry of p.
func (p Property) Value() PropertyEntry {
	return p.Entries[0]
}

// PropertyEntry holds a single entry of a property of an entity. Each
// PropertyEntry has a date and fields even if the property's definition does
// not specify them.
type PropertyEntry struct {
	Value  Stringer                 `json:"value,omitzero,omitempty"`
	Date   glxdate.Date             `json:"date,omitzero"` // The date when the value applied.
	Fields map[string]PropertyField `json:"fields,omitzero,omitempty"`
}

// Field returns the named field of e.
func (e PropertyEntry) Field(name string) PropertyField {
	return e.Fields[name]
}

// String returns the string representation of the value of e.
func (e PropertyEntry) String() string {
	return e.Value.String()
}

// An entry is valid if it has either a value or fields.
func (e PropertyEntry) valid() bool {
	return e.Value != nil || len(e.Fields) > 0
}

// PropertyField represents a field in a property entry.
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
		pctx := ppctx.Sub(name)
		def := defs[name]
		if def == nil {
			ppctx.Warnf("property discarded: %s", errUnknown{"property", name})
			continue
		}

		property := parseProperty(pctx, rawProperty, def)
		if len(property.Entries) == 0 {
			pctx.Warn("property discarded: no entries")
			continue
		}

		properties[name] = property
	}

	return properties
}

func appendOrWarn(ctx *context, entries []PropertyEntry, entry PropertyEntry) []PropertyEntry {
	if !entry.valid() {
		ctx.Warn("entry discarded: has neither value nor fields")
		return entries
	}
	return append(entries, entry)
}

func parseProperty(ctx *context, rawProperty any, def *glx.PropertyDefinition) Property {
	var entries []PropertyEntry

	switch v := rawProperty.(type) {
	case []any:
		// Convert each element of rawProperty a property entry map and parse it.
		for i, rawEntry := range v {
			entryCtx := ctx.Subi(i)
			rawEntryMap := asEntryMap(rawEntry)
			entry := parseEntry(entryCtx, rawEntryMap, def)
			entries = appendOrWarn(entryCtx, entries, entry)
		}
	case map[string]any:
		// Assume rawProperty is a property entry map and parse it.
		entry := parseEntry(ctx, v, def)
		entries = appendOrWarn(ctx, entries, entry)
	default:
		// Assume rawProperty is a scalar property value. Wrap it in a
		// property entry map and parse it.
		rawEntryMap := map[string]any{"value": v}
		entry := parseEntry(ctx, rawEntryMap, def)
		entries = appendOrWarn(ctx, entries, entry)
	}

	return Property{
		Definition: def,
		Entries:    entries,
	}
}

func asEntryMap(raw any) map[string]any {
	if rawMap, ok := raw.(map[string]any); ok {
		return rawMap
	}

	return map[string]any{"value": raw}
}

func parseEntry(ctx *context, rawEntry map[string]any, def *glx.PropertyDefinition) PropertyEntry {
	var entry PropertyEntry

	value, err := parseEntryValue(ctx, rawEntry["value"], def)
	if err != nil {
		ctx.Warnf("value discarded: %s", err)
		return entry
	}
	entry.Value = value
	entry.Date = parseEntryDate(ctx, rawEntry["date"])
	entry.Fields = parseEntryFields(ctx, rawEntry["fields"], def.Fields)

	return entry
}

func parseEntryDate(ctx *context, rawDate any) glxdate.Date {
	if rawDate == nil {
		return glxdate.Date{}
	}

	d, err := parseRawDate(rawDate)
	if err != nil {
		ctx.Warnf("date discarded: %s", err)
	}
	return d
}

func parseEntryFields(ctx *context, rawFields any, defs map[string]*glx.FieldDefinition) map[string]PropertyField {
	propertyFields := make(map[string]PropertyField)

	if rawFields == nil {
		return propertyFields
	}

	rawFieldsMap, ok := rawFields.(map[string]any)
	if !ok {
		ctx.Warnf("all fields discarded: %s", errParse{"object", rawFields})
		return propertyFields
	}

	for name, rawValue := range rawFieldsMap {
		def, ok := defs[name]
		if !ok {
			ctx.Warnf("field discarded: %s", errUnknown{"field", name})
			continue
		}
		value, err := parsePrimitiveValue(rawValue, def.ValueType)
		if err != nil {
			ctx.Warnf("field discarded: %s", err)
			continue
		}
		propertyFields[name] = PropertyField{Definition: def, Value: value}
	}

	return propertyFields
}

func parseEntryValue(ctx *context, rawValue any, def *glx.PropertyDefinition) (Stringer, error) {
	// TODO(validation): Spec says that a property entry can hava a value,
	// fields, or both. GLX validation requires a value for temporal
	// properties, but not for single- or multi-valued properties.
	switch {
	case def.ValueType != "":
		return parsePrimitiveValue(rawValue, def.ValueType)
	case def.ReferenceType != "":
		id, ok := rawValue.(string)
		if !ok {
			return nil, errParse{
				want: fmt.Sprintf("string for %q reference lookup", def.ReferenceType),
				got:  rawValue,
			}
		}
		// TODO(validation): Return error if no such entity.
		return ctx.Catalog.entity(id, def.ReferenceType), nil
	case def.VocabularyType != "":
		term, ok := rawValue.(string)
		if !ok {
			return nil, errParse{
				want: fmt.Sprintf("string for %q vocabulary lookup", def.VocabularyType),
				got:  rawValue,
			}
		}
		return newVocabularyValue(term, vocabulary(ctx.GLX, def.VocabularyType))
	default:
		// GLX validation guarantees that exactly one type field of def has a value.
		return nil, nil
	}
}
