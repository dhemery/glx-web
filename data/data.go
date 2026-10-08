// Package data defines the data items that glx-web applies templates to.
package data

import (
	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

// Site describes the site being rendered.
type Site struct {
	// The underlying GLX archive being rendered.
	GLX *glx.GLXFile
	// The catalog of entities being rendered.
	Catalog *entity.Catalog

	// TODO(feature): Add data parsed from a user-supplied site configuration
	// YAML file.
}

// EntityType describes the site and the set of entities of type T.
type EntityType[S ~[]T, T any] struct {
	Site
	// The type of entity being rendered.
	EntityType glx.EntityType
	// The entities being rendered, keyed by ID.
	Entities S
}

// Entity describes the site and an entity.
type Entity struct {
	Site
	// The entity being rendered.
	Entity any
}
