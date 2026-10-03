package entity

import (
	"github.com/genealogix/glx/go-glx"
)

// An SitePageData is sent to each template.
type SitePageData struct {
	// The underlying GLX archive being rendered.
	GLX *glx.GLXFile
	// The catalog of entities being rendered.
	Catalog *Catalog

	// TODO(dale): Add user-specified data parsed from some YAML file.
}

// An EntityTypePageData is sent to each entity type index template for entities of
// type T.
type EntityTypePageData[T any] struct {
	SitePageData
	// The type of entity being rendered.
	EntityType glx.EntityType
	// The map of entities being rendered, keyed by ID.
	Entities map[string]*T
}

// An EntityPageData is sent to each entity template.
type EntityPageData struct {
	SitePageData
	// The entity being rendered.
	Entity any
}
