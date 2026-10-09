package entity

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

type compiler interface {
	// Compile initializes the receiver's references to entities and
	// vocabularies in the context and notifies any referenced entities
	// that the receiver refers to them.
	compile(*context)
}

func compile[C compiler](ctx *context, compilers map[string]C) {
	for _, r := range compilers {
		r.compile(ctx)
	}
}

type context struct {
	GLX     *glx.GLXFile // Underlying GLX File.
	Catalog *Catalog     // Catalog being compiled.
	Entity  entity       // Entity being compiled.
	Path    []string     // Path to the current element within entity.
	Output  io.Writer    // Writer for warnings.
}

func (c *context) ForEntity(e entity) *context {
	out := *c
	out.Entity = e
	out.Path = []string{}
	return &out
}

func (c *context) Sub(p string) *context {
	out := *c
	out.Path = append(slices.Clone(c.Path), p)
	return &out
}

func (c *context) String() string {
	parts := make([]string, 0, len(c.Path)+1)

	parts = append(parts, c.Entity.Plural()+"["+c.Entity.id+"]")
	parts = append(parts, c.Path...)

	return strings.Join(parts, ".")
}

func (c *context) Warnf(format string, a ...any) {
	warning := fmt.Sprintf(format, a...)
	fmt.Fprintf(c.Output, "%s: %s\n", c, warning)
}
