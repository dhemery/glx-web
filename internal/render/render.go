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

// state captures the global state available to each template via dot or
// functions.
type state interface {
	Site() tmpl.Site
	PrepareExecute()
}

type Renderer struct {
	Clean     bool
	OutputDir string
	StaticDir string
	Templates templateSet
	State     state
	Err       error
}

func (r *Renderer) Render(catalog *entity.Catalog) error {
	if r.Clean {
		if err := os.RemoveAll(r.OutputDir); err != nil {
			return fmt.Errorf("removing output dir: %w", err)
		}
	}

	if err := os.Mkdir(r.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	if r.StaticDir != "" {
		if err := os.CopyFS(r.OutputDir, os.DirFS(r.StaticDir)); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("copying static files: %w", err)
			}
		}
	}

	for name, tmpl := range r.Templates.ForSite() {
		fname := filepath.Join(r.OutputDir, name)
		err := r.executeTemplate(fname, r.State.Site(), tmpl)
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

func (r *Renderer) renderEntities[S ~[]E, E ider](et glx.EntityType, entities S) {
	if r.Err != nil {
		return
	}

	entitySet := tmpl.EntitySet{
		EntityType: et,
		Entities:   entities,
	}

	for name, t := range r.Templates.ForType(et) {
		fname := filepath.Join(r.OutputDir, name)
		err := r.executeTemplate(fname, entitySet, t)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s: %w", fname, err)
			return
		}
	}

	for _, entity := range entities {
		id := entity.ID()
		for name, t := range r.Templates.ForEntity(et) {
			fname := filepath.Join(r.OutputDir, strings.ReplaceAll(name, "<id>", id))

			err := r.executeTemplate(fname, entity, t)
			if err != nil {
				r.Err = fmt.Errorf("rendering %s: %w", fname, err)
				return
			}
		}
	}
}

func (r *Renderer) executeTemplate(fname string, dot any, tmpl *template.Template) error {
	if err := os.MkdirAll(filepath.Dir(fname), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	r.State.PrepareExecute()

	return tmpl.Execute(f, dot)
}
