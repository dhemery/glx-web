package entity

import "github.com/genealogix/glx/go-glx"

// TODO: These struct dosn't fit in this package.

type Index struct {
	Title    string
	Entities any
}

// ArchiveData is the data sent to each template.
type ArchiveData struct {
	// The underlying GLX archive.
	GLX *glx.GLXFile
	// The catalog of entities being rendered.
	Catalog *Catalog
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

// TODO: Add UserData gathered from some YAML file.
