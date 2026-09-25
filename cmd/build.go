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

var (
	archivePath     = "."
	overwriteOutput = false
	includeLiving   = false
	outputPath      = "./public"
	staticPath      = "./static"
	templatePath    = "./templates"
)

func init() {
	buildCmd.Flags().StringVarP(&archivePath, "archive", "a", archivePath,
		"Archive `dir`.")
	buildCmd.Flags().BoolVarP(&overwriteOutput, "force", "f", overwriteOutput,
		"Overwrite existing output directory.")
	buildCmd.Flags().BoolVarP(&includeLiving, "living", "l", includeLiving,
		"Include living people.")
	buildCmd.Flags().StringVarP(&outputPath, "output", "o", outputPath,
		"Output `dir`.")
	buildCmd.Flags().StringVarP(&staticPath, "static", "s", staticPath,
		"Static `dir` of files to copy verbatim into the output.")
	buildCmd.Flags().StringVarP(&templatePath, "templates", "t", templatePath,
		"Template `dir`.")

	rootCmd.AddCommand(buildCmd)
}

func runBuild(c *cobra.Command, args []string) error {
	archivePath, err := c.Flags().GetString("archive")
	if err != nil {
		return err
	}

	wa, err := archive.Load(archivePath)
	if err != nil {
		return err
	}

	templates, err := loadTemplates()
	if err != nil {
		return err
	}

	return site.Render(wa, templates)
}

func loadTemplates() (map[glx.EntityType]*template.Template, error) {
	const tmplDir = "templates"
	baseTemplate, err := template.ParseFiles(filepath.Join(tmplDir, "base.gotmpl"))
	if err != nil {
		return nil, err
	}

	entityTemplates := map[glx.EntityType]*template.Template{}

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

		entityTemplates[t] = et
	}

	return entityTemplates, nil
}
