package cmd

import (
	"fmt"
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
	absArchiveDir, err := filepath.Abs(archiveDir)
	if err != nil {
		return fmt.Errorf("archive directory: %w", err)
	}

	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("output directory: %w", err)
	}

	absStaticDir, err := filepath.Abs(staticDir)
	if err != nil {
		return fmt.Errorf("static directory: %w", err)
	}

	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		return fmt.Errorf("template directory: %w", err)
	}

	templates, err := loadTemplates(absTemplateDir)
	if err != nil {
		return err
	}

	a, err := archive.Load(absArchiveDir)
	if err != nil {
		return err
	}

	r := &site.Renderer{
		Archive:   a,
		OutputDir: absOutputDir,
		Clean:     cleanOutput,
		StaticDir: absStaticDir,
		Templates: templates,
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
