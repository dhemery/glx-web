// Package cmd implements the glx-web command.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/render"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

func Execute() {
	err := webCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var webCmd = &cobra.Command{
	Use:   "glx-web",
	Short: "Build a website from a GLX archive and templates.",
	Args:  cobra.NoArgs,
	RunE:  run,
}

var (
	archiveDir          = "./archive"
	cleanBeforeBuilding = false
	outputDir           = "./public"
	staticDir           = "./static"
	templateDir         = "./templates"
)

func init() {
	webCmd.Flags().StringVarP(&archiveDir, "archive", "a", archiveDir,
		"the `dir` of the archive to render")
	webCmd.Flags().BoolVarP(&cleanBeforeBuilding, "clean", "c", cleanBeforeBuilding,
		"remove the existing output dir before building")
	webCmd.Flags().StringVarP(&outputDir, "output", "o", outputDir,
		"the output `dir`")
	webCmd.Flags().StringVarP(&staticDir, "static", "s", staticDir,
		"a `dir` of files to copy into the output dir, if any")
	webCmd.Flags().StringVarP(&templateDir, "templates", "t", templateDir,
		"the template `dir`")

}

func run(_ *cobra.Command, _ []string) error {
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

	templates, err := render.Templates(absTemplateDir)
	if err != nil {
		if validationError, ok := errors.AsType[*glx.StructuredValidationError](err); ok {
			for _, e := range validationError.Errors {
				fmt.Fprintln(os.Stderr, e)
			}
		}
		return err
	}

	if cleanBeforeBuilding {
		if err := os.RemoveAll(outputDir); err != nil {
			return fmt.Errorf("removing output dir: %w", err)
		}
	}

	if err := os.Mkdir(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	if staticDir != "" {
		if err := os.CopyFS(outputDir, os.DirFS(staticDir)); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("copying static files: %w", err)
			}
		}
	}

	r := &render.Renderer{
		OutputDir: absOutputDir,
		Templates: templates,
	}

	catalog, err := entity.NewCatalog(absArchiveDir, os.Stderr)
	if err != nil {
		return fmt.Errorf("loading GLX archive: %w", err)
	}

	if err := r.Render(catalog); err != nil {
		return fmt.Errorf("rendering: %w", err)
	}
	return nil
}
