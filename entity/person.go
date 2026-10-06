package entity

import (
	"fmt"
	"slices"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

// Person represents a GLX person entity.
type Person struct {
	entity
	*glx.Person

	Properties map[string]Property

	// Assertions about this person.
	Assertions AssertionList
	// Events this person participated in.
	Events EventList
	// Relationships this person participated in.
	Relationships RelationshipList
}

// String returns the display name
// composed by [glx.PersonDisplayName].
func (p *Person) String() string {
	return glx.PersonDisplayName(p.Person)
}

type PersonList []*Person

// Sort returns the persons sorted by String value.
func (l PersonList) Sort() PersonList {
	return sortStringers(l)
}

type PersonName struct {
	Prefix        string `json:"prefix"`
	Given         string `json:"given"`
	Nickname      string `json:"nickname"`
	SurnamePrefix string `json:"surname_prefix"`
	Surname       string `json:"surname"`
	Suffix        string `json:"suffix"`
}

// NewPersonName creates a name
// from parts supplied as successive key, value pairs.
// For each pair,
// the value is assigned to the name part identified by the key.
// Name parts are identified
// using the same keys as the "name" property
// in the standard GLX property vocabularyfor Person:
//
//	"prefix"         Prefix
//	"given"          Given
//	"nickname"       Nickname
//	"surname_prefix" SurnamePrefix
//	"surname"        Surname
//	"suffix"         Suffix
func NewPersonName(pairs []string) (*PersonName, error) {
	if len(pairs)%2 != 0 {
		return nil, ErrNotPairs(pairs)
	}

	var n PersonName
	var i int
	for i < len(pairs) {
		k := pairs[i]
		v := pairs[i+1]
		i = i + 2
		switch k {
		case "prefix":
			n.Prefix = v
		case "given":
			n.Given = v
		case "nickname":
			n.Nickname = v
		case "surname_prefix":
			n.SurnamePrefix = v
		case "surname":
			n.Surname = v
		case "suffix":
			n.Suffix = v
		default:
			return nil, ErrUnknownKey(k)
		}
	}

	return &n, nil
}

// String returns a string representation of n.
// The result includes the non-empty fields of n
// separated by spaces
// in this order:
// Prefix, Given, Nickname, SurnamePrefix, Surname, Suffix.
// Nickname is double-quoted.
func (n *PersonName) String() string {
	var pairs []string
	if n.Prefix != "" {
		pairs = append(pairs, n.Prefix)
	}
	if n.Given != "" {
		pairs = append(pairs, n.Given)
	}
	if n.Nickname != "" {
		pairs = append(pairs, `"`+n.Nickname+`"`)
	}
	if n.SurnamePrefix != "" {
		pairs = append(pairs, n.SurnamePrefix)
	}
	if n.Surname != "" {
		pairs = append(pairs, n.Surname)
	}
	if n.Suffix != "" {
		pairs = append(pairs, n.Suffix)
	}
	return strings.Join(pairs, " ")
}

// Formatted formets n according to a format specifier.
// Formatted uses the custom verbs
// implemented by [PersonName.Format].
//
// Unlike [fmt.Sprintf] and similar functions in the fmt package,
// which by default format one argument per verb,
// Formatted applies each verb to the same argument:
// the method's receiver:
//
//	n := %PersonName{Given: "George", Surname: "Washington"}
//	f := n.Formatted("%S, %g") // Yields "Washington, George".
func (n *PersonName) Formatted(format string) string {
	return fmt.Sprintf(format, n)
}

// Format implements [fmt.Formatter]
// to add custom verbs to format a PersonName.
// Format is called by [PersonName.Formatted]
// and by [fmt.Sprintf] and similar functions in the fmt package.
// It is not useful to call Format directly.
//
// Each new verb prints a field of the PersonName:
//
//	'<' Prefix
//	'g' Given
//	'n' Nickname
//	'{' SurnamePrefix
//	'S' Surname
//	'>' Suffix
//
// To format multiple fields of a single name,
// see [PersonName.Formatted].
//
// The functions in package fmt work differently.
// They format one argument per verb.
// To format more than one field of a name via [fmt.Sprintf]
// you must supply the name as multiple arguments:
//
//	mother := %PersonName{ ... }
//	child := &PersonName{ ... }
//	f := fmt.Sprintf("%S, %g was the mother of %S, %g", mother, mother, child, child)
//
// Alternatively, you can use explicit argument indexes
// to to apply multiple verbs to the same argument:
//
//	mother := %PersonName{ ... }
//	child := &PersonName{ ... }
//	f := fmt.Sprintf("%[1]S, %[1]g was the mother of %[2]S, %[2]g", mother, child)
func (n *PersonName) Format(f fmt.State, verb rune) {
	// Verbs defined by package fmt.
	switch verb {
	case 's':
		fmt.Fprint(f, n.String())
		return
	case 'q':
		fmt.Fprintf(f, "%q", n.String())
		return
	case 'v':
		fmt.Fprintf(f, "%v", n.GoString())
		return
	case 'x':
		fmt.Fprintf(f, "%x", n.String())
		return
	case 'X':
		fmt.Fprintf(f, "%X", n.String())
		return
	}

	// Verbs specific to PersonName
	switch verb {
	case '<':
		fmt.Fprint(f, n.Prefix)
	case 'g':
		fmt.Fprint(f, n.Given)
	case 'n':
		fmt.Fprint(f, n.Nickname)
	case '{':
		fmt.Fprint(f, n.SurnamePrefix)
	case 'S':
		fmt.Fprint(f, n.Surname)
	case '>':
		fmt.Fprint(f, n.Suffix)
	default:
		fmt.Fprintf(f, "PersonName unknown verb %q", verb)
	}
}

// WithDefaults returns a copy of n
// with its empty fields
// replaced by the corresponding fields of defaults.
func (n *PersonName) WithDefaults(defaults *PersonName) *PersonName {
	return nil
}

type PersonNameList []*PersonName

func (l PersonNameList) SortFormatted(format string) PersonNameList {
	return nil
}

type ErrNotPairs []string

func (e ErrNotPairs) Error() string {
	return fmt.Sprintf("%d args: %s", len(e), []string(e))
}

func (e ErrNotPairs) Is(target error) bool {
	if t, ok := target.(ErrNotPairs); ok {
		return slices.Equal(e, t)
	}
	return false
}

type ErrUnknownKey string

func (e ErrUnknownKey) Error() string {
	return fmt.Sprintf("unknown key: %q", string(e))
}
func (e ErrUnknownKey) Is(target error) bool {
	if t, ok := target.(ErrUnknownKey); ok {
		return e == t
	}
	return false
}

func newPerson(id string, inner *glx.Person) *Person {
	return &Person{
		Person:     inner,
		EntityType: glx.EntityTypePersons,
		id:         id,
	}
}

func (p *Person) compile(catalog *Catalog, glxfile *glx.GLXFile) {
	inner := p.Person
	p.Properties = newProperties(inner.Properties, glxfile.PersonProperties, catalog, glxfile)
}

func (p *Person) addAssertion(a *Assertion) {
	p.Assertions = append(p.Assertions, a)
}

func (p *Person) addEvent(e *Event) {
	p.Events = append(p.Events, e)
}

func (p *Person) addRelationship(e *Relationship) {
	p.Relationships = append(p.Relationships, e)
}

// GoString implements [fmt.GoStringer]
// to implement the 'v' verb used by package fmt.
func (n *PersonName) GoString() string {
	format := "&PersonName{Prefix: %q, Given: %q, Nickname: %q, SurnamePrefix: %q, Surname: %q, Suffix: %q}"
	return fmt.Sprintf(format, n.Prefix, n.Given, n.Nickname, n.SurnamePrefix, n.Surname, n.Suffix)
}
