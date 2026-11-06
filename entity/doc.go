/*
Package entity
supplements a GLX archive and its entities
with features that make them
more convenient to render via templates.

The top-level item of this package
is the [Catalog] of all entities
defined in your GLX archive.

# Entities

Each entity knows its ID and [glx.EntityType].

Each entity has a String method
that returns a reasonable value to use as a page title.
See each entity's String method for details.

# Entity Fields

Each entity field of vocabulary type
is represented as a type that implements [VocabularyValue].

Entity fields of reference type
are resolved to pointers to the referenced entities.
For example, Event's Place field
is resolved to a *Place.

Each entity embeds the corresponding GLX entity.
This makes is possible to access the underlying GLX entity’s fields
even if the outer entity hides them.
For example,
Each [Event] embeds a [*glx.Event]
that points to the original event from your GLX archive.
If dot is an Event,
you can access the fields of the embedded GLX Event:

	{{.Event.Place}}

# Properties

Each property is represented as a [Property].
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
is represented as a type that implements [VocabularyValue].

# Property Fields

Each PropertyField knows its definition,
which gives access to
its label, description, value type,
and other details.

Each PropertyField has a String method
that yeilds the field's value as a string.

# Vocabulary Values

Each value from a GLX controlled vocabulary
is represented as a type that implements [VocabularyValue].
Each type defines methods that return the value's
label, [*glx.VocabularyEntry], vocabulary name, and raw value.
The String method for each type returns the value's label.

# Entity and Participant Lists

Each entity type has an associated list type.
Each list type
has an EntityType method that
returns the [glx.EntityType] of its entities.

There is also a [ParticipantList] type.

Some list types
have methods to sort, filter, group, and map the elements
in ways often useful for templates.
For example, [PersonList.GroupBySurname]
groups persons by surname.
[ParticipantList.Persons] maps the participants to persons.

# Back-References

Each entity compiles lists of "back-references,"
direct references to the entity from other entities.
For example:
  - [Person.Events] is the list of events in which the person participated.
  - [Event.Assertions] is the list of assertions about the event.
  - [Source.Citations] is the list of citations that cite the source.

Currently this back-reference compilation does not include:
  - References via properties.
  - References via an Assertion's Property field.
  - Indirect references,
    such as an Assertion's indirect reference to a source
    via a direct reference to a Citation.

# Assertions

Each Assertion's Property field
is resolved to a Label
for the asserted field or property
if applicable.

Each Assertion's Value field
is resolved to an entity or VocabularyValue
if applicable.
Otherwise it resolves to a StringValue.

The original, unresolved Property and Value strings
are available through the embedded [*glx.Assertion].
*/
package entity
