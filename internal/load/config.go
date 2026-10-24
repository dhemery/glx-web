package load

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/config"
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
	var spec config.Config

	configFileName := filepath.Join(dir, "config.yaml")
	configFile, err := os.Open(configFileName)
	if err != nil {
		return cfg{}, err
	}
	defer configFile.Close()

	dec := yaml.NewDecoder(configFile)
	dec.KnownFields(true)

	err = dec.Decode(&spec)
	if err != nil {
		return cfg{}, fmt.Errorf("%s: %w", "config.yaml", err)
	}

	tt, err := compileTemplates(os.DirFS(dir), spec.Templates, funcs)
	if err != nil {
		return cfg{}, fmt.Errorf("loading templates: %w", err)
	}

	return cfg{data: spec.Data, templates: tt}, err
}
