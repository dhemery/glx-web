package load

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/tmpl"
	"go.yaml.in/yaml/v3"
)

func Config(templateDir string) (tmpl.Config, error) {
	configFileName := filepath.Join(templateDir, "config.yaml")
	configFile, err := os.Open(configFileName)
	if err != nil {
		return tmpl.Config{}, err
	}
	defer configFile.Close()

	dec := yaml.NewDecoder(configFile)
	// Complain about unknown fields.
	dec.KnownFields(true)

	var tmplConfig tmpl.Config
	if err = dec.Decode(&tmplConfig); err != nil {
		return tmpl.Config{}, fmt.Errorf("%s: %w", configFileName, err)
	}

	return tmplConfig, nil
}
