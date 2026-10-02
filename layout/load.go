// Package layout loads and supplies templates for entity and index pages.
package layout

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/genealogix/glx/go-glx"
	"gopkg.in/yaml.v3"
)

type SiteTemplates struct {
	// Templates for each entity type.
	Entities map[glx.EntityType]EntityTemplates
	// Templates for top-level files.
	Indexes map[string]*template.Template
}

type EntityTemplates struct {
	// The template for entities of the ssociated type.
	Entity *template.Template
	// Templates for indexes for each entity type.
	Indexes map[string]*template.Template
}

// TemplateSpec specifies how to build a template. If a base is named, the
// template is constructed by parsing the patterns into a clone of the named
// base template. Otherwise the patterns are parsed into a new, empty template.
type TemplateSpec struct {
	// The name of the base template, if any.
	Base string `json:"base"`
	// The glob patterns of the files to parse to build the template.
	// The patterns should be relative to the template directory.
	Patterns []string `json:"patterns"`
}

// EntityTypeTemplateSpec specifies tne entity and index templates for an
// entity type.
type EntityTypeTemplateSpec struct {
	// Specification of the template to apply to each entity of the
	// associated type.
	Entity TemplateSpec `json:"entity"`
	// Specification of the index templates for the associated entity type.
	// The .Content for each index template is the map of entities keyed by
	// ID.
	Indexes map[string]TemplateSpec `json:"indexes"`
}

type SiteTemplateSpec struct {
	// Specifications of the base templates that the other templates build on.
	Bases map[string]TemplateSpec `json:"bases"`
	// Specifications of templates for top-level pages. The .Content for
	// each top-level template is nil.
	Indexes map[string]TemplateSpec `json:"indexes"`
	// Specifications of the templates for each entity type.
	Entities map[glx.EntityType]EntityTypeTemplateSpec `json:"entities"`
}

const (
	templateExtensionDefault    = "gotmpl"
	templatePatternDefaultRoot  = "base/root." + templateExtensionDefault
	templatePatternDefaultList  = "base/list." + templateExtensionDefault
	templatePatternDefaultIndex = "base/index." + templateExtensionDefault
)

func newSpec() SiteTemplateSpec {
	return SiteTemplateSpec{
		Bases:    map[string]TemplateSpec{},
		Indexes:  map[string]TemplateSpec{},
		Entities: map[glx.EntityType]EntityTypeTemplateSpec{},
	}
}

func Load(templateDir string) (SiteTemplates, error) {
	spec, err := loadSpec(templateDir)
	if err != nil {
		return SiteTemplates{}, err
	}

	return newSiteTemplates(os.DirFS(templateDir), spec)
}

type loader struct {
	fsys          fs.FS
	baseSpecs     map[string]TemplateSpec
	baseTemplates map[string]*template.Template
}

func newSiteTemplates(fsys fs.FS, siteSpec SiteTemplateSpec) (SiteTemplates, error) {
	st := SiteTemplates{
		Entities: map[glx.EntityType]EntityTemplates{},
		Indexes:  map[string]*template.Template{},
	}

	l := loader{
		fsys:          fsys,
		baseSpecs:     siteSpec.Bases,
		baseTemplates: make(map[string]*template.Template),
	}

	for name, indexSpec := range siteSpec.Indexes {
		t, err := l.load(indexSpec)
		if err != nil {
			return st, fmt.Errorf("loading top-level template %q: %w", name, err)
		}
		st.Indexes[name] = t
	}

	return st, nil
}

func (l *loader) load(spec TemplateSpec) (*template.Template, error) {
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

func (l *loader) cloneBase(name string) (*template.Template, error) {
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

func loadSpec(dir string) (SiteTemplateSpec, error) {
	var spec SiteTemplateSpec

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
