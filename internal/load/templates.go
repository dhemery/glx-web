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
	site       map[string]*template.Template
	entityType map[glx.EntityType]map[string]*template.Template
	entity     map[glx.EntityType]map[string]*template.Template
	partial    map[string]*template.Template
}

func (s templateSet) ForSite() map[string]*template.Template {
	return s.site
}

func (s templateSet) ForType(e glx.EntityType) map[string]*template.Template {
	return s.entityType[e]
}

func (s templateSet) ForEntity(e glx.EntityType) map[string]*template.Template {
	return s.entity[e]
}

func (s templateSet) ForPartials() map[string]*template.Template {
	return s.partial
}

func Templates(templateDir string, specs tmpl.Config, funcs template.FuncMap) (templateSet, error) {
	templates := templateSet{
		site:       make(map[string]*template.Template),
		entityType: make(map[glx.EntityType]map[string]*template.Template),
		entity:     make(map[glx.EntityType]map[string]*template.Template),
	}

	c := templateCompiler{
		templateDir: templateDir,
		funcs:       funcs,
		specs:       specs,
		bases:       make(map[string]*template.Template),
	}

	templates.site = c.compileTemplates("site", specs.Templates.Site)

	for et, typeSpec := range c.specs.Templates.Types {
		templates.entityType[et] = c.compileTemplates("type", typeSpec.Type)
		templates.entity[et] = c.compileTemplates("entity", typeSpec.Entity)
	}

	templates.partial = c.compilePartials(specs.Partials)

	return templates, c.err
}

type templateCompiler struct {
	templateDir string
	funcs       template.FuncMap
	specs       tmpl.Config
	bases       map[string]*template.Template
	err         error
}

func (c *templateCompiler) compileTemplates(context string, specs map[string]tmpl.Spec) map[string]*template.Template {
	if c.err != nil {
		return nil
	}
	tt := make(map[string]*template.Template)
	for name, spec := range specs {
		t, err := c.compileSpec(spec)
		if err != nil {
			c.err = fmt.Errorf("%s template %q: %w", context, name, err)
			return nil
		}
		tt[name] = t
	}
	return tt
}

func (c *templateCompiler) compilePartials(specs map[string]tmpl.FileList) map[string]*template.Template {
	if c.err != nil {
		return nil
	}

	tt := make(map[string]*template.Template)
	for name, files := range specs {
		t, err := c.compileFileList(files)
		if err != nil {
			c.err = fmt.Errorf("partial %q: %w", name, err)
			return nil
		}
		tt[name] = t
	}
	return tt
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

	return c.parseFiles(base, spec.Files)
}

func (c *templateCompiler) compileFileList(files tmpl.FileList) (*template.Template, error) {
	if len(files) == 0 {
		return nil, errors.New("no files")
	}

	name := filepath.Base(files[0])
	return c.parseFiles(template.New(name).Funcs(c.funcs), files)
}

func (c *templateCompiler) parseFiles(t *template.Template, fnames tmpl.FileList) (*template.Template, error) {
	for _, fname := range fnames {
		_, err := t.ParseFiles(filepath.Join(c.templateDir, fname))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
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
		baseSpec, ok := c.specs.Templates.Bases[name]
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
