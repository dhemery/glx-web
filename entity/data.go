package entity

import "github.com/genealogix/glx/go-glx"

// Data is the package of data sent to a template to render.
type Data struct {
	// The content item being rendered. For an entity template, this is the entity being rendered.
	// For an index template, it is a map of all entities of the relevant type, keyed by ID.
	Content any

	// The archive being rendered.
	Archive *Catalog

	// The underlying [glx.GLXFile] of the archive.
	GLX *glx.GLXFile
}

// TODO: Add UserData gathered from some YAML file.
