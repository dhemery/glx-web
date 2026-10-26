package load

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"

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

func compileTemplates(fsys fs.FS, specs tmpl.Specs, funcs template.FuncMap) (templateSet, error) {
	templates := templateSet{
		site:      make(map[string]*template.Template),
		entitySet: make(map[glx.EntityType]map[string]*template.Template),
		entity:    make(map[glx.EntityType]map[string]*template.Template),
	}

	c := templateCompiler{
		fsys:  fsys,
		specs: specs,
		// Each compiled template ultimately uses root as a base, and
		// so gains funcs.
		root:  template.New("glx-web-root").Funcs(funcs),
		bases: make(map[string]*template.Template),
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
	fsys fs.FS
	// The default base to use when a Spec does not name a base.
	root  *template.Template
	specs tmpl.Specs
	bases map[string]*template.Template
}

func (c *templateCompiler) compileTemplate(spec tmpl.Spec) (*template.Template, error) {
	if spec.Base == "" && len(spec.Patterns) == 0 {
		return nil, errors.New("no base or patterns")
	}

	t, err := c.cloneBase(spec.Base)
	if err != nil {
		return nil, err
	}

	for _, pattern := range spec.Patterns {
		files, err := fs.Glob(c.fsys, pattern)
		if err != nil {
			return nil, fmt.Errorf("pattern %q: %w", pattern, err)
		}
		for _, file := range files {
			b, err := fs.ReadFile(c.fsys, file)
			if err != nil {
				return nil, fmt.Errorf("pattern %q: template file %q: %w",
					pattern, file, err)
			}
			t, err = t.Parse(string(b))
			if err != nil {
				return nil, fmt.Errorf("pattern %q: template file %q: %w",
					pattern, file, err)
			}
		}
	}
	return t, nil
}

// Clone the named base template, compiling it first if necessary.
func (c *templateCompiler) cloneBase(name string) (*template.Template, error) {
	// If no name, use root, which has funcs installed.
	if name == "" {
		return c.root.Clone()
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
