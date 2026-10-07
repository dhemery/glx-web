package entity

import (
	"slices"
	"strings"
)

func sortValues[S ~[]E, E Stringer](in S) S {
	sorted := slices.Clone(in)
	slices.SortStableFunc(sorted, compareValues)
	return sorted
}

func compareValues[V Stringer](a, b V) int {
	return strings.Compare(a.String(), b.String())
}
