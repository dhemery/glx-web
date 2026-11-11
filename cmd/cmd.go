// Package cmd implements the glx-web command.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dhemery/glx-web/context"
	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/internal/load"
	"github.com/dhemery/glx-web/internal/render"
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
	Short: "Build a website from a GLX archive and templ ffates.",
	Args:  cobra.NoArgs,
	RunE:  run,
}

var (
	archiveDir  = "./archive"
	templateDir = "./templates"
	outputDir   = "./public"
	staticDir   = "./static"
	clean       = false
)

func init() {
	webCmd.Flags().StringVarP(&archiveDir, "archive", "a", archiveDir,
		"the `dir` of the archive to render")
	webCmd.Flags().BoolVarP(&clean, "clean", "c", clean,
		"remove the existing output dir before building")
	webCmd.Flags().StringVarP(&outputDir, "output", "o", outputDir,
		"the output `dir`")
	webCmd.Flags().StringVarP(&staticDir, "static", "s", staticDir,
		"a `dir` of files to copy into the output dir, if any")
	webCmd.Flags().StringVarP(&templateDir, "templates", "t", templateDir,
		"the template `dir`")
}

func run(_ *cobra.Command, _ []string) error {
	g, err := load.GLX(os.DirFS(archiveDir))
	if err != nil {
		return fmt.Errorf("loading GLX archive: %w", err)
	}

	catalog, err := entity.NewCatalog(g, os.Stderr)
	if err != nil {
		if validationError, ok := errors.AsType[*glx.StructuredValidationError](err); ok {
			for _, e := range validationError.Errors {
				fmt.Fprintln(os.Stderr, e)
			}
		}
		return err
	}

	cfg, err := load.Config(templateDir)
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	state, funcs := context.Create(cfg.Data)

	tmpls, err := load.Templates(templateDir, cfg.Templates, funcs)
	if err != nil {
		return fmt.Errorf("loading templates: %w", err)
	}

	state.SetPartials(tmpls.ForPartials())

	renderer := &render.Renderer{
		Clean:     clean,
		OutputDir: outputDir,
		StaticDir: staticDir,
		Templates: tmpls,
		State:     state,
	}

	return renderer.Render(catalog)
}
