package entity

import (
	"strconv"

	"github.com/genealogix/glx/go-glx"
)

type Stringer interface {
	String() string
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

func parsePrimitiveValue(rawValue any, valueType string) (Stringer, error) {
	switch valueType {
	case "boolean":
		b, ok := rawValue.(bool)
		if !ok {
			return nil, errParse{"boolean", rawValue}
		}
		return BoolValue(b), nil
	case "date":
		return parseDateString(rawValue)
	case "integer":
		i, ok := rawValue.(int)
		if !ok {
			return nil, errParse{"integer", rawValue}
		}
		return IntValue(i), nil
	case "string", "":
		s, ok := rawValue.(string)
		if !ok {
			return nil, errParse{"string", rawValue}
		}
		return StringValue(s), nil
	}

	return nil, errUnknown{"value type", valueType}
}

func parseDateString(rawDate any) (glx.DateString, error) {
	ds, ok := rawDate.(string)
	if !ok {
		return "", errParse{"date string", rawDate}
	}
	return glx.DateString(ds), nil
}

// VocabularyValue represents a value from a GLX vocabulary. Templates can
// check whether a value is a VocabularyValue by calling the
// "is_vocabulary_value" template function.
type VocabularyValue struct {
	Definition *glx.VocabularyEntry `json:"definition"` // The GLX vocabulary's definition of the value.
	Value      string               `json:"value"`
}

// String returns v's label, or "" if it is undefined.
func (v *VocabularyValue) String() string {
	if v == nil {
		return ""
	}
	return v.Definition.Label
}

// IsVocabularyValue reports whether v is a [VocabularyValue]. Templates can call
// this function via the "is_vocabulary_value" template function.
func IsVocabularyValue(v any) bool {
	_, ok := v.(VocabularyValue)
	return ok
}

func newVocabularyValue(term string, vocabulary map[string]*glx.VocabularyEntry) (*VocabularyValue, error) {
	def, ok := vocabulary[term]
	if !ok {
		return nil, errUnknown{"vocabulary value", term}
	}
	return &VocabularyValue{Value: term, Definition: def}, nil
}
