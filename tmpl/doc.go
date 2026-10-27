/*
Package tmpl configures the templates that glx-web uses to render files.
It also defines the data supplied to each template,
and the functions that a template can use
to access and manage data not passed directly to the template.

# Templates

The glx-web command
applies user-specified templates to render files.
It does not supply any built-in templates.

To specify the templates,
write a config.yaml file in the template directory.
The file is deserialized into a [Config].

The configuration defines four kinds of templates:

  - Base templates,
    declared in [SiteSpec.Bases].
    A base template is never invoked directly,
    but can be extended by other templates,
    including by other base templates.

  - Site templates,
    declared in [SiteSpec.Site].
    Each site template is applied to a [Site]
    to render a file about the site as a whole.
    You can render the site information to multiple files
    by specifying multiple site templates.
    The map key for the template
    is used as the path for the rendered file
    relative to the output directory.

  - Type templates,
    declared in [TypeSpec.Type].
    Each type template is applied to an [EntitySet]
    to render a file
    about the set of entities of a given type
    or about the entity type as a whole.
    You can render an entity set to multiple files
    by specifying multiple type templates.
    The map key for the template
    is used as the path for the rendered file
    relative to the output directory.

  - Entity templates,
    declared in [TypeSpec.Entity].
    Each entity template is applied to an entity
    to render a file about that entity.
    You can render an entity to multiple files
    by specifying multiple Entity templates.
    Each occurrence of "<id>" in the map key for each template
    is replaced with the entity's ID,
    and the result is used as the path for the rendered file
    relative to the output directory.

# Custom Data

The configuration also has a [Config.Data] field
that you can use to supply custom data about the site.
This data is copied to [Site.Data],
making it available to all templates.

# Functions

The glx-web command defines several functions
that all templates can access:

  - The store function returns a key/value [Store]
    that persists for the life of the currently executing template.
    You can use this Store
    to compile data for use later in the template.

  - The site function returns the [Site],
    which includes the custom data supplied by the configuration file
    and a global key/value [Store]
    that persists across the execution of all templates.
    You can use this Store
    to compile data to share across all templates.
*/
package tmpl
