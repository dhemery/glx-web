package context

import (
	"fmt"
	"io"
	"strings"
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
	site     Site
	stores   []Store
	partials map[string]*template.Template
}

func (s *state) Site() Site {
	return s.site
}

func (s *state) Store() Store {
	return s.stores[len(s.stores)-1]
}

func (s *state) Push() {
	s.stores = append(s.stores, make(Store))
}

func (s *state) Pop() {
	s.stores = s.stores[:len(s.stores)-1]
}

func (s *state) SetPartials(partials map[string]*template.Template) {
	s.partials = partials
}

func (s *state) eval(name string, data any) (Store, error) {
	result, err := s.executePartial(name, data, io.Discard)
	if err != nil {
		return nil, fmt.Errorf("eval: %w", err)
	}
	return result, nil
}

func (s *state) format(name string, data any) (string, error) {
	var sb strings.Builder
	if _, err := s.executePartial(name, data, &sb); err != nil {
		return "", fmt.Errorf("format: %w", err)
	}
	return sb.String(), nil
}

func (s *state) executePartial(name string, data any, w io.Writer) (Store, error) {
	format, ok := s.partials[name]
	if !ok {
		return nil, fmt.Errorf("no such partial template: %q", name)
	}

	s.Push()
	defer s.Pop()

	if err := format.Execute(w, data); err != nil {
		return nil, err
	}
	return s.stores[len(s.stores)-1], nil
}

// Create creates the state and functions used by the rendering engine.
func Create(data any) (*state, map[string]any) {
	site := Site{Data: data, Store: make(Store)}
	state := &state{site: site}
	return state, newFuncMap(state)
}
