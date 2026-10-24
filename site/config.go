package site

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/genealogix/glx/go-glx"
	"go.yaml.in/yaml/v3"
)

// Import for comment links.
var _ = filepath.Match

type Config struct {
	// Data available to templates via the "data" function.
	Data map[string]any `json:"data"`
	// Specification of templates to render the site.
	Templates Spec `json:"templates"`
}

// Spec specifies the templates to render the site.
type Spec struct {
	// Templates that other templates can build on. The key is used as the
	// name of the base. Base templates can build on other base templates.
	Bases map[string]TemplateSpec `json:"bases"`
	// Templates for pages about the site. Each template is applied to the
	// [Config.Data] read from the config.yaml file in the templates
	// directory. The key is used as the base file name for the rendered
	// file.
	Pages map[string]TemplateSpec `json:"pages"`
	// Templates for each entity type. The key is the plural form of the
	// entity type name, e.g. persons.
	Entities map[glx.EntityType]EntityTypeSpec `json:"entities"`
}

// EntityTypeSpec specifies templates for pages about entities and sets
// of entities of a given type.
type EntityTypeSpec struct {
	// Templates to render index pages or other pages about the type as a
	// whole. Each template is applied to the [EntitySet] for the type. The
	// map key is used as the base file name for the rendered file.
	Pages map[string]TemplateSpec `json:"pages"`
	// The template to apply to each entity of the type to render the page
	// about the entity.
	Entity TemplateSpec `json:"entity"`
}

// TemplateSpec specifies how to build a template. If both a base and patterns
// are specified, the template is constructed by parsing the patterns into a
// clone of the base template. If no patterns are listed, the resulting
// template is a clone of the base. If no base is named, the template is
// constructed by parsing the files that match the patterns.
//
// Patterns are applied according to the sementics of [filepath.Match].
type TemplateSpec struct {
	// The name of the base template to build on.
	Base string `json:"base"`
	// The glob patterns that identify the template files to parse,
	// relative to the template directory. If the spec has no Base, the
	// first pattern must supply the template body.
	Patterns []string `json:"patterns"`
}

func loadConfig(dir string) (Config, error) {
	var spec Config

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
