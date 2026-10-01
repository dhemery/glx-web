// Package site renders a website from a compiled GLX archive.
package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

type Renderer struct {
	OutputDir string
	Clean     bool
	StaticDir string
	Templates map[glx.EntityType]*template.Template
	Err       error
}

func (r *Renderer) Render(catalog *entity.Catalog, glxfile *glx.GLXFile) error {
	if err := os.Mkdir(r.OutputDir, 0755); err != nil {
		return err
	}

	r.renderEntities(glx.EntityTypeAssertions, catalog.Assertions, catalog, glxfile)
	r.renderEntities(glx.EntityTypeCitations, catalog.Citations, catalog, glxfile)
	r.renderEntities(glx.EntityTypeEvents, catalog.Events, catalog, glxfile)
	r.renderEntities(glx.EntityTypeMedia, catalog.Media, catalog, glxfile)
	r.renderEntities(glx.EntityTypePersons, catalog.Persons, catalog, glxfile)
	r.renderEntities(glx.EntityTypePlaces, catalog.Places, catalog, glxfile)
	r.renderEntities(glx.EntityTypeRelationships, catalog.Relationships, catalog, glxfile)
	r.renderEntities(glx.EntityTypeRepositories, catalog.Repositories, catalog, glxfile)
	r.renderEntities(glx.EntityTypeSources, catalog.Sources, catalog, glxfile)

	return r.Err
}

func (r *Renderer) renderEntities[T any](t glx.EntityType, entities map[string]*T, catalog *entity.Catalog, glxfile *glx.GLXFile) {
	if r.Err != nil {
		return
	}

	typeDir := filepath.Join(r.OutputDir, t.Plural())
	if err := os.Mkdir(typeDir, 0755); err != nil {
		r.Err = err
		return
	}

	data := entity.Data{
		Catalog: catalog,
		GLX:     glxfile,
	}
	entityTemplate := r.Templates[t]

	for id, entity := range entities {
		entityDir := filepath.Join(typeDir, id)
		if err := os.Mkdir(entityDir, 0755); err != nil {
			r.Err = fmt.Errorf("creating dir for %s[%s]: %w", t, id, err)
			return
		}

		data.Content = entity
		fname := filepath.Join(entityDir, "index.html")

		if err := render(fname, data, entityTemplate); err != nil {
			r.Err = fmt.Errorf("rendering %s[%s]: %w", t, id, err)
			return
		}
	}
}

func render(fname string, data any, tmpl *template.Template) error {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
