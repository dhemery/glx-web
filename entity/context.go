package entity

import (
	"io"

	"github.com/genealogix/glx/go-glx"
)

type context struct {
	C      *Catalog
	G      *glx.GLXFile
	ErrOut io.Writer
}
