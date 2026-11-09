package context

import (
	"fmt"
	"maps"
	"reflect"
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
		"difference":          Difference,
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

// Difference returns a new list
// of the elements of x
// that are not present in y.
// Each argument must be a slice or an array.
// Difference is available to templates as "difference".
func Difference(x, y any) ([]any, error) {
	listValue := reflect.ValueOf(x)
	omitValue := reflect.ValueOf(y)

	switch listValue.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		return nil, fmt.Errorf("first argument must be a slice but got %T", x)
	}
	switch omitValue.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		return nil, fmt.Errorf("second argument must be a slice but got %T", x)
	}

	omitLength := omitValue.Len()
	omitSet := make(map[any]struct{}, omitLength)
	for i := range omitLength {
		omitSet[omitValue.Index(i).Interface()] = struct{}{}

	}

	length := listValue.Len()
	result := make([]any, 0, length)

	for i := range length {
		item := listValue.Index(i).Interface()
		if _, found := omitSet[item]; !found {
			result = append(result, item)
		}
	}

	return result, nil
}
