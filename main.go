// The glx-web command
// builds a website
// that describes the entities in a GLX archive.
// It applies user-specified templates to render pages:
//   - for each entity,
//   - for the set of entities of each type,
//   - and for the website as a whole.
//
// The glx-web command presents the data in a GLX archive
// to templates
// in a form that makes it relatively straightforward to render.
//
// # Known Limitations
//
// The glx-web project
// is a work in progress.
// Every feature is subject to change without notice.
// Each limitation
// mentioned here and in the following sections
// is under consideration.
//
// The glx-web command
// does not render research logs or studies.
// It renders only:
//   - Assertions
//   - Citations
//   - Events
//   - Media
//   - Persons
//   - Places
//   - Relationships
//   - Repositories
//   - Sources
//
// The glx-web command
// generally expects each entity's properties to be well-formed,
// even beyond the guarantees enforced by GLX deserialization.
//
// Templates have no way to sort or filter
// the slices and maps
// that the glx-web command delivers to them,
// or to group elements by some property
// (in the generic sense of the word).
// For example, there is no way to:
//   - Filter participants by role or ID.
//   - Sort events or property values by date.
//   - Group events or relationships by type.
//   - Group persons by surname.
//   - Group places by parent place.
//
// It is akward for templates
// to handle properties differently
// depending on their types:
//
//	{{ range .Properties }}
//	  {{ if .Definition.ReferenceType }}
//	    ... handle reference value
//	  {{ else if .Definition.VocabularyType }}
//	    ... handle vocabulary value
//	  {{ else }}
//	    ... handle primitive value
//	  {{ end }}
//	{{ end ))
//
// The glx-web command
// uses only HTML templates
// to render files.
// There is no way to use text templates instead.
//
// The glx-web command
// does not privatize living persons.
//
// # Entities
//
// Each entity embeds the corresponding GLX entity.
// This makes is possible to access the underlying GLX entity’s fields
// even if the outer entity hides them.
//
// Each entity knows its ID and [glx.EntityType].
//
// Each entity has a PagePath method
// that returns the path to the directory
// where the glx-web command renders the entity,
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
// Entity fields of reference type
// are resolved to pointers to the referenced entities.
// For example, Citation's Source field
// is resolved to a *Source.
// Assertion's Citations field
// is resolved to a []*Citation.
//
// Entity fields of vocabulary type
// are represented as [entity.VocabularyValue].
//
// # Vocabulary Values
//
// Each VocabularyValue
// knows its definition
// (see [glx.VocabularyEntry]),
// which gives access
// to its label, description,
// and other details.
//
// VocabularyValue has a String method
// that yields the label associated with the value.
//
// # Assertion Limitations
//
// Assertion's String method
// currently yields a value
// that merely identifies the assertion
// instead of describing it.
//
// Assertion's Value field
// is currently presented as a plain string,
// even if the property being asserted
// has reference type or vocabulary type.
//
// # Back-References
//
// Each entity compiles lists of "back-references,"
// direct references to the entity from other entities.
// For example:
//   - [entity.Person.Events] yields the list of events in which the person is a participant.
//   - [entity.Event.Assertions] yields the list of assertions where the event is the subject.
//   - [entity.Source.Citations] yields the list of citations that cite the source.
//
// Currently this back-reference compilation does not include:
//   - References via properties.
//   - References via an Assertion's Property field.
//   - Indirect references,
//     such as an Assertion's indirect reference to a source
//     via a direct reference to a Citation.
//
// # Properties
//
// Each property is represented as an [entity.Property].
// In many cases
// this allows rendering any property
// in a straightforward way,
// without (in many cases) having to reason about
// the property definition
// or its YAML representation.
//
// Each Property knows its definition
// (see [glx.PropertyDefinition]),
// which gives access to
// its label, description, value type,
// and other details.
//
// Each Property's values
// are represented as a list of [entity.PropertyValue].
// The values are presented as a list
// even if the entity's GLX file gives only one value,
// and even if the property's definition specifies a single value
// A template that cares only about the first value
// can access it via the Property's Value and String methods.
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
// The fields are represented as a list of [entity.PropertyField].
//
// Each PropertyValue has a String method.
// that yields the value as a string.
//
// If a Property's type is a vocabulary type,
// the Value of each PropertyValue
// is represented as an [entity.VocabularyValue].
//
// # Property Fields
//
// Each PropertyField knows its definition,
// which gives access to
// its label, description, value type,
// and other details.
//
// Each PropertyField has a String method
// that yeilds the field's value as a string.
//
// # Dates
//
// Dates in the following values
// are represented as [glxdate.Date]:
//   - Entity fields of type "date",
//     such as [entity.Event.Date].
//   - Each [entity.PropertyValue.Date],
//     the date when the property value applied.
//   - Each [entity.PropertyValue.Value]
//     in properties of type "date".
//
// # Configuration
//
// The glx-web command
// applies user-specified templates to render pages.
// It does not supply any built-in templates.
//
// To specify the templates,
// write a config.yaml file in the template directory.
// The file is deserialized into a [config.SiteTemplates].
//
// The configuration defines four kinds of templates:
//
//   - Base templates,
//     declared in [config.SiteTemplates.Bases].
//     A base template is never invoked directly,
//     but can be extended by other templates,
//     including by other base templates.
//
//   - Site page templates,
//     declared in [config.SiteTemplates.Pages].
//     A site page template renders a page
//     about the site.
//     It is applied to a [data.Site]
//     and its output is written to
//     <output-dir>/<index-name>.html.
//     You can render multiple top-level pages
//     by writing multiple site page templates.
//     You will typically want a site page template named "index"
//     to render the home page for the site.
//
//   - Entity type page templates.
//     declared in [config.EntityTypeTemplates.Pages].
//     An entity type page template renders a page
//     about the collection of entities of a given type.
//     It is applied to a [data.EntityType]
//     and its output is written to
//     <output-dir>/<entity-type-plural>/<template-name>.html.
//     You can render the collection in a variety of ways
//     by writing multiple entity type page templates.
//     You will typically want an entity type page named "index"
//     for each type,
//     to render a listing of the entities of the type.
//
//   - Entity templates
//     declared in [config.EntityTypeTemplates.Entity].
//     An entity template renders a page about an entity.
//     It is applied to a [data.Entity]
//     and its output is written to
//     <output-dir>/<entity-type-plural>/<entity-id>/index.html.
package main

import (
	"github.com/dhemery/glx-web/cmd"
	"github.com/dhemery/glx-web/config"
	"github.com/dhemery/glx-web/data"
	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

// Imports to support doc comments.
var _ config.SiteTemplates
var _ data.Site
var _ entity.IntValue
var _ glx.EntityType
var _ glxdate.Date

func main() {
	cmd.Execute()
}
