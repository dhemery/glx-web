/*
The glx-web command renders a website
from a GLX archive.
It applies user-specified templates to render files:
  - for each entity,
  - for the list of all entities of each type,
  - and for the website as a whole.

Package [entity]
supplements a GLX archive and its entities
with convenient features that templates can use to render web pages.

Package [tmpl] package
defines how templates are configured
and provides additional data and functions
to templates.

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
