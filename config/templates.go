// Package config defines types to configure glx-web. Currently the
// configuration specifies only the templates to apply to render the output.
//
// To configure the templates, write a config.yaml file in the template
// directory. See [SiteTemplates] for the fields.
//
// The configuration defines four kinds of templates:
//
//   - Base templates,
//     declared in [SiteTemplates.Bases].
//     A base template is never invoked directly,
//     but can be extended by other templates,
//     including by other base templates.
//
//   - Site page templates,
//     declared in [SiteTemplates.Pages].
//     A site page template renders a page
//     about the site.
//     It is applied to an [entity.SiteData],
//     and its output is written to
//     <output-dir>/<index-name>.html.
//     You can render multiple top-level pages
//     by writing multiple site page templates.
//     You will typically want a site page template named "index"
//     to render the home page for the site.
//
//   - Entity type page templates.
//     declared in [EntityTypeTemplates.Pages].
//     An entity type page template renders a page
//     about the collection of entities of a given type.
//     It is applied to an [entity.EntityTypeData],
//     and its output is written to
//     <output-dir>/<entity-type-plural>/<template-name>.html.
//     You can render the collection in a variety of ways
//     by writing multiple entity type page templates.
//     You will typically want an entity type page named "index"
//     for each type,
//     to render a listing of the entities of the type.
//
//   - Entity templates
//     declared in [EntityTypeTemplates.Entity].
//     An entity template renders a page about an entity.
//     It is applied to an [entity.EntityTypeData],
//     and its output is written to
//     <output-dir>/<entity-type-plural>/<entity-id>/index.html.
package config

import (
	"path/filepath"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

var _ = entity.SiteData{}
var _ = filepath.Match

// SiteTemplates specifies the templates for each site page, entity type
// page, and entity page.
type SiteTemplates struct {
	// Templates that other templates can build on. The key is used as the
	// name of the base. Base templates can build on other base templates.
	Bases map[string]TemplateSpec `json:"bases"`
	// Templates for pages about the site. The key is used as the base file
	// name for the rendered file.
	Pages map[string]TemplateSpec `json:"pages"`
	// Templates for each entity type. The key is the plural form of the
	// entity type name, e.g. persons.
	Entities map[glx.EntityType]EntityTypeTemplates `json:"entities"`
}

// EntityTypeTemplates specifies templates for pages about entities and sets
// of entities of a given type.
type EntityTypeTemplates struct {
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
