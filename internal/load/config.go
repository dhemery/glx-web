package load

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/tmpl"
	"go.yaml.in/yaml/v3"
)

type config struct {
	data      map[string]any
	templates templateSet
}

func (c config) Data() map[string]any {
	return c.data
}

func (c config) Templates() templateSet {
	return c.templates
}

func Config(templateDir string, funcs template.FuncMap) (config, error) {
	configFileName := filepath.Join(templateDir, "config.yaml")
	configFile, err := os.Open(configFileName)
	if err != nil {
		return config{}, err
	}
	defer configFile.Close()

	dec := yaml.NewDecoder(configFile)
	// Complain about unknown fields.
	dec.KnownFields(true)

	var tmplConfig tmpl.Config
	if err = dec.Decode(&tmplConfig); err != nil {
		return config{}, fmt.Errorf("%s: %w", configFileName, err)
	}

	tt, err := compileTemplateSet(templateDir, tmplConfig.Templates, funcs)
	if err != nil {
		return config{}, err
	}

	return config{data: tmplConfig.Data, templates: tt}, err
}
