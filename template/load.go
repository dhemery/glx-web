// Package template loads and supplies templates.
package template

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/genealogix/glx/go-glx"
)

type TemplateSpec struct {
	// The name of the base template, if any.
	Base string
	// The patterns to parse.
	Patterns []string
}

type EntityTypeTemplateSpec struct {
	// Specification of the template for pages for entities of this type.
	Entity TemplateSpec
	// Specification of the index templates for this entity type.
	Indexes map[string]TemplateSpec
}

type SiteTemplateSpec struct {
	// The default extension to use to identify template files.
	Ext string
	// Specifications of the base templates that the other templates build on.
	Bases map[string]TemplateSpec
	// Specifications of templates for top-level pages.
	Indexes map[string]TemplateSpec
	// Specifications of the index and entity templates for each entity type.
	Entities map[glx.EntityType]EntityTypeTemplateSpec
}

func Load(templateDir string) (map[glx.EntityType]*template.Template, error) {
	baseFileName := filepath.Join(templateDir, "base.gotmpl")
	baseTemplate, err := template.ParseFiles(baseFileName)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", baseFileName, err)
	}

	tmpl := map[glx.EntityType]*template.Template{}

	for _, t := range glx.AllEntityTypes {
		glob := filepath.Join(templateDir, t.Plural(), "*.gotmpl")

		et, err := baseTemplate.Clone()
		if err != nil {
			return nil, fmt.Errorf("cloning base template: %w", err)
		}

		et, err = et.ParseGlob(glob)
		if err != nil {
			continue // TODO: Load each type explicitly and handle error
		}

		tmpl[t] = et
	}

	it, err := baseTemplate.Clone()
	if err != nil {
		return nil, fmt.Errorf("cloning base template: %w", err)
	}
	indexFileName := filepath.Join(templateDir, "index.gotmpl")
	it, err = it.ParseFiles(indexFileName)
	if err != nil {
		return nil, fmt.Errorf("parsing index template %s: %w", indexFileName, err)
	}
	tmpl["index"] = it

	return tmpl, nil
}
