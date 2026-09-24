// Package site renders a website from a compiled GLX archive.
package site

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/archive"
	"github.com/genealogix/glx/go-glx"
)

type renderer struct {
	siteDir   string
	templates map[glx.EntityType]*template.Template
	err       error
}

func Render(a *archive.Archive, templates map[glx.EntityType]*template.Template) error {
	siteDir := "public"
	if err := os.Mkdir(siteDir, 0755); err != nil {
		return err
	}

	r := renderer{siteDir: siteDir, templates: templates}

	r.renderEntities(glx.EntityTypeEvents, a.Events)
	r.renderEntities(glx.EntityTypePersons, a.Persons)
	r.renderEntities(glx.EntityTypePlaces, a.Places)

	return r.err
}

func (r *renderer) renderEntities[T any](t glx.EntityType, entities map[string]*T) {
	if r.err != nil {
		return
	}
	typeDir := filepath.Join(r.siteDir, t.Plural())
	if err := os.Mkdir(typeDir, 0755); err != nil {
		r.err = err
		return
	}

	for id, e := range entities {
		entityDir := filepath.Join(typeDir, id)
		if err := os.Mkdir(entityDir, 0755); err != nil {
			r.err = err
			return
		}

		if err := render(entityDir, e, r.templates[t]); err != nil {
			r.err = err
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
