// Package cmd implements the glx-web command.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/internal/load"
	"github.com/dhemery/glx-web/site"
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
	archiveDir = "./archive"
	opts       = &site.Options{
		Clean:       false,
		OutputDir:   "./public",
		StaticDir:   "./static",
		TemplateDir: "./templates",
	}
)

func init() {
	webCmd.Flags().StringVarP(&archiveDir, "archive", "a", archiveDir,
		"the `dir` of the archive to render")
	webCmd.Flags().BoolVarP(&opts.Clean, "clean", "c", opts.Clean,
		"remove the existing output dir before building")
	webCmd.Flags().StringVarP(&opts.OutputDir, "output", "o", opts.OutputDir,
		"the output `dir`")
	webCmd.Flags().StringVarP(&opts.StaticDir, "static", "s", opts.StaticDir,
		"a `dir` of files to copy into the output dir, if any")
	webCmd.Flags().StringVarP(&opts.TemplateDir, "templates", "t", opts.TemplateDir,
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

	if err := site.Render(catalog, opts); err != nil {
		return fmt.Errorf("rendering: %w", err)
	}
	return nil
}
