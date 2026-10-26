package load

import (
	"errors"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/dhemery/glx-web/tmpl"
	"github.com/genealogix/glx/go-glx"
)

type templateSet struct {
	site      map[string]*template.Template
	entitySet map[glx.EntityType]map[string]*template.Template
	entity    map[glx.EntityType]map[string]*template.Template
}

func (s templateSet) ForSite() map[string]*template.Template {
	return s.site
}

func (s templateSet) ForSet(e glx.EntityType) map[string]*template.Template {
	return s.entitySet[e]
}

func (s templateSet) ForEntity(e glx.EntityType) map[string]*template.Template {
	return s.entity[e]
}

func compileTemplates(templateDir string, specs tmpl.Specs, funcs template.FuncMap) (templateSet, error) {
	templates := templateSet{
		site:      make(map[string]*template.Template),
		entitySet: make(map[glx.EntityType]map[string]*template.Template),
		entity:    make(map[glx.EntityType]map[string]*template.Template),
	}

	c := templateCompiler{
		templateDir: templateDir,
		funcs:       funcs,
		specs:       specs,
		bases:       make(map[string]*template.Template),
	}

	for name, siteSpec := range c.specs.Site {
		t, err := c.compileTemplate(siteSpec)
		if err != nil {
			return templateSet{}, fmt.Errorf("site template %q: %w", name, err)
		}
		templates.site[name] = t
	}

	for et, typeSpec := range c.specs.Types {
		setTemplates := make(map[string]*template.Template)

		for name, spec := range typeSpec.Type {
			t, err := c.compileTemplate(spec)
			if err != nil {
				return templateSet{}, fmt.Errorf("type template %q: %w", name, err)
			}
			setTemplates[name] = t
		}
		templates.entitySet[et] = setTemplates

		entityTemplates := make(map[string]*template.Template)
		for name, spec := range typeSpec.Entity {
			t, err := c.compileTemplate(spec)
			if err != nil {
				return templateSet{}, fmt.Errorf("entity template %q: %w", name, err)
			}
			entityTemplates[name] = t
		}
		templates.entity[et] = entityTemplates
	}

	return templates, nil
}

type templateCompiler struct {
	templateDir string
	funcs       template.FuncMap
	specs       tmpl.Specs
	bases       map[string]*template.Template
}

func (c *templateCompiler) compileTemplate(spec tmpl.Spec) (*template.Template, error) {
	base, err := c.cloneBase(spec.Base)
	if err != nil {
		return nil, err
	}

	if len(spec.Files) == 0 {
		if base == nil {
			return nil, errors.New("no base and no files")
		}
		// No files, so use the cloned base.
		return base, nil
	}

	// If there's no base, use an undefined template named after the first
	// file, with funcs installed.
	if base == nil {
		name := filepath.Base(spec.Files[0])
		base = template.New(name).Funcs(c.funcs)
	}

	var fnames []string
	for _, fname := range spec.Files {
		fnames = append(fnames, filepath.Join(c.templateDir, fname))
	}
	return base.ParseFiles(fnames...)
}

// Clone the named base template, compiling it first if necessary.
func (c *templateCompiler) cloneBase(name string) (*template.Template, error) {
	if name == "" {
		return nil, nil
	}

	var base *template.Template
	var err error
	base, ok := c.bases[name]

	// If base has not been compiled, compile it now.
	if !ok {
		baseSpec, ok := c.specs.Bases[name]
		if !ok {
			return nil, fmt.Errorf("base template %q: no such base", name)
		}
		base, err = c.compileTemplate(baseSpec)
		if err != nil {
			return nil, fmt.Errorf("base template %q: %w", name, err)
		}
		c.bases[name] = base
	}

	return base.Clone()
}
