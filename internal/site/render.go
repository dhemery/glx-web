// Package site renders a website from a compiled GLX archive.
package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/layout"
	"github.com/genealogix/glx/go-glx"
)

type Renderer struct {
	OutputDir string
	Clean     bool
	StaticDir string
	Templates layout.SiteTemplates
	Err       error
}

func (r *Renderer) Render(catalog *entity.Catalog, glxfile *glx.GLXFile) error {
	if err := os.Mkdir(r.OutputDir, 0755); err != nil {
		return err
	}

	err := r.renderTopLevelIndexes(catalog, glxfile)
	if err != nil {
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

func (r *Renderer) renderTopLevelIndexes(catalog *entity.Catalog, glxfile *glx.GLXFile) error {
	data := entity.Data{
		GLX:     glxfile,
		Catalog: catalog,
		Content: nil,
	}

	for name, tmpl := range r.Templates.Indexes {
		fname := filepath.Join(r.OutputDir, name+".html")
		err := render(fname, data, tmpl)
		if err != nil {
			return fmt.Errorf("rendering top-level template %q: %w", name, err)
		}
	}

	return nil
}

func (r *Renderer) renderEntities[T any](_ glx.EntityType, _ map[string]*T, _ *entity.Catalog, _ *glx.GLXFile) {
	if r.Err != nil {
		return
	}
}

func render(fname string, data any, tmpl *template.Template) error {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
