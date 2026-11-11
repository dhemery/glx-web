/*
Package tmpl
defines the templates, data, and functions
that glx-web uses to render files.

The glx-web command
applies user-specified templates to render files.
It does not supply any built-in templates.

To specify the templates and to supply custom data,
write a config.yaml file in the template directory.
The file is deserialized into a [Config].

# Templates

The [Config.Templates] field defines five kinds of templates:

  - Base templates,
    declared in [SiteSpec.Bases].
    A base template is never invoked directly,
    but can be extended by other templates,
    including by other base templates.

  - Site templates,
    declared in [SiteSpec.Site].
    Each site template is applied to a [Site].
    You can render the site information to multiple files
    by specifying multiple site templates.
    The map key for the template
    is used as the path for the rendered file
    relative to the output directory.

  - Type templates,
    declared in [TypeSpec.Type].
    Each type template is applied to a list of entities.
    Each list type has an EntityType method
    that returns the [glx.EntityType]
    of the entities in the list.
    You can render an entity list to multiple files
    by specifying multiple type templates.
    The map key for the template
    is used as the path for the rendered file
    relative to the output directory.

  - Entity templates,
    declared in [TypeSpec.Entity].
    Each entity template is applied to an entity.
    You can render an entity to multiple files
    by specifying multiple Entity templates.
    Each occurrence of "<id>" in the map key for each template
    is replaced with the entity's ID,
    and the result is used as the path for the rendered file
    relative to the output directory.

  - Partial templates.
    A partial template computes one or more values
    instead of rendering to a file.
    To execute a partial template,
    call the eval or format function.

# Custom Data

Use the [Config.Data] field
to supply custom data about the site.
This data is copied to [Site.Data],
making it available to all templates.

# Functions

The glx-web command defines several functions
that all templates can use:

	is_entity
	    Reports whether its argument satisfies entity.Entity.

	is_property
	    Reports whether its argument's type is entity.Property.

	is_vocabulary_value
	    Reports whether its argument's type satisfies entity.VocabularyValue.

	site
	    Returns the Site, which includes the custom data supplied by the
	    configuration file and a global key/value Store that persists
	    across the execution of all templates. You can use this site Store
	    to compile data to share across all templates.

	store
	    Returns a Store that persists for the life of the currently
	    executing template. You can use this Store to compile data for use
	    later in the template.

See [Site] and [Store] for descriptions of the values returned by the
site and store functions.
*/
package tmpl
