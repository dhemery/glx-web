// Package site renders a website from a compiled GLX archive.
package site

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

type Options struct {
	Clean       bool
	OutputDir   string
	StaticDir   string
	TemplateDir string
}

type templateSet interface {
	ForSite() map[string]*template.Template
	ForEntitySet(glx.EntityType) map[string]*template.Template
	ForEntity(glx.EntityType) *template.Template
}

type renderer struct {
	OutputDir string
	Templates templateSet
	Err       error
}

func Render(catalog *entity.Catalog, opts *Options) error {
	r := renderer{OutputDir: opts.OutputDir}
	var err error

	if r.Templates, err = loadTemplates(opts.TemplateDir); err != nil {
		return err
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

	if err = r.renderTopLevelPages(); err != nil {
		return err
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

func (r *renderer) renderTopLevelPages() error {
	dot := Site()
	for name, tmpl := range r.Templates.ForSite() {
		fname := filepath.Join(r.OutputDir, name+".html")
		err := render(fname, dot, tmpl)
		if err != nil {
			return fmt.Errorf("rendering top-level page with template %q: %w", name, err)
		}
	}

	return r.Err
}

type ider interface {
	ID() string
}

func (r *renderer) renderEntities[S ~[]E, E ider](et glx.EntityType, entities S) {
	if r.Err != nil {
		return
	}

	entityTypeDir := filepath.Join(r.OutputDir, et.Plural())

	if err := os.Mkdir(entityTypeDir, 0755); err != nil {
		r.Err = fmt.Errorf("rendering %s: %w", et, err)
		return
	}

	entitySet := EntitySet{
		EntityType: et,
		Entities:   entities,
	}

	for name, tmpl := range r.Templates.ForEntitySet(et) {
		fname := filepath.Join(entityTypeDir, name+".html")
		err := render(fname, entitySet, tmpl)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s with page template %q: %w", et, name, err)
			return
		}
	}

	entityTemplate := r.Templates.ForEntity(et)
	if entityTemplate == nil {
		r.Err = fmt.Errorf("no entity template for %s", et)
		return
	}

	for _, entity := range entities {
		id := entity.ID()
		entityDir := filepath.Join(entityTypeDir, id)
		if err := os.Mkdir(entityDir, 0755); err != nil {
			r.Err = fmt.Errorf("rendering %s[%s]: %w", et, id, err)
			return
		}

		fname := filepath.Join(entityDir, "index.html")

		err := render(fname, entity, entityTemplate)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s[%s]: %w", et, id, err)
			return
		}
	}
}

func render(fname string, dot any, tmpl *template.Template) error {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	ResetPage()

	return tmpl.Execute(f, dot)
}
