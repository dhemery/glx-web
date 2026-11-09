/*
Package context defines the data and functions available to all templates.

The glx-web command makes many functions
available to all templates.
Most are provided by the [Sprout] library.
Templates can access all of the functions
in Sprout's [hermetic] function group.

In addition, templates can access several functions specific to glx-web:

	is_entity
	    Reports whether its argument satisfies entity.Entity.

	is_property
	    Reports whether its argument's type is entity.Property.

	is_vocabulary_value
	    Reports whether its argument satisfies entity.VocabularyValue.

	site
	    Returns the Site, which includes the custom data supplied by the
	    configuration file and a global key/value Store that persists
	    across the execution of all templates. You can use this site Store
	    to compile data to share across all templates.

	store
	    Returns a Store that persists for the life of the currently
	    executing template. You can use this Store to compile data for use
	    anywhere in the template, including in associated templates.

[hermetic]: https://docs.atom.codes/sprout/groups/hermetic
[Sprout]: https://docs.atom.codes/sprout
*/
package context
