package render

import (
	"path/filepath"

	"github.com/genealogix/glx/go-glx"
)

// Import for comment links.
var _ = filepath.Match

type Config struct {
	Site      SiteConfig       `json:"site"`
	Templates SiteTemplateSpec `json:"templates"`
}

type SiteConfig struct {
	// Data available to templates via the "site.data" function.
	Data map[string]any `json:"data"`
}

// SiteTemplateSpec specifies the templates to render the site.
type SiteTemplateSpec struct {
	// Templates that other templates can build on. The key is used as the
	// name of the base. Base templates can build on other base templates.
	Bases map[string]TemplateSpec `json:"bases"`
	// Templates for pages about the site. The key is used as the base file
	// name for the rendered file.
	Pages map[string]TemplateSpec `json:"pages"`
	// Templates for each entity type. The key is the plural form of the
	// entity type name, e.g. persons.
	Entities map[glx.EntityType]EntityTypeTemplateSpec `json:"entities"`
}

// EntityTypeTemplateSpec specifies templates for pages about entities and sets
// of entities of a given type.
type EntityTypeTemplateSpec struct {
	// Templates for pages about the set of entities of a given type. The
	// key is used as the base file name for the rendered file.
	Pages map[string]TemplateSpec `json:"pages"`
	// The template for the page about an entity of the associated type.
	Entity TemplateSpec `json:"entity"`
}

// TemplateSpec specifies how to build a template. If both a base and patterns
// are specified, the template is constructed by parsing the patterns into a
// clone of the base template. If no patterns are listed, the resulting
// template is a clone of base. If no base is named, the template is
// constructed by parsing the files that match the patterns.
//
// Patterns are applied according to the sementics of [filepath.Match].
type TemplateSpec struct {
	// The name of the base template to build on.
	Base string `json:"base"`
	// The glob patterns that identify the files to parse, relative to the
	// template directory.
	Patterns []string `json:"patterns"`
}
