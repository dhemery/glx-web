package context

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
)

// Context holds information about the current state of rendering.
// Templates can access the Context via the context function.
type Context interface {
	// Data returns the usr-supplied custom data from the config file.
	Data() any

	// Globals returns a map that templates can use to share information
	// globally. Globals is initially empty. Its state persists from one
	// template to the next.
	Globals() map[string]any

	// Locals returns a map that templates can use to share information
	// while rendering a single file. Each Site, Type, and Entity template
	// starts with a new, empty Locals.
	Locals() map[string]any

	// Path returns the path to the file being rendered, relative to the
	// output directory.
	Path() string
}

var _ Context = &ctx{}

// New creates a new execution context with the given data.
// The result implements [Context].
func New(data any) (*ctx, template.FuncMap) {
	c := &ctx{
		data:    data,
		globals: make(map[string]any),
	}
	return c, newFuncMap(c)
}

type ctx struct {
	data    any
	globals map[string]any
	locals  map[string]any
	path    string

	partials map[string]*template.Template
	results  resultStack
}

func (c *ctx) Data() any {
	return c.data
}

func (c *ctx) Globals() map[string]any {
	return c.globals
}

func (c *ctx) Locals() map[string]any {
	return c.locals
}

func (c *ctx) Path() string {
	return c.path
}

func (c *ctx) StartPage(path string) {
	c.locals = make(map[string]any)
	c.path = path
}

func (c *ctx) InstallPartials(partials map[string]*template.Template) {
	c.partials = partials
}

func (c *ctx) eval(name string, data any) (any, error) {
	tmpl, ok := c.partials[name]
	if !ok {
		return nil, fmt.Errorf("no such partial template: %q", name)
	}

	c.results.push()
	defer c.results.pop()

	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return nil, err
	}

	result := c.results.top()

	// If the partial set a result, return the result.
	if result != nil {
		return *result, nil
	}

	// The partial did not set a result. Return its output.
	return output.String(), nil
}

var errEmptyResultStack = errors.New("no active partial")

// resultStack is the stack of results set by the currently executing partials.
type resultStack []*any

// push pushes a new unset result onto the stack.
func (s *resultStack) push() {
	*s = append(*s, nil)
}

// set sets the result at the top of the stack.
func (s *resultStack) set(result any) (any, error) {
	top := len(*s) - 1
	if top < 0 {
		return nil, errEmptyResultStack
	}

	(*s)[top] = &result

	return result, nil
}

// top returns the top result from the stack.
func (s *resultStack) top() *any {
	return (*s)[len(*s)-1]
}

// pop removes the top result from the stack.
func (s *resultStack) pop() {
	*s = (*s)[:len(*s)-1]
}
