// Package data defines the data items that glx-web applies templates to.
package data

import (
	"github.com/genealogix/glx/go-glx"
)

// EntitySet is the set of all entities of a given type.
type EntitySet struct {
	// The type of entity in the set.
	glx.EntityType
	// A map of the entities, keyed by ID.
	Entities any
}

var (
	theSite     = site{Store: make(map[string]any)}
	currentPage page
)

var TemplateFuncs = map[string]any{
	"page": func() page { return currentPage },
	"site": func() site { return theSite },
}

func ResetPage() {
	currentPage = page{Store: make(map[string]any)}
}

func SetSiteData(data map[string]any) {
	theSite.Data = data
}

func Site() site {
	return theSite
}

type site struct {
	// Data from the site field of congig.yaml
	Data map[string]any
	// Store shared by all pages.
	Store Store
}

type page struct {
	// Store scoped to a page.
	Store Store
}

type Store map[string]any

func (s Store) Get(name string) any {
	return s[name]
}

func (s Store) Set(name string, value any) string {
	s[name] = value
	return ""
}

func (s Store) Delete(name string) string {
	delete(s, name)
	return ""
}
