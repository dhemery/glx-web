// Package config TBD
package config

import (
	"path/filepath"

	"github.com/genealogix/glx/go-glx"
)

// Import for comment links.
var _ = filepath.Match

type Config struct {
	// Data available to templates via the "data" function.
	Data map[string]any `json:"data"`
	// Specification of templates to render the site.
	Templates Templates `json:"templates"`
}

// Templates specifies the templates to render pages for the site, each entity,
// and the entity set for each entity type.
type Templates struct {
	// Templates that other templates can build on. The key is used as the
	// name of the base. Base templates can build on other base templates.
	Bases map[string]TemplateSpec `json:"bases"`
	// Templates to render pages about the site. Each template is applied
	// to the [Config.Data] read from the config.yaml file in the templates
	// directory. The key is used as the file name for the rendered file.
	Site map[string]TemplateSpec `json:"site"`
	// Templates for each entity type. The key is the plural form of the
	// entity type name, e.g. persons.
	Types map[glx.EntityType]TypeTemplates `json:"types"`
}

// TypeTemplates specifies templates for pages about entities and sets of entities
// of a given entity type.
type TypeTemplates struct {
	// Templates to render pages about the set of entities of the type or
	// about the type as a whole. Each template is applied to the type's
	// [EntitySet]. The map key is used as the file name for the rendered
	// file.
	Set map[string]TemplateSpec `json:"set"`
	// The template to render the page about each entity of the type.
	Entity TemplateSpec `json:"entity"`
}

// TemplateSpec specifies how to build a template. If both a base and patterns
// are specified, the template is constructed by parsing the files that match
// the patterns into a clone of the base template. If no patterns are listed,
// the result is a clone of the base. If no base is named, the template is
// constructed by parsing the files that match the patterns. In this case, the
// first pattern should match exactly one file, and the file should define the
// body of the template.
//
// Patterns are applied according to the sementics of [filepath.Match].
type TemplateSpec struct {
	// The name of the base template to build on.
	Base string `json:"base"`
	// The glob patterns that identify the template files to parse,
	// relative to the template directory. If the spec has no Base, the
	// first pattern must supply the template body.
	Patterns []string `json:"patterns"`
}
