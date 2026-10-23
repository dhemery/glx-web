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

	err := r.renderTopLevelPages()
	if err != nil {
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

func (r *Renderer) renderTopLevelPages() error {
	dot := data.Site()
	for name, tmpl := range r.Templates.Pages {
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

func (r *Renderer) renderEntities[S ~[]E, E ider](t glx.EntityType, entities S) {
	if r.Err != nil {
		return
	}

	entityTypeDir := filepath.Join(r.OutputDir, t.Plural())

	if err := os.Mkdir(entityTypeDir, 0755); err != nil {
		r.Err = fmt.Errorf("rendering %s: %w", t, err)
		return
	}

	templates := r.Templates.EntityTypes[t]
	entitySet := data.EntitySet{
		EntityType: t,
		Entities:   entities,
	}

	for name, tmpl := range templates.Pages {
		fname := filepath.Join(entityTypeDir, name+".html")
		err := render(fname, entitySet, tmpl)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s with page template %q: %w", t, name, err)
			return
		}
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

		err := render(fname, entity, entityTemplate)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s[%s]: %w", t, id, err)
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

	data.ResetPage()

	return tmpl.Execute(f, dot)
}
