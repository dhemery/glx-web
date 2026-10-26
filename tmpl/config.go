// Package tmpl configures the templates that glx-web uses to render pages.
package tmpl

import (
	"github.com/genealogix/glx/go-glx"
)

// Config configures the templates used by glx-web.
type Config struct {
	// User-defined data to make available to templates via the data
	// function.
	Data map[string]any `json:"data"`
	// Templates to render the site.
	Templates SiteSpec `json:"templates"`
}

// SiteSpec specifies the templates for the site and for each entity type.
type SiteSpec struct {
	// Templates that other templates can build on. The map key is used as
	// the name of the base. Base templates can build on other base
	// templates.
	Bases map[string]Spec `json:"bases"`
	// Templates for the site as a whole. Each template is applied to the
	// [Site]. The map key is used as the path for the rendered file
	// relative to the output directory.
	Site map[string]Spec `json:"site"`
	// Templates for each entity type.
	Types map[glx.EntityType]TypeSpec `json:"types"`
}

// TypeSpec specifies templates for an entity type and for entities of that
// type.
type TypeSpec struct {
	// Templates for the entity type as a whole. Each template is applied
	// to the type's [EntitySet]. The map key is used as the path for the
	// rendered file relative to the output directory.
	Type map[string]Spec `json:"type"`
	// Templates for each entity of the type. Each template is applied to
	// the entity. Each occurrence of "<id>" in the map key is replaced
	// by the entity's ID, and the result is used as the path for the
	// rendered file relative to the output directory.
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
