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

type siteTemplates struct {
	// Templates to render top-level pages.
	sitePageTemplates map[string]*template.Template
	// Templates for each entity type.
	typeTemplates map[glx.EntityType]entityTypeTemplates
}

type entityTypeTemplates struct {
	// Templates to render pages about the set of entities of some type.
	typePageTemplates map[string]*template.Template
	// Template to render the page about each entity of the type.
	entityPageTemplate *template.Template
}

func loadTemplates(templateDir string) (siteTemplates, error) {
	config, err := loadConfig(templateDir)
	if err != nil {
		return siteTemplates{}, fmt.Errorf("loading template configuration: %w", err)
	}

	return loadSiteTemplates(os.DirFS(templateDir), config.Templates)
}

type templateLoader struct {
	fsys          fs.FS
	baseSpecs     map[string]TemplateSpec
	baseTemplates map[string]*template.Template
}

func loadSiteTemplates(fsys fs.FS, spec Spec) (siteTemplates, error) {
	st := siteTemplates{
		sitePageTemplates: make(map[string]*template.Template),
		typeTemplates:     make(map[glx.EntityType]entityTypeTemplates),
	}

	l := templateLoader{
		fsys:          fsys,
		baseSpecs:     spec.Bases,
		baseTemplates: make(map[string]*template.Template),
	}

	// Load the templates for top-level pages.
	for name, indexSpec := range spec.Pages {
		tmpl, err := l.loadTemplate(indexSpec)
		if err != nil {
			return st, fmt.Errorf("loading top-level template %q: %w", name, err)
		}
		st.sitePageTemplates[name] = tmpl
	}

	// Load the templates for each entity type.
	for et, entityTypeSpec := range spec.Entities {
		tmpls, err := l.loadEntityTypeTemplates(entityTypeSpec)
		if err != nil {
			return st, fmt.Errorf("loading %s templates: %w", et, err)
		}
		st.typeTemplates[et] = tmpls
	}

	return st, nil
}

func (l templateLoader) loadEntityTypeTemplates(spec EntityTypeSpec) (entityTypeTemplates, error) {
	et := entityTypeTemplates{
		typePageTemplates: make(map[string]*template.Template),
	}

	// Load the templates to render the pages about the set of entities of
	// this type.
	for name, entityTypePageSpec := range spec.Pages {
		tmpl, err := l.loadTemplate(entityTypePageSpec)
		if err != nil {
			return et, fmt.Errorf("index template %q: %w", name, err)
		}
		et.typePageTemplates[name] = tmpl
	}

	// Load the template to render the page about each entity.
	tmpl, err := l.loadTemplate(spec.Entity)
	if err != nil {
		return et, fmt.Errorf("entity template: %w", err)
	}
	et.entityPageTemplate = tmpl
	return et, nil

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
