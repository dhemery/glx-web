/*
The glx-web command renders a website
from a GLX archive.
It applies user-specified templates to render pages:
  - for each entity,
  - for the set of entities of each type,
  - and for the website as a whole.

See [entity.Catalog]
and the entity package documentation
for details of how glx-web
supplements a GLX archive and its entities
with features that make them
more convenient to render via templates.

See [tmpl.Config]
and the tmpl package documentation
for details of how to specify templates
and use the data and functions
that glx-web makes available to templates.

The glx-web project
is a work in progress.
Every feature is subject to change without notice.

The glx-web command:

  - Expects each entity's properties to be well-formed,
    even beyond the guarantees enforced by GLX deserialization.

  - Currently does not privatize living persons.

  - Currently does not render [glx.ResearchLog] and [glx.Study] entities,
    but does render every other GLX entity.
*/
package main

import (
	"github.com/dhemery/glx-web/cmd"
	"github.com/dhemery/glx-web/entity"
	"github.com/dhemery/glx-web/tmpl"
	"github.com/genealogix/glx/go-glx"
)

// Imports to support doc comments.
var _ entity.Catalog
var _ tmpl.Config
var _ glx.ResearchLog

func main() {
	cmd.Execute()
}
