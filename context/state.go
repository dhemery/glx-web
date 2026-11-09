package context

import (
	"text/template"
)

// Site is the data passed to each Site template. It is also available to all
// templates via the site function.
type Site struct {
	// The data parsed from the configuration file's data field.
	Data any
	// A key/value store scoped to the site. The store is maintained
	// throughout the rendering of the site. You can use this store to
	// compile data to share across all templates.
	Store Store
}

// Store is a key-value store.
type Store map[string]any

// Set associates key with value.
// It returns an empty string so that it produces no output
// when called without assigning its return value.
func (s Store) Set(key string, value any) string {
	s[key] = value
	return ""
}

// Get returns the value associated with key.
func (s Store) Get(key string) any {
	return s[key]
}

// Delete deletes key and its associated value from s.
// It returns an empty string so that it produces no output
// when called without assigning its return value.
func (s Store) Delete(key string) string {
	delete(s, key)
	return ""
}

type state struct {
	site  Site
	store Store
}

func (s *state) Site() Site {
	return s.site
}

func (s *state) Store() Store {
	return s.store
}

func (s *state) PrepareExecute() {
	s.store = make(Store)
}

// Create creates the state and functions used by the rendering engine.
func Create(data any) (*state, template.FuncMap) {
	site := Site{Data: data, Store: make(Store)}
	state := &state{site: site}
	return state, newFuncMap(state)
}
