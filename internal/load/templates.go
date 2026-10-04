package load

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/config"
	"github.com/dhemery/glx-web/internal/render"
	"github.com/genealogix/glx/go-glx"
	"gopkg.in/yaml.v3"
)

func Templates(templateDir string) (render.SiteTemplates, error) {
	spec, err := loadSpec(templateDir)
	if err != nil {
		return render.SiteTemplates{}, err
	}

	return newSiteTemplates(os.DirFS(templateDir), spec)
}

type templateLoader struct {
	fsys          fs.FS
	baseSpecs     map[string]config.TemplateSpec
	baseTemplates map[string]*template.Template
}

func newSiteTemplates(fsys fs.FS, siteSpec config.SiteTemplates) (render.SiteTemplates, error) {
	st := render.SiteTemplates{
		EntityTypes: map[glx.EntityType]render.EntityTypeTemplates{},
		Pages:       map[string]*template.Template{},
	}

	l := templateLoader{
		fsys:          fsys,
		baseSpecs:     siteSpec.Bases,
		baseTemplates: make(map[string]*template.Template),
	}

	for name, indexSpec := range siteSpec.Pages {
		t, err := l.load(indexSpec)
		if err != nil {
			return st, fmt.Errorf("loading top-level template %q: %w", name, err)
		}
		st.Pages[name] = t
	}

	for et, entityTypeSpec := range siteSpec.Entities {
		templates, err := l.loadEntityTypeTemplates(entityTypeSpec)
		if err != nil {
			return st, fmt.Errorf("loading %s templates: %w", et, err)
		}
		st.EntityTypes[et] = templates
	}

	return st, nil
}

func (l templateLoader) loadEntityTypeTemplates(spec config.EntityTypeTemplates) (render.EntityTypeTemplates, error) {
	et := render.EntityTypeTemplates{
		Pages: make(map[string]*template.Template),
	}

	for name, indexSpec := range spec.Pages {
		t, err := l.load(indexSpec)
		if err != nil {
			return et, fmt.Errorf("index template %q: %w", name, err)
		}
		et.Pages[name] = t
	}
	entityTemplate, err := l.load(spec.Entity)
	if err != nil {
		return et, fmt.Errorf("entity template: %w", err)
	}
	et.Entity = entityTemplate
	return et, nil

}

func (l *templateLoader) load(spec config.TemplateSpec) (*template.Template, error) {
	base, err := l.cloneBase(spec.Base)
	if err != nil {
		return nil, err
	}

	if base == nil {
		return template.ParseFS(l.fsys, spec.Patterns...)
	}

	if len(spec.Patterns) == 0 {
		return base, nil
	}

	return base.ParseFS(l.fsys, spec.Patterns...)
}

func (l *templateLoader) cloneBase(name string) (*template.Template, error) {
	if name == "" {
		return nil, nil
	}

	var base *template.Template
	var err error
	base, ok := l.baseTemplates[name]
	if !ok {
		spec, ok := l.baseSpecs[name]
		if !ok {
			return nil, fmt.Errorf("no spec for base template %q", name)
		}
		base, err = l.load(spec)
		if err != nil {
			return nil, fmt.Errorf("loading base template %q: %w", name, err)
		}
		l.baseTemplates[name] = base
	}

	return base.Clone()
}

func loadSpec(dir string) (config.SiteTemplates, error) {
	var spec config.SiteTemplates

	configFileName := filepath.Join(dir, "config.yaml")
	configFile, err := os.Open(configFileName)
	if err != nil {
		return spec, err
	}
	defer configFile.Close()

	dec := yaml.NewDecoder(configFile)
	dec.KnownFields(true)

	err = dec.Decode(&spec)
	if err != nil {
		return spec, fmt.Errorf("%s: %w", "config.yaml", err)
	}
	return spec, err
}
