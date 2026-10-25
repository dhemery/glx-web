package load

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/tmpl"
	"go.yaml.in/yaml/v3"
)

type cfg struct {
	data      map[string]any
	templates templateSet
}

func (c cfg) Data() map[string]any {
	return c.data
}

func (c cfg) Templates() templateSet {
	return c.templates
}

func Config(dir string, funcs template.FuncMap) (cfg, error) {
	var tmplConfig tmpl.Config

	configFileName := filepath.Join(dir, "config.yaml")
	configFile, err := os.Open(configFileName)
	if err != nil {
		return cfg{}, err
	}
	defer configFile.Close()

	dec := yaml.NewDecoder(configFile)
	dec.KnownFields(true)

	err = dec.Decode(&tmplConfig)
	if err != nil {
		return cfg{}, fmt.Errorf("%s: %w", "config.yaml", err)
	}

	tt, err := compileTemplates(os.DirFS(dir), tmplConfig.Templates, funcs)
	if err != nil {
		return cfg{}, fmt.Errorf("loading templates: %w", err)
	}

	return cfg{data: tmplConfig.Data, templates: tt}, err
}
