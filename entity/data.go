package entity

import (
	"github.com/genealogix/glx/go-glx"
)

// SiteData describes the site.
type SiteData struct {
	// The underlying GLX archive being rendered.
	GLX *glx.GLXFile
	// The catalog of entities being rendered.
	Catalog *Catalog

	// TODO(dale): Add user-specified data parsed from some YAML file.
}

// EntityTypeData describes the set of entities of type T and the site.
type EntityTypeData[T any] struct {
	SiteData
	// The type of entity being rendered.
	EntityType glx.EntityType
	// The map of entities being rendered, keyed by ID.
	Entities map[string]*T
}

// EntityData describes an entity and the site.
type EntityData struct {
	SiteData
	// The entity being rendered.
	Entity any
}
