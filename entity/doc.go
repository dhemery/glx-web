// Package entity presents a GLX archive and its entities
// in a form that makes them easier to render via templates.
//
// # Entities
//
// Each entity knows its ID and [glx.EntityType].
//
// Each entity has a PagePath method
// that returns the path to the directory
// where glx-web renders the entity,
// relative to the output directory.
//
// Each entity has a String method that,
// for most entity types,
// yields a reasonable value to use as a page title.
// See each entity's String method for details.
//
// You can use the PagePath and String methods
// link to any entity using this construction
// (using .Source as an example):
//
//	{{ with .Source }}<a href="/{{ .PagePath }}">{{ . }}</a>{{ end }}
//
// Each entity's field references
// are resolved to pointers to the referenced entities.
// For example, each Citation's Source field
// is resolved to a Source entity.
// Each Assertion's Citations field
// is resolved to a list of *Citation.
//
// Currently references via properties are not resolved in this way.
//
// # Back-References
//
// Each entity compiles lists of "back-references,"
// direct references to the entity from other entities.
// For example:
//   - [Person.Events] yields the list of events in which the person is a participant.
//   - [Event.Assertions] yields the list of assertions where the event is the subject.
//   - [Source.Citations] yields the list of citations that cite the source.
//
// Currently this back-reference compilation does not include:
//   - References via properties.
//   - References via an Assertion's property field.
//   - Indirect references,
//     such as an Assertion's indirect reference to a source
//     via a direct reference to a Citation.
//
// # Properties
//
// Each property is represented as a [Property].
// In many cases
// this allows rendering any property
// in a straightforward way,
// without having to reason about
// its property definition
// or its YAML representation.
//
// Each Property knows its definition,
// which makes it easy to access
// its label, description, value type,
// and other details.
//
// Each Property's values
// are represented as a list of [PropertyValue],
// even if the property's definition specifies a single value.
//
// Each Property has a Value method
// that yields the property's first value.
//
// Each Property has a String method
// that yields the property's first value as a string.
//
// # Property Values
//
// Each PropertyValue has a Date and Fields,
// even if the property definition does not specify them.
// The fields are represented as a list of [PropertyField].
//
// # Property Fields
//
// Each PropertyField knows its definition,
// which makes it easy to access
// its label, description, value type,
// and other details.
//
// Each PropertyField has a String method
// that yeilds the field's value as a string.
//
// # Limitations
//
// This package
// generally expects each entity's properties to be well-formed,
// even beyond the guarantees provided by GLX deserialization.
//
// [Assertion.String] yields a value
// that merely identifies the assertion
// instead of describing it.
//
// Templates have no way to sort or filter
// the slices and maps delivered by this package,
// or to group elements by some property
// (in the generic sense of the word).
// For example:
//   - Filter participants by role or ID
//   - Sort events or property values by date.
//   - Group events or relationships by type.
//   - Group persons by surname.
//   - Group places by parent place.
package entity
