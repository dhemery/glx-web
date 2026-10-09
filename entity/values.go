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

func parsePrimitiveValue(ctx *context, raw any, valueType string) Stringer {
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
		ctx.Warnf("cannot parse primitive value type %T: using string value %q instead",
			raw, coerced)
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

func newVocabularyValue(ctx *context, term string, vocabulary map[string]*glx.VocabularyEntry) VocabularyValue {
	def, ok := vocabulary[term]
	if !ok {

		def = synthesizeVocabularyEntry(term)
		ctx.Warnf("unknown vocabulary term %q: using synthesized vocabulary value with label %q",
			term, def.Label)
	}
	return VocabularyValue{Value: term, Definition: def}
}
