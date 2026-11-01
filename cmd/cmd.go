// Package cmd implements the glx-web command.
package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"os"

	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/internal/load"
	"github.com/dhemery/glx-web/internal/render"
	"github.com/dhemery/glx-web/tmpl"
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

type state struct {
	site  tmpl.Site
	store tmpl.Store
}

func (s *state) Site() tmpl.Site {
	return s.site
}

func (s *state) Store() tmpl.Store {
	return s.store
}

func (s *state) PrepareExecute() {
	s.store = make(tmpl.Store)
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

	site := tmpl.Site{Data: cfg.Data, Store: make(tmpl.Store)}
	state := &state{site: site}

	funcs := template.FuncMap{
		"is_entity":           entity.IsEntity,
		"is_property":         entity.IsProperty,
		"is_vocabulary_value": entity.IsVocabularyValue,
		"site":                state.Site,
		"store":               state.Store,
	}

	tmpls, err := load.Templates(templateDir, cfg.Templates, funcs)
	if err != nil {
		return fmt.Errorf("loading templates: %w", err)
	}

	renderer := &render.Renderer{
		Clean:     clean,
		OutputDir: outputDir,
		StaticDir: staticDir,
		Templates: tmpls,
		State:     state,
	}

	return renderer.Render(catalog)
}
