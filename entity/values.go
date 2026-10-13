package entity

import (
	"strconv"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

// TODO(errors): These low-level functions should return error instead of
// logging. Let the caller log a warning.

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

func parsePrimitiveValue(raw any, valueType string) (Stringer, error) {
	switch valueType {
	case "boolean":
		b, ok := raw.(bool)
		if !ok {
			return nil, errParse{"boolean", raw}
		}
		return BoolValue(b), nil
	case "date":
		return parseDateRaw(raw)
	case "integer":
		i, ok := raw.(int)
		if !ok {
			return nil, errParse{"integer", raw}
		}
		return IntValue(i), nil
	case "string", "":
		s, ok := raw.(string)
		if !ok {
			return nil, errParse{"string", raw}
		}
		return StringValue(s), nil
	}

	return nil, errUnknown{"value type", valueType}
}

func parseDateRaw(raw any) (glxdate.Date, error) {
	s, ok := raw.(string)
	if !ok {
		return glxdate.Date{}, errParse{"date", raw}
	}
	return parseDate(s), nil
}

func parseDate(s string) glxdate.Date {
	date, _ := glxdate.Parse(s)
	return date
}

// VocabularyValue represents a value from a GLX vocabulary.
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

func newVocabularyValue(term string, vocabulary map[string]*glx.VocabularyEntry) (*VocabularyValue, error) {
	def, ok := vocabulary[term]
	if !ok {
		return nil, errUnknown{"vocabulary value", term}
	}
	return &VocabularyValue{Value: term, Definition: def}, nil
}
