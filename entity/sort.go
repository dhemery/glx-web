package entity

import (
	"fmt"
	"slices"
	"strings"
)

func sortStringers[S ~[]E, E fmt.Stringer](in S) S {
	sorted := slices.Clone(in)
	slices.SortStableFunc(sorted, compareStringers)
	return sorted
}

func compareStringers[T fmt.Stringer](a, b T) int {
	return strings.Compare(a.String(), b.String())
}
