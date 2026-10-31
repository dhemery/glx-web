package entity

import (
	"github.com/genealogix/glx/go-glx"
)

// Entity is the type satisfied by every entity in a [Catalog]. Templates can
// check whether a value is an entity by calling the "is_entity" template
// function.
type Entity interface {
	EntityType() glx.EntityType
	ID() string
}

// IsEntity reports whether v satisfies [Entity]. Templates can call this
// function via the "is_entity" template function.
func IsEntity(v any) bool {
	_, ok := v.(Entity)
	return ok
}
