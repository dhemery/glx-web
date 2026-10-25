// Package tmpl configures the templates that glx-web uses to render a GLX
// archive.
package tmpl

import (
	"path/filepath"

	"github.com/genealogix/glx/go-glx"
)

// Import for comment links.
var _ = filepath.Match

// Config configures the templates used by glx-web.
type Config struct {
	// Data available to templates via the "data" function.
	Data map[string]any `json:"data"`
	// Specification of templates to render the site.
	Templates Specs `json:"templates"`
}

// Specs specifies the templates for the site and for each entity type.
type Specs struct {
	// Templates that other templates can build on. The key is used as the
	// name of the base. Base templates can build on other base templates.
	Bases map[string]Spec `json:"bases"`
	// Templates for the site as a whole. Each template is applied to the
	// [Site]. The map key is used as the file name for the rendered file,
	// which is written to the output directory.
	Site map[string]Spec `json:"site"`
	// Templates for each entity type.
	Types map[glx.EntityType]Types `json:"types"`
}

// Types specifies templates for an entity type and for entities of that type.
type Types struct {
	// Templates for the entity type as a whole. Each template is applied
	// to the type's [EntitySet]. The map key is used as the file name for
	// the rendered file, which is written to the entity type directory.
	Type map[string]Spec `json:"type"`
	// The template for each entity of the type.
	Entity Spec `json:"entity"`
}

// Spec specifies a template. If both a base and patterns are specified, the
// template is constructed by parsing the files that match the patterns into a
// clone of the base template. If no patterns are listed, the result is a clone
// of the base. If no base is named, the template is constructed by parsing the
// files that match the patterns. In this case, the first pattern should match
// exactly one file, and the file should define the body of the template.
//
// Patterns are applied according to the sementics of [filepath.Match].
type Spec struct {
	// The name of the base template to build on.
	Base string `json:"base"`
	// The glob patterns that identify the template files to parse,
	// relative to the template directory. If the spec has no Base, the
	// first pattern must supply the template body.
	Patterns []string `json:"patterns"`
}
