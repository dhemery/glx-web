// Package context defines the environment in which templates execute. It
// includes the types passed to templates and the functions that glx-web makes
// available to templates.
package context

import (
	"github.com/genealogix/glx/go-glx"
)

// Site is the type passed to each Site template.
type Site struct {
	// The data parsed from the configuration file's data field.
	Data map[string]any
	// A store that persists through the rendering of all pages.
	Store Store
}

// EntitySet is the type passed to each Set template.
type EntitySet struct {
	// The type of entity in the set.
	glx.EntityType
	// A map of the entities, keyed by ID.
	Entities any
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
