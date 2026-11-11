/*
Package context defines the data and functions available to all templates.

The glx-web command makes many functions
available to all templates.
Most are provided by the [Sprout] library.
Templates can access all of the functions
in Sprout's [All] registry group.

# Context Functions

These functions supply data from the context to the template.

	site
	    Returns the Site, which includes the custom data supplied by the
	    configuration file and a global key/value Store that persists
	    across the execution of all templates. You can use this site Store
	    to compile data to share across all templates.

	store
	    Returns a Store that persists for the life of the template
	    rendering the current file. You can use this Store to compile data
	    for use anywhere in the template, including in associated
	    templates.

# Partial Execution Functions

Each of these functions
executes a named partial template to compute a result.
Partial templates are specified in the config.yaml file.

Each call establishes a new context for the partial template,
including a store that persists for the life of the partial.

	eval <name> <value>
	    Applies the named partial template to the value and returns the
	    resulting store. The output written by the template is discarded.

	format <name> <value>
	    Applies the named partial template to the value and returns the
	    output written by the template. The partial's store is discarded.

# Type Test Functions

These functions report about the types of their arguments.

	is_entity <value>
	    Reports whether the value satisfies entity.Entity.

	is_property <value>
	    Reports whether the value's type is *entity.Property.

	is_vocabulary_value <value>
	    Reports whether the value satisfies entity.VocabularyValue.

# Set-Like Functions

These functions are analogous to set operations,
with two differences:

  - If the arguments contain duplicates
    the result may contain duplicates.
    Use the Sprout "uniq" function
    to remove duplicates from the arguments or the results.
  - The result preserves the order of elements in the inputs.

See the descriptions for details.

	difference <list1> <list2>
	    Returns a new slice of the elements of list1 that are not present
	    in list2. Both arguments must be a slices. Any element duplicated
	    in list1 is duplicated in the result if it is not present in list2.
	    The result preserves the order of elements in list1.


	intersection <list1> <list2>
	    Returns a new slice of the elements that are present in both list1
	    and list2. Both arguments must be a slices. Any element duplicated
	    in list1 is duplicated in the result if it is present in list2. The
	    result preserves the order of elements in list1.

	union <list1> <list2>
	    Returns a new slice of the elements that appear in list1, list2, or
	    both. Both arguments must be a slices. Each element is present only
	    once in the result. The result preserves the order of each
	    element's first occurrence in the concatenation of list1 and list2.

[All]: https://docs.atom.codes/sprout/groups/all
[Sprout]: https://docs.atom.codes/sprout
*/
package context
