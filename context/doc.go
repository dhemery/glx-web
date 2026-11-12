/*
Package context defines the data and functions available to all templates.

The glx-web command makes many functions
available to all templates.
Most are provided by the [Sprout] library.
Templates can access all of the functions
in Sprout's [All] registry group.

# Context Functions

See [Context] for the methods of the returned Context.

	context
	    Returns the Context, which includes the custom data supplied by the
	    configuration file as well as global and local data stores that
	    templates can use to share data.

# Partial Template Functions

A partial template (or simply "partial")
is a template executed by eval
to compute a result
instead of writing to an output file.
Partial templates are specified in the config.yaml file.

	eval <name> <value>
	    Returns the result of applying the named partial template to the
	    value. By default eval returns the text written by the partial. If
	    the partial calls return, eval instead discards the partial's
	    output and returns the result set by the last call to return.

	return <value>
	    Sets the result for the partial at the top of the stack (see note
	    below) and returns the value.

Context maintains a stack of partials.
Executing a partial via the eval function
pushes the partial onto the stack.
When the the completes,
the partial is popped off the stack.
A partial can set its result
by calling the return function.

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
