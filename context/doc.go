/*
Package context defines the data and functions available to all templates.

The glx-web command makes many functions
available to all templates.
Most are provided by the [Sprout] library.
Templates can access all of the functions
in Sprout's [All] registry group.

# Data Functions

These functions supply data from the context to the template.

	site
	    Returns the Site, which includes the custom data supplied by the
	    configuration file and a global key/value Store that persists
	    across the execution of all templates. You can use this site Store
	    to compile data to share across all templates.

	store
	    Returns a Store that persists for the life of the currently
	    executing template. You can use this Store to compile data for use
	    anywhere in the template, including in associated templates.

# Type Test Functions

These functions report about the types of their arguments.

	is_entity
	    Reports whether its argument satisfies entity.Entity.

	is_property
	    Reports whether its argument's type is *entity.Property.

	is_vocabulary_value
	    Reports whether its argument satisfies entity.VocabularyValue.

# Set-Like Functions

These functions are analogous to set operations,
with two differences:

  - If the arguments contain duplicates
    the result may contain duplicates.
    Use the Sprout "uniq" function
    to remove duplicates from the arguments or the results.
  - The result preserves the order of elements in the inputs.

See the descriptions for details.

	difference
	    Returns a new slice of the elements of arg1 that are not present in
	    arg2. Both arguments must be a slices. Any element duplicated in
	    arg1 is duplicated in the result if it is not present in arg2. The
	    result preserves the order of elements in arg1.


	intersection
	    Returns a new slice of the elements that are present in both arg1
	    and arg2. Both arguments must be a slices. Any element duplicated
	    in arg1 is duplicated in the result if it is present in arg2. The
	    result preserves the order of elements in arg1.

	union
	    Returns a new slice of the elements that appear in arg1, arg2, or
	    both. Both arguments must be a slices. Each element is present only
	    once in the result. The result preserves the order of each
	    element's first occurrence in the concatenation of arg1 and arg2.

[All]: https://docs.atom.codes/sprout/groups/all
[Sprout]: https://docs.atom.codes/sprout
*/
package context
