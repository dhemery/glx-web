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
