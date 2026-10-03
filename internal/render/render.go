// Package render renders a website from a compiled GLX archive.
package render

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

type SiteTemplates struct {
	// Templates for each entity type.
	Entities map[glx.EntityType]EntityTemplates
	// Templates for top-level files.
	Indexes map[string]*template.Template
}

type EntityTemplates struct {
	// The template for entities of the ssociated type.
	Entity *template.Template
	// Templates for indexes for each entity type.
	Indexes map[string]*template.Template
}

type Renderer struct {
	OutputDir string
	Clean     bool
	StaticDir string
	Templates SiteTemplates
	Err       error
}

func (r *Renderer) Render(catalog *entity.Catalog, glxfile *glx.GLXFile) error {
	if err := os.Mkdir(r.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	data := entity.SitePageData{
		GLX:     glxfile,
		Catalog: catalog,
	}

	err := r.renderTopLevelIndexes(data)
	if err != nil {
		return err
	}

	r.renderEntities(glx.EntityTypeAssertions, catalog.Assertions, data)
	r.renderEntities(glx.EntityTypeCitations, catalog.Citations, data)
	r.renderEntities(glx.EntityTypeEvents, catalog.Events, data)
	r.renderEntities(glx.EntityTypeMedia, catalog.Media, data)
	r.renderEntities(glx.EntityTypePersons, catalog.Persons, data)
	r.renderEntities(glx.EntityTypePlaces, catalog.Places, data)
	r.renderEntities(glx.EntityTypeRelationships, catalog.Relationships, data)
	r.renderEntities(glx.EntityTypeRepositories, catalog.Repositories, data)
	r.renderEntities(glx.EntityTypeSources, catalog.Sources, data)

	return r.Err
}

func (r *Renderer) renderTopLevelIndexes(data entity.SitePageData) error {
	for name, tmpl := range r.Templates.Indexes {
		fname := filepath.Join(r.OutputDir, name+".html")
		err := render(fname, data, tmpl)
		if err != nil {
			return fmt.Errorf("rendering top-level index with template %q: %w", name, err)
		}
	}

	return r.Err
}

func (r *Renderer) renderEntities[T any](t glx.EntityType, entities map[string]*T, archiveData entity.SitePageData) {
	if r.Err != nil {
		return
	}

	entityTypeData := entity.EntityTypePageData[T]{
		SitePageData: archiveData,
		EntityType:   t,
		Entities:     entities,
	}

	entityTypeDir := filepath.Join(r.OutputDir, t.Plural())

	if err := os.Mkdir(entityTypeDir, 0755); err != nil {
		r.Err = fmt.Errorf("rendering %s: %w", t, err)
		return
	}

	templates := r.Templates.Entities[t]
	for name, tmpl := range templates.Indexes {
		fname := filepath.Join(entityTypeDir, name+".html")
		err := render(fname, entityTypeData, tmpl)
		if err != nil {
			r.Err = fmt.Errorf("rendering %s with index template %q: %w", t, name, err)
			return
		}
	}

	entityData := entity.EntityPageData{
		SitePageData: archiveData,
	}
	entityTemplate := templates.Entity
	if entityTemplate == nil {
		r.Err = fmt.Errorf("no entity template for %s", t)
		return
	}

	for id, entity := range entities {
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
