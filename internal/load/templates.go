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

func compileSiteSpec(templateDir string, specs tmpl.SiteSpec, funcs template.FuncMap) (templateSet, error) {
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

	for name, spec := range c.specs.Site {
		t, err := c.compileSpec(spec)
		if err != nil {
			return templateSet{}, fmt.Errorf("site template %q: %w", name, err)
		}
		templates.site[name] = t
	}

	for et, typeSpec := range c.specs.Types {
		templates.entitySet[et] = make(map[string]*template.Template)
		for name, spec := range typeSpec.Type {
			t, err := c.compileSpec(spec)
			if err != nil {
				return templateSet{}, fmt.Errorf("type template %q: %w", name, err)
			}
			templates.entitySet[et][name] = t
		}

		templates.entity[et] = make(map[string]*template.Template)
		for name, spec := range typeSpec.Entity {
			t, err := c.compileSpec(spec)
			if err != nil {
				return templateSet{}, fmt.Errorf("entity template %q: %w", name, err)
			}
			templates.entity[et][name] = t
		}
	}

	return templates, nil
}

type templateCompiler struct {
	templateDir string
	funcs       template.FuncMap
	specs       tmpl.SiteSpec
	bases       map[string]*template.Template
}

func (c *templateCompiler) compileSpec(spec tmpl.Spec) (*template.Template, error) {
	base, err := c.base(spec.Base)
	if err != nil {
		return nil, err
	}

	switch {
	case base == nil && len(spec.Files) == 0:
		return nil, errors.New("no base and no files")
	case len(spec.Files) == 0:
		// There's a base but no files to parse. Use the base as is.
		return base, nil
	case base == nil:
		// There are files to parse but no base. Use an undefined base
		// with funcs. Name it after the first file so that parsing the
		// file will define the template.
		name := filepath.Base(spec.Files[0])
		base = template.New(name).Funcs(c.funcs)
	default:
		// Both base and files are specified. Clone the base.
		base, err = base.Clone()
		if err != nil {
			return nil, err
		}
	}

	var fnames []string
	for _, fname := range spec.Files {
		fnames = append(fnames, filepath.Join(c.templateDir, fname))
	}
	return base.ParseFiles(fnames...)
}

// Return the named base template, compiling it first if necessary.
func (c *templateCompiler) base(name string) (*template.Template, error) {
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
		base, err = c.compileSpec(baseSpec)
		if err != nil {
			return nil, fmt.Errorf("base template %q: %w", name, err)
		}
		c.bases[name] = base
	}

	return base, nil
}
