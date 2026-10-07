package entity

import (
	"path"

	"github.com/genealogix/glx/go-glx"
)

type entity struct {
	glx.EntityType
	id string
}

func (e entity) ID() string {
	return e.id
}

// PagePath returns the path to the directory where glx-web renders the entity,
// relative to the site's output directory.
func (e entity) PagePath() string {
	return path.Join(e.Plural(), e.id)
}
