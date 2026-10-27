// Package render renders a website from a compiled GLX archive.
package render

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/tmpl"
	"github.com/genealogix/glx/go-glx"
)

type templateSet interface {
	ForSite() map[string]*template.Template
	ForType(glx.EntityType) map[string]*template.Template
	ForEntity(glx.EntityType) map[string]*template.Template
}

type Options struct {
	Clean     bool
	Data      map[string]any
	OutputDir string
	StaticDir string
	Templates templateSet
}

type renderer struct {
	Templates templateSet
	OutputDir string
	Err       error
}

var (
	site      tmpl.Site
	pageStore = make(tmpl.Store)
)

// TemplateFuncs is a map of functions to add to each template.
var TemplateFuncs = map[string]any{
	"site":  func() tmpl.Site { return site },
	"store": func() tmpl.Store { return pageStore },
}

func Site(catalog *entity.Catalog, opts *Options) error {
	site = tmpl.Site{
		Data:  opts.Data,
		Store: make(tmpl.Store),
	}

	r := renderer{
		Templates: opts.Templates,
		OutputDir: opts.OutputDir,
	}

	if opts.Clean {
		if err := os.RemoveAll(opts.OutputDir); err != nil {
			return fmt.Errorf("removing output dir: %w", err)
		}
	}

	if err := os.Mkdir(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	if opts.StaticDir != "" {
		if err := os.CopyFS(opts.OutputDir, os.DirFS(opts.StaticDir)); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("copying static files: %w", err)
			}
		}
	}

	for name, tmpl := range r.Templates.ForSite() {
		fname := filepath.Join(r.OutputDir, name)
		err := render(fname, site, tmpl)
		if err != nil {
			return fmt.Errorf("rendering %s: %w", fname, err)
		}
	}

	r.renderEntities(glx.EntityTypeAssertions, catalog.Assertions())
	r.renderEntities(glx.EntityTypeCitations, catalog.Citations())
	r.renderEntities(glx.EntityTypeEvents, catalog.Events())
	r.renderEntities(glx.EntityTypeMedia, catalog.Media())
	r.renderEntities(glx.EntityTypePersons, catalog.Persons())
	r.renderEntities(glx.EntityTypePlaces, catalog.Places())
	r.renderEntities(glx.EntityTypeRelationships, catalog.Relationships())
	r.renderEntities(glx.EntityTypeRepositories, catalog.Repositories())
	r.renderEntities(glx.EntityTypeSources, catalog.Sources())

	return r.Err
}

type ider interface {
	ID() string
}

func (r *renderer) renderEntities[S ~[]E, E ider](et glx.EntityType, entities S) {
	if r.Err != nil {
		return
	}

	entitySet := tmpl.EntitySet{
		EntityType: et,
		Entities:   entities,
	}

	for name, t := range r.Templates.ForType(et) {
		fname := filepath.Join(r.OutputDir, name)
		err := render(fname, entitySet, t)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s: %w", fname, err)
			return
		}
	}

	for _, entity := range entities {
		id := entity.ID()
		for name, t := range r.Templates.ForEntity(et) {
			fname := filepath.Join(r.OutputDir, strings.ReplaceAll(name, "<id>", id))

			err := render(fname, entity, t)
			if err != nil {
				r.Err = fmt.Errorf("rendering %s: %w", fname, err)
				return
			}
		}
	}
}

func render(fname string, dot any, tmpl *template.Template) error {
	if err := os.MkdirAll(filepath.Dir(fname), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	clear(pageStore)

	return tmpl.Execute(f, dot)
}
