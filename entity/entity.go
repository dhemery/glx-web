package entity

import (
	"github.com/genealogix/glx/go-glx"
)

type entity struct {
	glx.EntityType `json:"-"`
	id             string
}

func (e entity) ID() string {
	return e.id
}
