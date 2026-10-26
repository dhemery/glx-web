// Package tmpl configures the templates that glx-web uses to render pages.
package tmpl

import (
	"path/filepath"

	"github.com/genealogix/glx/go-glx"
)

// Config configures the templates used by glx-web.
type Config struct {
	// User-defined data to make available to templates via the data
	// function.
	Data map[string]any `json:"data"`
	// Templates to render the site.
	Templates Specs `json:"templates"`
}

// Specs specifies the templates for the site and for each entity type.
type Specs struct {
	// Templates that other templates can build on. The key is used as the
	// name of the base. Base templates can build on other base templates.
	Bases map[string]Spec `json:"bases"`
	// Templates for the site as a whole. Each template is applied to the
	// [Site]. The map key is used as the path for the rendered file
	// relative to the output directory.
	Site map[string]Spec `json:"site"`
	// Templates for each entity type.
	Types map[glx.EntityType]Types `json:"types"`
}

// Types specifies templates for an entity type and for entities of that type.
type Types struct {
	// Templates for the entity type as a whole. Each template is applied
	// to the type's [EntitySet]. The map key is used as the path for the
	// rendered file relative to the output directory.
	Type map[string]Spec `json:"type"`
	// Templates for each entity of the type. Each template is applied to
	// the entity. Each occurrence of "<id>" in the map key is replaced
	// with the entity's ID, and the result is used as the path for the
	// rendered file, relative to the output directory.
	Entity map[string]Spec `json:"entity"`
}

// Spec specifies a template. Each pattern is matched to files relative to the
// template directory. The template is compiled by cloning the base and parsing
// the content of each matching file into it. If no base is specified, an
// undefined template named glx-web-root is used as the base.
//
// Patterns are matched according to the sementics of [filepath.Match].
type Spec struct {
	// The name of the base template to build on.
	Base string `json:"base"`
	// The glob patterns that identify the template files to parse,
	// relative to the template directory.
	Patterns []string `json:"patterns"`
}

// Import for comment links.
var _ = filepath.Match
