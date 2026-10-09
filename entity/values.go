package entity

import (
	"fmt"
	"log/slog"
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

func parsePrimitiveValue(raw any, valueType string, l *slog.Logger) Stringer {
	l = l.With("value_type", valueType)
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
		// TODO(errors): Maybe return error instead.
		l.Warn("cannot parse primitive value", "raw_type", fmt.Sprintf("%T", raw))
		return nil
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

func newVocabularyValue(key string, vocabulary map[string]*glx.VocabularyEntry, l *slog.Logger) VocabularyValue {
	def, ok := vocabulary[key]
	if !ok {
		// TODO(errors): Handle unknown vocabulary value.
		// Maybe return STringValue.
		l.Error("unknown vocabulary value", "value", key)
	}
	return VocabularyValue{Value: key, Definition: def}
}
