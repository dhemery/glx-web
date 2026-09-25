// Package site renders a website from a compiled GLX archive.
package site

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/web"
	"github.com/genealogix/glx/go-glx"
)

type Renderer struct {
	Archive   *web.Archive
	OutputDir string
	Clean     bool
	StaticDir string
	Templates map[glx.EntityType]*template.Template
	Err       error
}

func (r *Renderer) Render() error {
	if err := os.Mkdir(r.OutputDir, 0755); err != nil {
		return err
	}

	r.renderEntities(glx.EntityTypeAssertions, r.Archive.Assertions)
	r.renderEntities(glx.EntityTypeCitations, r.Archive.Citations)
	r.renderEntities(glx.EntityTypeEvents, r.Archive.Events)
	r.renderEntities(glx.EntityTypeMedia, r.Archive.Media)
	r.renderEntities(glx.EntityTypePersons, r.Archive.Persons)
	r.renderEntities(glx.EntityTypePlaces, r.Archive.Places)
	r.renderEntities(glx.EntityTypeRelationships, r.Archive.Relationships)
	r.renderEntities(glx.EntityTypeRepositories, r.Archive.Repositories)
	r.renderEntities(glx.EntityTypeSources, r.Archive.Sources)

	return r.Err
}

func (r *Renderer) renderEntities[T any](t glx.EntityType, entities map[string]*T) {
	if r.Err != nil {
		return
	}
	typeDir := filepath.Join(r.OutputDir, t.Plural())
	if err := os.Mkdir(typeDir, 0755); err != nil {
		r.Err = err
		return
	}

	for id, e := range entities {
		entityDir := filepath.Join(typeDir, id)
		if err := os.Mkdir(entityDir, 0755); err != nil {
			r.Err = err
			return
		}

		if err := render(entityDir, e, r.Templates[t]); err != nil {
			r.Err = err
			return
		}
	}
}

func render(dir string, data any, tmpl *template.Template) error {
	if data == nil {
		return nil
	}

	fname := filepath.Join(dir, "index.html")
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
