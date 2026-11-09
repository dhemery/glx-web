package context

import (
	"maps"
	"text/template"

	"github.com/dhemery/glx-web/entity"
	"github.com/go-sprout/sprout"
	"github.com/go-sprout/sprout/group/hermetic"
)

func newFuncMap(state *state) template.FuncMap {
	funcs := template.FuncMap{
		"is_entity":           IsEntity,
		"is_property":         IsProperty,
		"is_vocabulary_value": IsVocabularyValue,
		"site":                state.Site,
		"store":               state.Store,
	}

	sh := sprout.New()
	sh.AddGroups(hermetic.RegistryGroup())
	maps.Insert(funcs, maps.All(sh.Build()))

	return funcs

}

// IsEntity reports whether v satisfies [entity.Entity]. Templates can call this
// function via the "is_entity" template function.
func IsEntity(v any) bool {
	_, ok := v.(entity.Entity)
	return ok
}

// IsProperty reports whether v's type is [*entity.Property]. Templates can
// call this function via the "is_property" template function.
func IsProperty(v any) bool {
	_, ok := v.(*entity.Property)
	return ok
}

// IsVocabularyValue reports whether v satisfies [entity.VocabularyValue]. Templates can call
// this function via the "is_vocabulary_value" template function.
func IsVocabularyValue(v any) bool {
	_, ok := v.(entity.VocabularyValue)
	return ok
}
