package site

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"

	"github.com/genealogix/glx/go-glx"
)

type templates struct {
	site      map[string]*template.Template
	entitySet map[glx.EntityType]map[string]*template.Template
	entity    map[glx.EntityType]*template.Template
}

func (s templates) ForSite() map[string]*template.Template {
	return s.site
}

func (s templates) ForEntitySet(e glx.EntityType) map[string]*template.Template {
	return s.entitySet[e]
}

func (s templates) ForEntity(e glx.EntityType) *template.Template {
	return s.entity[e]
}

func loadTemplates(templateDir string) (templates, error) {
	config, err := loadConfig(templateDir)
	if err != nil {
		return templates{}, fmt.Errorf("loading template configuration: %w", err)
	}

	return loadSiteTemplates(os.DirFS(templateDir), config.Templates)
}

type templateLoader struct {
	fsys          fs.FS
	baseSpecs     map[string]TemplateSpec
	baseTemplates map[string]*template.Template
}

func loadSiteTemplates(fsys fs.FS, spec Spec) (templates, error) {
	siteTemplates := templates{
		site:      make(map[string]*template.Template),
		entitySet: make(map[glx.EntityType]map[string]*template.Template),
		entity:    make(map[glx.EntityType]*template.Template),
	}

	l := templateLoader{
		fsys:          fsys,
		baseSpecs:     spec.Bases,
		baseTemplates: make(map[string]*template.Template),
	}

	// Load the templates for top-level pages.
	for name, sitePageSpec := range spec.Pages {
		siteTemplate, err := l.loadTemplate(sitePageSpec)
		if err != nil {
			return templates{}, fmt.Errorf("loading top-level template %q: %w", name, err)
		}
		siteTemplates.site[name] = siteTemplate
	}

	// Load the templates for each entity type.
	for et, entityTypeSpec := range spec.Entities {
		entitySetTemplates := make(map[string]*template.Template)

		// Load the entity set template.
		for name, setPageSpec := range entityTypeSpec.Pages {
			setTemplate, err := l.loadTemplate(setPageSpec)
			if err != nil {
				return templates{}, fmt.Errorf("index template %q: %w", name, err)
			}
			entitySetTemplates[name] = setTemplate
		}
		siteTemplates.entitySet[et] = entitySetTemplates

		// Load the entity template.
		entityTemplate, err := l.loadTemplate(entityTypeSpec.Entity)
		if err != nil {
			return templates{}, fmt.Errorf("entity template: %w", err)
		}
		siteTemplates.entity[et] = entityTemplate
	}

	return siteTemplates, nil
}

func (l *templateLoader) loadTemplate(spec TemplateSpec) (*template.Template, error) {
	if spec.Base == "" && len(spec.Patterns) == 0 {
		return nil, errors.New("no base or patterns")
	}

	base, err := l.cloneBase(spec.Base)
	if err != nil {
		return nil, err
	}

	// If no base, create a new base with funcs but no body. Name it after
	// the first pattern so that the first parsed file will supply the
	// body.
	if base == nil {
		name := path.Base(spec.Patterns[0])
		base = template.New(name).Funcs(TemplateFuncs)
	}

	// If no patterns, return the cloned/new base as is.
	if len(spec.Patterns) == 0 {
		return base, nil
	}

	return base.ParseFS(l.fsys, spec.Patterns...)
}

// Clone the named base template, loading it first if necessary.
func (l *templateLoader) cloneBase(name string) (*template.Template, error) {
	if name == "" {
		return nil, nil
	}

	var base *template.Template
	var err error
	base, ok := l.baseTemplates[name]
	if !ok { // base has not been loaded. Load it now.
		spec, ok := l.baseSpecs[name]
		if !ok {
			return nil, fmt.Errorf("no spec for base template %q", name)
		}
		base, err = l.loadTemplate(spec)
		if err != nil {
			return nil, fmt.Errorf("loading base template %q: %w", name, err)
		}
		l.baseTemplates[name] = base
	}

	return base.Clone()
}
