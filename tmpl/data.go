package tmpl

import (
	"github.com/genealogix/glx/go-glx"
)

// Site is the data passed to each Site template. It is also available to all
// templates via the site function.
type Site struct {
	// The data parsed from the configuration file's data field.
	Data map[string]any
	// A key/value store scoped to the site. The store is maintained
	// throughout the rendering of the site.
	Store Store
}

// EntitySet is the data passed to each Type template.
type EntitySet struct {
	// The type of entity in the set.
	glx.EntityType
	// A map of the entities, keyed by ID.
	Entities any
}

// Store is a key-value store.
type Store map[string]any

// Set associates key with value.
func (s Store) Set(key string, value any) string {
	s[key] = value
	return ""
}

// Get returns the value associated with key.
func (s Store) Get(key string) any {
	return s[key]
}

// Delete deletes key and its associated value from s.
func (s Store) Delete(key string) string {
	delete(s, key)
	return ""
}
