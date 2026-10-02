// Package site renders a website from a compiled GLX archive.
package site

import (
	"io"
	"os"

	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/template"
	"github.com/genealogix/glx/go-glx"
)

type Renderer struct {
	OutputDir string
	Clean     bool
	StaticDir string
	Templates template.SiteTemplates
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
}

type exec interface {
	Execute(w io.Writer, data any) error
}

func render(fname string, data any, tmpl exec) error {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
