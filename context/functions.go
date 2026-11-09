package context

import (
	"fmt"
	"maps"
	"reflect"
	"text/template"

	"github.com/dhemery/glx-web/entity"
	"github.com/go-sprout/sprout"
	"github.com/go-sprout/sprout/group/all"
)

func newFuncMap(state *state) template.FuncMap {
	funcs := template.FuncMap{
		// Data functions
		"site":  state.Site,
		"store": state.Store,

		// Type test functions
		"is_entity":           isEntity,
		"is_property":         isProperty,
		"is_vocabulary_value": isVocabularyValue,

		// Set-like list functions
		"difference":   difference,
		"intersection": intersection,
		"union":        union,
	}

	sh := sprout.New()
	sh.AddGroups(all.RegistryGroup())
	maps.Insert(funcs, maps.All(sh.Build()))

	return funcs

}

func isEntity(v any) bool {
	_, ok := v.(entity.Entity)
	return ok
}

func isProperty(v any) bool {
	_, ok := v.(*entity.Property)
	return ok
}

func isVocabularyValue(v any) bool {
	_, ok := v.(entity.VocabularyValue)
	return ok
}

func difference(x, y any) ([]any, error) {
	listValue := reflect.ValueOf(x)
	omitValue := reflect.ValueOf(y)

	switch listValue.Kind() {
	case reflect.Slice:
	default:
		return nil, fmt.Errorf("first argument must be a slice but got %T", x)
	}
	switch omitValue.Kind() {
	case reflect.Slice:
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
		if _, omit := omitSet[item]; !omit {
			result = append(result, item)
		}
	}

	return result, nil
}

func intersection(x, y any) ([]any, error) {
	listValue := reflect.ValueOf(x)
	keepValue := reflect.ValueOf(y)

	switch listValue.Kind() {
	case reflect.Slice:
	default:
		return nil, fmt.Errorf("first argument must be a slice but got %T", x)
	}
	switch keepValue.Kind() {
	case reflect.Slice:
	default:
		return nil, fmt.Errorf("second argument must be a slice but got %T", x)
	}

	keepLength := keepValue.Len()
	keepSet := make(map[any]struct{}, keepLength)
	for i := range keepLength {
		keepSet[keepValue.Index(i).Interface()] = struct{}{}

	}

	length := listValue.Len()
	result := make([]any, 0, length)

	for i := range length {
		item := listValue.Index(i).Interface()
		if _, keep := keepSet[item]; keep {
			result = append(result, item)
		}
	}

	return result, nil
}

func union(x, y any) ([]any, error) {
	leftValue := reflect.ValueOf(x)
	rightValue := reflect.ValueOf(y)

	switch leftValue.Kind() {
	case reflect.Slice:
	default:
		return nil, fmt.Errorf("first argument must be a slice but got %T", x)
	}
	switch rightValue.Kind() {
	case reflect.Slice:
	default:
		return nil, fmt.Errorf("second argument must be a slice but got %T", x)
	}

	leftLength := leftValue.Len()
	rightLength := rightValue.Len()

	addedSet := make(map[any]struct{}, leftLength+rightLength)

	var result []any

	for v := range leftValue.Seq() {
		item := v.Interface()
		if _, added := addedSet[item]; !added {
			addedSet[v] = struct{}{}
			result = append(result, item)
		}
	}

	return result, nil
}
