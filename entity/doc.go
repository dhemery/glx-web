/*
Package entity
supplements a GLX archive and its entities
with features that make them
more convenient to render via templates.

# Entities

Each entity embeds the corresponding GLX entity.
This makes is possible to access the underlying GLX entity’s fields
even if the outer entity hides them.

Each entity knows its ID and [glx.EntityType].

Each entity has a String method that,
for most entity types,
yields a reasonable value to use as a page title.
See each entity's String method for details.

Entity fields of reference type
are resolved to pointers to the referenced entities.
For example, Citation's Source field
is resolved to a *Source.
Assertion's Citations field
is resolved to a []*Citation.

Entity fields of vocabulary type
are represented as [VocabularyValue].

# Properties

Each property is represented as a [Property].
In many cases
this allows rendering any property
in a straightforward way,
without (in many cases) having to reason about
the property definition
or its YAML representation.

Each Property knows its definition
(see [glx.PropertyDefinition]),
which gives access to
its label, description, value type,
and other details.

Each Property has a list of [PropertyEntry].
Each entry always has a value, a date, and fields,
even if the property definition
does not specify all of those.
Each property always has a list of entries
even if the entity's GLX file gives only one value,
and even if the property's definition specifies a single value.
A template that cares only about the first entry
can access it via the Property's Entry and String methods.

Each Property has a Value method
that returns the value of the first entry.

Each Property has a String method
that returns the value of the first entry
as a string.

# Property Entries

Each PropertyEntry has a Value, a Date, and Fields,
even if the property definition
does not specify all of those.
The fields are represented as a list of [PropertyField].

Each PropertyEntry has a String method
that returns the entry's value as a string.

If a Property's type is a vocabulary type,
each entry's value
is represented as a [VocabularyValue].

# Property Fields

Each PropertyField knows its definition,
which gives access to
its label, description, value type,
and other details.

Each PropertyField has a String method
that yeilds the field's value as a string.

# Vocabulary Values

Each VocabularyValue
knows its definition
(see [glx.VocabularyEntry]),
which gives access
to its label, description,
and other details.

VocabularyValue has a String method
that yields the label associated with the value.

# Back-References

Each entity compiles lists of "back-references,"
direct references to the entity from other entities.
For example:
  - [Person.Events] yields the list of events in which the person participated.
  - [Event.Assertions] yields the list of assertions about the event.
  - [Source.Citations] yields the list of citations that cite the source.

Currently this back-reference compilation does not include:
  - References via properties.
  - References via an Assertion's Property field.
  - Indirect references,
    such as an Assertion's indirect reference to a source
    via a direct reference to a Citation.

# Dates

Dates in the following values
are represented as [glxdate.Date]:
  - Entity fields of type "date",
    such as [Event.Date].
  - Each [PropertyEntry.Date],
    the date when the entry's value applied.
  - Each [PropertyEntry.Value]
    in properties of type "date".

# Assertion Limitations

Assertion's String method
currently yields a value
that merely identifies the assertion
instead of describing it.

Assertion's Value field
is currently presented as a plain string,
even if the property being asserted
has reference type or vocabulary type.
*/
package entity
