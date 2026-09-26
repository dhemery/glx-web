package web

import (
	"fmt"
	"strconv"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

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

// VocabularyValue represents a value from a GLX vocabulary.
type VocabularyValue struct {
	Definition *glx.VocabularyEntry // The GLX vocabulary's definition of the value.
	Value      string
}

// String returns v's label.
func (v VocabularyValue) String() string {
	return v.Definition.Label
}

// newDate parses the input as a [glxdate.Date], ignoring errors.
func newDate(s string) glxdate.Date {
	date, _ := glxdate.Parse(s)
	return date
}

func newPrimitiveValue(in any, valueType string) fmt.Stringer {
	switch typedIn := in.(type) {
	case string:
		if valueType == "date" {
			return newDate(typedIn)
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

func newReferenceValue(referentID string, referentType string, a *Archive) fmt.Stringer {
	switch referentType {
	case "citations":
		return a.Citations[referentID]
	case "persons":
		return a.Persons[referentID]
	case "places":
		return a.Places[referentID]
	default:
		panic("newReferenceValue unimplemented entity type: " + referentType)
	}
}

func newVocabularyValue(key string, vocabularyType string, a *Archive) VocabularyValue {
	defs := a.vocabulary(vocabularyType)
	return VocabularyValue{Value: key, Definition: defs[key]}
}
