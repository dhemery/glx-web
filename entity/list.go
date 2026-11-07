package entity

import (
	"iter"
	"slices"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

// Interfaces and impelementations of common list types.
// The methods common to many lists
// are implemented here.
// Methods peculiar to a given element type
// are implemented in the file that implements the element.

// List is the interface
// implemented by each entity list type
// and by [ParticipantList].
// Many implementations offer additional specialized methods
// to sort, filter, group, and map their elements.
// See [EntityList] for additional methods
// impelemented by all entity list types.
type List[S ~[]E, E comparable] interface {
	// Difference returns the unique elements of the list that are not
	// present in other. The result preserves the order of each element's
	// fist occurrence in the list.
	Difference(other S) S

	// Intersection returns the unique elements of the list that are
	// present in other. The result preserves the order of each element's
	// first occurrence in the list.
	Intersection(other S) S

	// Sort returns a copy of the list ordered by comparing the String
	// values of its elements.
	Sort() S

	// Union returns the unique elements of the concatenation of the list
	// and other. The result preserves the order of each element's first
	// occurrence in the concatenation of the list and other.
	Union(other S) S

	// Unique returns the unique elements of the list in the order of their
	// first occurrence.
	Unique() S
}

// EntityList is the interface
// implemented by each entity list type
// (e.g. [PersonList]).
// Many implementations offer additional specialized methods
// to sort, filter, group, and map their elements.
type EntityList[S ~[]E, E comparable] interface {
	List[S, E]
	// All returns an iterator that yields the list elements in order.
	All() iter.Seq[E]

	// EntityType returns the [glx.EntityType] of the elements of the list.
	EntityType() glx.EntityType
}

// Type checks for list types.
var (
	_ List[ParticipantList, *Participant]         = new(ParticipantList)
	_ EntityList[AssertionList, *Assertion]       = new(AssertionList)
	_ EntityList[CitationList, *Citation]         = new(CitationList)
	_ EntityList[EventList, *Event]               = new(EventList)
	_ EntityList[MediaList, *Media]               = new(MediaList)
	_ EntityList[PersonList, *Person]             = new(PersonList)
	_ EntityList[PlaceList, *Place]               = new(PlaceList)
	_ EntityList[RelationshipList, *Relationship] = new(RelationshipList)
	_ EntityList[RepositoryList, *Repository]     = new(RepositoryList)
	_ EntityList[SourceList, *Source]             = new(SourceList)
)

// AssertionList is a list of [*Assertion].
// AssertionList implements [EntityList].
type AssertionList []*Assertion

func (l AssertionList) All() iter.Seq[*Assertion] {
	return slices.Values(l)
}

func (l AssertionList) Difference(other AssertionList) AssertionList {
	return difference(l, other)
}

func (l AssertionList) EntityType() glx.EntityType {
	return glx.EntityTypeAssertions
}

func (l AssertionList) Intersection(other AssertionList) AssertionList {
	return intersection(l, other)
}

func (l AssertionList) Sort() AssertionList {
	return sortedByString(l)
}

func (l AssertionList) Union(other AssertionList) AssertionList {
	return union(l, other)
}

func (l AssertionList) Unique() AssertionList {
	return unique(l)
}

// CitationList is a list of [*Citation].
// CitationList implements [EntityList].
type CitationList []*Citation

func (l CitationList) All() iter.Seq[*Citation] {
	return slices.Values(l)
}

func (l CitationList) EntityType() glx.EntityType {
	return glx.EntityTypeCitations
}

// Sort returns the citations sorted by String value.
func (l CitationList) Sort() CitationList {
	return sortedByString(l)
}

func (l CitationList) Difference(other CitationList) CitationList {
	return difference(l, other)
}

func (l CitationList) Intersection(other CitationList) CitationList {
	return intersection(l, other)
}

func (l CitationList) Union(other CitationList) CitationList {
	return union(l, other)
}

func (l CitationList) Unique() CitationList {
	return unique(l)
}

// EventList is a list of [*Event].
// EventList impelements [EntityList].
type EventList []*Event

func (l EventList) All() iter.Seq[*Event] {
	return slices.Values(l)
}

func (l EventList) EntityType() glx.EntityType {
	return glx.EntityTypeEvents
}

func (l EventList) Sort() EventList {
	return sortedByString(l)
}

func (l EventList) Difference(other EventList) EventList {
	return difference(l, other)
}

func (l EventList) Intersection(other EventList) EventList {
	return intersection(l, other)
}

func (l EventList) Union(other EventList) EventList {
	return union(l, other)
}

func (l EventList) Unique() EventList {
	return unique(l)
}

// MediaList is a list of [*Media].
// MediaList implements [EntityList].
type MediaList []*Media

func (l MediaList) All() iter.Seq[*Media] {
	return slices.Values(l)
}

func (l MediaList) EntityType() glx.EntityType {
	return glx.EntityTypeMedia
}

func (l MediaList) Sort() MediaList {
	return sortedByString(l)
}

func (l MediaList) Difference(other MediaList) MediaList {
	return difference(l, other)
}

func (l MediaList) Intersection(other MediaList) MediaList {
	return intersection(l, other)
}

func (l MediaList) Union(other MediaList) MediaList {
	return union(l, other)
}

func (l MediaList) Unique() MediaList {
	return unique(l)
}

func (l ParticipantList) Difference(other ParticipantList) ParticipantList {
	return difference(l, other)
}

func (l ParticipantList) Intersection(other ParticipantList) ParticipantList {
	return intersection(l, other)
}

func (l ParticipantList) Union(other ParticipantList) ParticipantList {
	return union(l, other)
}

func (l ParticipantList) Unique() ParticipantList {
	return unique(l)
}

// PersonList is a list of [*Person].
// PersonList implements [EntityList].
type PersonList []*Person

func (l PersonList) All() iter.Seq[*Person] {
	return slices.Values(l)
}

func (l PersonList) EntityType() glx.EntityType {
	return glx.EntityTypePersons
}

func (l PersonList) Sort() PersonList {
	return sortedByString(l)
}

func (l PersonList) Difference(other PersonList) PersonList {
	return difference(l, other)
}

func (l PersonList) Intersection(other PersonList) PersonList {
	return intersection(l, other)
}

func (l PersonList) Union(other PersonList) PersonList {
	return union(l, other)
}

func (l PersonList) Unique() PersonList {
	return unique(l)
}

// PlaceList is a list of [*Place].
// PlaceList implements [EntityList].
type PlaceList []*Place

func (l PlaceList) All() iter.Seq[*Place] {
	return slices.Values(l)
}

func (l PlaceList) EntityType() glx.EntityType {
	return glx.EntityTypePlaces
}

func (l PlaceList) Sort() PlaceList {
	return sortedByString(l)
}

func (l PlaceList) Difference(other PlaceList) PlaceList {
	return difference(l, other)
}

func (l PlaceList) Intersection(other PlaceList) PlaceList {
	return intersection(l, other)
}

func (l PlaceList) Union(other PlaceList) PlaceList {
	return union(l, other)
}

func (l PlaceList) Unique() PlaceList {
	return unique(l)
}

// RelationshipList is a list of [*Relationship].
// RelationshipList implements [EntityList].
type RelationshipList []*Relationship

func (l RelationshipList) All() iter.Seq[*Relationship] {
	return slices.Values(l)
}

func (l RelationshipList) EntityType() glx.EntityType {
	return glx.EntityTypeRelationships
}

func (l RelationshipList) Sort() RelationshipList {
	return sortedByString(l)
}

func (l RelationshipList) Difference(other RelationshipList) RelationshipList {
	return difference(l, other)
}

func (l RelationshipList) Intersection(other RelationshipList) RelationshipList {
	return intersection(l, other)
}

func (l RelationshipList) Union(other RelationshipList) RelationshipList {
	return union(l, other)
}

func (l RelationshipList) Unique() RelationshipList {
	return unique(l)
}

// RepositoryList is a list of [*Repository].
// RepositoryList implements [EntityList].
type RepositoryList []*Repository

func (l RepositoryList) All() iter.Seq[*Repository] {
	return slices.Values(l)
}

func (l RepositoryList) EntityType() glx.EntityType {
	return glx.EntityTypeRepositories
}

func (l RepositoryList) Sort() RepositoryList {
	return sortedByString(l)
}
func (l RepositoryList) Difference(other RepositoryList) RepositoryList {
	return difference(l, other)
}

func (l RepositoryList) Intersection(other RepositoryList) RepositoryList {
	return intersection(l, other)
}

func (l RepositoryList) Union(other RepositoryList) RepositoryList {
	return union(l, other)
}

func (l RepositoryList) Unique() RepositoryList {
	return unique(l)
}

// SourceList is a list of [*Source].
// SourceList implements [EntityList].
type SourceList []*Source

func (l SourceList) All() iter.Seq[*Source] {
	return slices.Values(l)
}

func (l SourceList) Difference(other SourceList) SourceList {
	return difference(l, other)
}

func (l SourceList) EntityType() glx.EntityType {
	return glx.EntityTypeSources
}

func (l SourceList) Intersection(other SourceList) SourceList {
	return intersection(l, other)
}

func (l SourceList) Sort() SourceList {
	return sortedByString(l)
}

func (l SourceList) Union(other SourceList) SourceList {
	return union(l, other)
}

func (l SourceList) Unique() SourceList {
	return unique(l)
}

func difference[S ~[]E, E comparable](x, y S) S {
	return slices.DeleteFunc(unique(x), func(element E) bool {
		return slices.Contains(y, element)
	})
}

func intersection[S ~[]E, E comparable](x, y S) S {
	var out S
	for _, element := range unique(x) {
		if slices.Contains(y, element) {
			out = append(out, element)
		}
	}
	return out
}

func union[S ~[]E, E comparable](x, y S) S {
	return unique(slices.Concat(x, y))
}

func unique[S ~[]E, E comparable](x S) S {
	var out S
	for _, element := range x {
		if !slices.Contains(out, element) {
			out = append(out, element)
		}
	}
	return out
}

func sortedByString[S ~[]E, E Stringer](s S) S {
	return slices.SortedStableFunc(slices.Values(s), compareString)
}

func compareString[T Stringer](a, b T) int {
	return strings.Compare(a.String(), b.String())
}
