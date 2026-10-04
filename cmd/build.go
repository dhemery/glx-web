package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/internal/load"
	"github.com/dhemery/glx-web/internal/render"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a website from a GLX archive",
	Args:  cobra.NoArgs,
	RunE:  runBuild,
}

var (
	archiveDir  = "."
	outputDir   = "./public"
	templateDir = "./templates"
)

func init() {
	buildCmd.Flags().StringVarP(&archiveDir, "archive", "a", archiveDir,
		"archive `dir`")
	// buildCmd.Flags().BoolVarP(&cleanOutput, "clean", "c", cleanOutput,
	// 	"remove existing output directory before building")
	// buildCmd.Flags().BoolVarP(&includeLiving, "living", "l", includeLiving,
	// 	"include living people")
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", outputDir,
		"output `dir`")
	// buildCmd.Flags().StringVarP(&staticDir, "static", "s", staticDir,
	// 	"static `dir` of files to copy into the output dir")
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

	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		return fmt.Errorf("template directory: %w", err)
	}

	templates, err := load.Templates(absTemplateDir)
	if err != nil {
		return err
	}

	glxFile, err := load.Archive(absArchiveDir)
	if err != nil {
		return err
	}

	r := &render.Renderer{
		OutputDir: absOutputDir,
		Templates: templates,
	}

	return r.Render(entity.NewCatalog(glxFile), glxFile)
}
