// Package site renders a website from a compiled GLX archive.
package site

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/dhemery/glx-web/archive"
	"github.com/genealogix/glx/go-glx"
)

func Render(a *archive.Archive, templates map[string]*template.Template) error {
	siteDir := "public"
	if err := os.Mkdir(siteDir, 0755); err != nil {
		return err
	}

	placesDir := filepath.Join(siteDir, glx.EntityTypePlaces.Plural())
	if err := os.Mkdir(placesDir, 0755); err != nil {
		return err
	}

	placeTemplate := templates[glx.EntityTypePlaces.Plural()]

	for id, e := range a.Places {
		placeDir := filepath.Join(placesDir, id)
		if err := os.Mkdir(placeDir, 0755); err != nil {
			return err
		}

		if err := render(placeDir, e, placeTemplate); err != nil {
			return err
		}
	}
	return nil
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
