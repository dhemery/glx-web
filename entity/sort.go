package entity

import (
	"slices"
	"strings"
)

func sortedByString[S ~[]E, E Stringer](s S) S {
	return slices.SortedStableFunc(slices.Values(s), compareString)
}

func compareString[T Stringer](a, b T) int {
	return strings.Compare(a.String(), b.String())
}
