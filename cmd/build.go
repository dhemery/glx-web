package cmd

import (
	"html/template"
	"path/filepath"

	"github.com/dhemery/glx-web/archive"
	"github.com/dhemery/glx-web/site"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a website from a GLX archive",
	RunE:  runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuild(c *cobra.Command, args []string) error {
	archivePath, err := c.Flags().GetString("archive")
	if err != nil {
		return err
	}

	wa, err := archive.Load(archivePath)

	templates, err := loadTemplates()
	if err != nil {
		return err
	}

	return site.Render(wa, templates)
}

func loadTemplates() (map[string]*template.Template, error) {
	const tmplDir = "templates"
	baseTemplate, err := template.ParseFiles(filepath.Join(tmplDir, "base.gotmpl"))
	if err != nil {
		return nil, err
	}

	entityTemplates := map[string]*template.Template{}

	for _, t := range glx.AllEntityTypes {
		glob := filepath.Join(tmplDir, t.Plural(), "*.gotmpl")

		et, err := baseTemplate.Clone()
		if err != nil {
			return nil, err
		}

		et, err = et.ParseGlob(glob)
		if err != nil {
			// Ignore failures for now
			continue
		}

		entityTemplates[t.Plural()] = et
	}

	return entityTemplates, nil
}
