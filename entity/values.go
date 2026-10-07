package entity

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

// TODO(dale): Probably these functions should return error instead of logging.

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

func newPrimitiveValue(raw any, valueType string, l *slog.Logger) Stringer {
	l = l.With("value_type", valueType)
	switch v := raw.(type) {
	case string:
		if valueType == "date" {
			return newDate(v, l)
		}
		return StringValue(v)
	case int:
		return IntValue(v)
	case bool:
		return BoolValue(v)
	default:
		l.Warn("value dropped: cannot parse type", "type", fmt.Sprintf("%T", raw))
		return nil
	}
}

func newDate(raw any, l *slog.Logger) glxdate.Date {
	s, ok := raw.(string)
	if !ok {
		l.Warn("date dropped: cannot parse type", "type", fmt.Sprintf("%T", raw))
		return glxdate.Date{}
	}

	date, err := glxdate.Parse(s)
	if err != nil {
		l.Warn("date dropped: error parsing date", "error", err)
		return glxdate.Date{}
	}

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
		l.Warn("vocabulary value incomplete: no vocabulary entry", "value", key)
	}
	return VocabularyValue{Value: key, Definition: def}
}
