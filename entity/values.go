package entity

import (
	"fmt"
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

func parsePrimitiveValue(_ *context, raw any, valueType string) Stringer {
	switch v := raw.(type) {
	case string:
		if valueType == "date" {
			return parseDate(v)
		}
		return StringValue(v)
	case int:
		return IntValue(v)
	case bool:
		return BoolValue(v)
	default:
		coerced := fmt.Sprint(raw)
		return StringValue(coerced)
	}
}

func parseDate(s string) glxdate.Date {
	date, _ := glxdate.Parse(s)

	return date
}

// VocabularyValue represents a value from a GLX vocabulary.
type VocabularyValue struct {
	Definition *glx.VocabularyEntry // The GLX vocabulary's definition of the value.
	Value      string
}

// String returns v's label.
func (v VocabularyValue) String() string {
	return v.Definition.Label
}

func synthesizeVocabularyEntry(term string) *glx.VocabularyEntry {
	return &glx.VocabularyEntry{
		Label:       term,
		Description: "Vocabulary entry synthesized by glx-web",
	}
}

func newVocabularyValue(_ *context, key string, vocabulary map[string]*glx.VocabularyEntry) VocabularyValue {
	def, ok := vocabulary[key]
	if !ok {
		def = synthesizeVocabularyEntry(key)
	}
	return VocabularyValue{Value: key, Definition: def}
}
