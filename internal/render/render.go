// Package render renders a website from a compiled GLX archive.
package render

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/data"
	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

type SiteTemplates struct {
	// Templates to render top-level pages.
	Pages map[string]*template.Template
	// Templates for each entity type.
	EntityTypes map[glx.EntityType]EntityTypeTemplates
}

type EntityTypeTemplates struct {
	// Templates to render pages about the entity type.
	Pages map[string]*template.Template
	// Template to render the page about an entity of the type.
	Entity *template.Template
}

type Renderer struct {
	OutputDir string
	Clean     bool
	StaticDir string
	Templates SiteTemplates
	Err       error
}

func (r *Renderer) Render(catalog *entity.Catalog, glxfile *glx.GLXFile) error {
	data := data.Site{
		GLX:     glxfile,
		Catalog: catalog,
	}

	err := r.renderTopLevelPages(data)
	if err != nil {
		return err
	}

	r.renderEntities(glx.EntityTypeAssertions, catalog.Assertions(), data)
	r.renderEntities(glx.EntityTypeCitations, catalog.Citations(), data)
	r.renderEntities(glx.EntityTypeEvents, catalog.Events(), data)
	r.renderEntities(glx.EntityTypeMedia, catalog.Media(), data)
	r.renderEntities(glx.EntityTypePersons, catalog.Persons(), data)
	r.renderEntities(glx.EntityTypePlaces, catalog.Places(), data)
	r.renderEntities(glx.EntityTypeRelationships, catalog.Relationships(), data)
	r.renderEntities(glx.EntityTypeRepositories, catalog.Repositories(), data)
	r.renderEntities(glx.EntityTypeSources, catalog.Sources(), data)

	return r.Err
}

func (r *Renderer) renderTopLevelPages(data data.Site) error {
	for name, tmpl := range r.Templates.Pages {
		fname := filepath.Join(r.OutputDir, name+".html")
		err := render(fname, data, tmpl)
		if err != nil {
			return fmt.Errorf("rendering top-level page with template %q: %w", name, err)
		}
	}

	return r.Err
}

type ider interface {
	ID() string
}

func (r *Renderer) renderEntities[S ~[]E, E ider](t glx.EntityType, entities S, siteData data.Site) {
	if r.Err != nil {
		return
	}

	entityTypeData := data.EntityType[S, E]{
		Site:       siteData,
		EntityType: t,
		Entities:   entities,
	}

	entityTypeDir := filepath.Join(r.OutputDir, t.Plural())

	if err := os.Mkdir(entityTypeDir, 0755); err != nil {
		r.Err = fmt.Errorf("rendering %s: %w", t, err)
		return
	}

	templates := r.Templates.EntityTypes[t]
	for name, tmpl := range templates.Pages {
		fname := filepath.Join(entityTypeDir, name+".html")
		err := render(fname, entityTypeData, tmpl)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s with page template %q: %w", t, name, err)
			return
		}
	}

	entityData := data.Entity{
		Site: siteData,
	}
	entityTemplate := templates.Entity
	if entityTemplate == nil {
		r.Err = fmt.Errorf("no entity template for %s", t)
		return
	}

	for _, entity := range entities {
		id := entity.ID()
		entityDir := filepath.Join(entityTypeDir, id)
		if err := os.Mkdir(entityDir, 0755); err != nil {
			r.Err = fmt.Errorf("rendering %s[%s]: %w", t, id, err)
			return
		}

		fname := filepath.Join(entityDir, "index.html")
		entityData.Entity = entity

		err := render(fname, entityData, entityTemplate)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s[%s]: %w", t, id, err)
			return
		}
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
