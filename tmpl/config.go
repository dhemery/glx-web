/*
Package tmpl compiles the templates
that glx-web uses to render files.

The glx-web command
applies user-specified templates to render files.
It does not supply any built-in templates.

To specify the templates and to supply custom data,
write a config.yaml file in the template directory.
The file is deserialized into a [Config].

# Templates

The [Config.Templates] field defines four kinds of templates:

  - Base templates,
    declared in [SiteSpec.Bases].
    A base template is never invoked directly,
    but can be extended by other templates,
    including by other base templates.

  - Site templates,
    declared in [SiteSpec.Site].
    Each site template is applied to a [funcs.Site].
    You can render the site information to multiple files
    by specifying multiple site templates.
    The map key for the template
    is used as the path for the rendered file
    relative to the output directory.

  - Type templates,
    declared in [TypeSpec.Type].
    Each type template is applied to a list of entities.
    Each list type has an EntityType method
    that returns the [glx.EntityType]
    of the entities in the list.
    You can render an entity list to multiple files
    by specifying multiple type templates.
    The map key for the template
    is used as the path for the rendered file
    relative to the output directory.

  - Entity templates,
    declared in [TypeSpec.Entity].
    Each entity template is applied to an entity.
    You can render an entity to multiple files
    by specifying multiple Entity templates.
    Each occurrence of "<id>" in the map key for each template
    is replaced with the entity's ID,
    and the result is used as the path for the rendered file
    relative to the output directory.

# Custom Data

Use the [Config.Data] field
to supply custom data about the site.
This data is copied to [context.Site.Data],
making it available to all templates.
*/
package tmpl

import (
	"github.com/dhemery/glx-web/context"
	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

// Config configures the templates used by glx-web.
type Config struct {
	// User-defined data to make available to templates via the data
	// function.
	Data any `json:"data"`
	// Templates to render the site.
	Templates SiteSpec `json:"templates"`
}

// SiteSpec specifies the templates for the site and for each entity type.
// Each template can use the functions defined by package [context].
type SiteSpec struct {
	// Each Base template defines a base that other templates can build on.
	// The map key is used as the name of the base. Base templates can
	// build on other base templates.
	Bases map[string]Spec `json:"bases"`
	// Each Site template is applied to the [Site]. The map key is used as
	// the path for the rendered file relative to the output directory.
	Site map[string]Spec `json:"site"`
	// Templates for each entity type.
	Types map[glx.EntityType]TypeSpec `json:"types"`
	// Each Partial template can be executed by the format and eval
	// functions to return a string or other result. See the descriptions
	// of those functions for details. The map key is used as the name of
	// the partial.
	Partials map[string]Spec `json:"partials"`
}

// TypeSpec specifies templates for an entity type and for entities of that
// type.
type TypeSpec struct {
	// Each Type template is applied to the list of all entities of the
	// type. Each entity list has a distinct type with methods relevant to
	// the type of entity. For example, the list for
	// [glx.EntityTypePersons] is [entity.PersonList], which has a
	// GroupBySurname method. Each list type has an EntityType method that
	// returns the [glx.EntityType] of its elements. Some lists have
	// methods to sort, filter, group, or map the elements in ways useful
	// to templates. The map key is used as the path for the rendered file
	// relative to the output directory.
	Type map[string]Spec `json:"type"`
	// Each Entity template is applied to each entity of the type. Each
	// occurrence of "<id>" in the map key is replaced by the entity's ID,
	// and the result is used as the path for the rendered file relative to
	// the output directory.
	Entity map[string]Spec `json:"entity"`
}

// Spec specifies a template. The template is compiled by cloning the base and
// parsing the files into it. If no base is named, there must be at least one
// file, and the first file must define the body of the template.
type Spec struct {
	// The name of the base template to build on.
	Base string `json:"base"`
	// The paths to the template files to parse relative to the template
	// directory.
	Files []string `json:"files"`
}

// Import for doc reference.
var _ entity.PersonList
var _ context.Site
