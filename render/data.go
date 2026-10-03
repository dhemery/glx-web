package render

import (
	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

// ArchiveData is data sent to each template to describe the archive being
// published.
type ArchiveData struct {
	// The underlying GLX archive.
	GLX *glx.GLXFile
	// The catalog of entities being rendered.
	Catalog *entity.Catalog

	// TODO(dale): Add user-specified data gathered from some YAML file.
}

// EntityListData is the data sent to each index template for an entity type.
type EntityListData[T any] struct {
	ArchiveData
	// The type of entity being rendered.
	EntityType glx.EntityType
	// The map entities being rendered, keyed by ID.
	Entities map[string]*T
}

// EntityData is the data sent to each entity template.
type EntityData struct {
	ArchiveData
	// The entity being rendered.
	Entity any
}
