package entity

import (
	"fmt"
	"io"

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
	Path    string       // Path to the current element within entity.
	Output  io.Writer    // Writer for warnings.
}

func (c *context) ForEntity(e entity) *context {
	out := *c
	out.Entity = e
	out.Path = ""
	return &out
}

func (c *context) Sub(p string) *context {
	out := *c
	out.Path = out.Path + "." + p
	return &out
}

func (c *context) Subi(i int) *context {
	out := *c
	out.Path = fmt.Sprintf("%s[%d]", c.Path, i)
	return &out
}

func (c *context) String() string {
	return fmt.Sprintf("%s[%s]%s", c.Entity.Plural(), c.Entity.id, c.Path)
}

func (c *context) Warn(warning any) {
	fmt.Fprintf(c.Output, "%s: %s\n", c, warning)
}

func (c *context) Warnf(format string, a ...any) {
	c.Warn(fmt.Sprintf(format, a...))
}

type errParse struct {
	want string // The expected type with optional explanation.
	got  any    // The raw value.
}

func (e errParse) Error() string {
	return fmt.Sprintf("expected %s, got %T", e.want, e.got)
}

type errUnknown struct {
	kind string // The kind of unknown thing.
	name string // The unknown name.
}

func (e errUnknown) Error() string {
	return fmt.Sprintf("unknown %s %q", e.kind, e.name)
}
