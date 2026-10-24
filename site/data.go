package site

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
	site      = Site{Data: make(map[string]any), Store: make(Store)}
	siteStore = make(Store)
	pageStore = make(Store)
)

type Site struct {
	Data  map[string]any
	Store Store
}

var templateFuncs = map[string]any{
	"site":  func() Site { return site },
	"store": func() Store { return pageStore },
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
