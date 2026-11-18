package load

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/tmpl"
	"go.yaml.in/yaml/v4"
)

func Config(templateDir string) (tmpl.Config, error) {
	configFileName := filepath.Join(templateDir, "config.yaml")
	configFile, err := os.Open(configFileName)
	if err != nil {
		return tmpl.Config{}, err
	}
	defer configFile.Close()

	loader, err := yaml.NewLoader(configFile, yaml.WithKnownFields())
	if err != nil {
		return tmpl.Config{}, err
	}

	var tmplConfig tmpl.Config
	if err = loader.Load(&tmplConfig); err != nil {
		return tmpl.Config{}, fmt.Errorf("%s: %w", configFileName, err)
	}

	return tmplConfig, nil
}
