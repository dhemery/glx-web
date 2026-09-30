// Package site renders a website from a compiled GLX archive.
package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/web"
	"github.com/genealogix/glx/go-glx"
)

type Renderer struct {
	OutputDir string
	Clean     bool
	StaticDir string
	Templates map[glx.EntityType]*template.Template
	Err       error
}

func (r *Renderer) Render(archive *web.Archive, glxfile *glx.GLXFile) error {
	if err := os.Mkdir(r.OutputDir, 0755); err != nil {
		return err
	}

	data := web.Data{
		Archive: archive,
		GLX:     glxfile,
		Content: nil,
	}

	r.renderEntities(glx.EntityTypeAssertions, archive.Assertions, data)
	r.renderEntities(glx.EntityTypeCitations, archive.Citations, data)
	r.renderEntities(glx.EntityTypeEvents, archive.Events, data)
	r.renderEntities(glx.EntityTypeMedia, archive.Media, data)
	r.renderEntities(glx.EntityTypePersons, archive.Persons, data)
	r.renderEntities(glx.EntityTypePlaces, archive.Places, data)
	r.renderEntities(glx.EntityTypeRelationships, archive.Relationships, data)
	r.renderEntities(glx.EntityTypeRepositories, archive.Repositories, data)
	r.renderEntities(glx.EntityTypeSources, archive.Sources, data)

	return r.Err
}

func (r *Renderer) renderEntities[T any](t glx.EntityType, entities map[string]*T, data web.Data) {
	if r.Err != nil {
		return
	}

	typeDir := filepath.Join(r.OutputDir, t.Plural())
	if err := os.Mkdir(typeDir, 0755); err != nil {
		r.Err = err
		return
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
