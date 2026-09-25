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
	Args:  cobra.NoArgs,
	RunE:  runBuild,
}

var (
	archiveDir    = "."
	cleanOutput   = false
	includeLiving = false
	outputDir     = "./public"
	staticDir     = "./static"
	templateDir   = "./templates"
)

func init() {
	buildCmd.Flags().StringVarP(&archiveDir, "archive", "a", archiveDir,
		"archive `dir`")
	buildCmd.Flags().BoolVarP(&cleanOutput, "clean", "c", cleanOutput,
		"remove existing output directory before building")
	buildCmd.Flags().BoolVarP(&includeLiving, "living", "l", includeLiving,
		"include living people")
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", outputDir,
		"output `dir`")
	buildCmd.Flags().StringVarP(&staticDir, "static", "s", staticDir,
		"static `dir` of files to copy into the output dir")
	buildCmd.Flags().StringVarP(&templateDir, "templates", "t", templateDir,
		"template `dir`")

	rootCmd.AddCommand(buildCmd)
}

func runBuild(_ *cobra.Command, _ []string) error {
	templates, err := loadTemplates(templateDir)
	if err != nil {
		return err
	}

	a, err := archive.Load(archiveDir)
	if err != nil {
		return err
	}

	r := &site.Renderer{
		Archive:   a,
		OutputDir: outputDir,
		Clean:     cleanOutput,
		StaticDir: staticDir,
		Templates: templates,
		Err:       err,
	}
	return r.Render()
}

func loadTemplates(templateDir string) (map[glx.EntityType]*template.Template, error) {
	baseTemplate, err := template.ParseFiles(filepath.Join(templateDir, "base.gotmpl"))
	if err != nil {
		return nil, err
	}

	entityTemplates := map[glx.EntityType]*template.Template{}

	for _, t := range glx.AllEntityTypes {
		glob := filepath.Join(templateDir, t.Plural(), "*.gotmpl")

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
