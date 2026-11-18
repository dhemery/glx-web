package tmpl

import (
	"fmt"

	"github.com/dhemery/glx-web/context"
	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
	"go.yaml.in/yaml/v3"
)

// Config specifies the templates and data used by glx-web.
type Config struct {
	// User-defined data available to each template. Each Site template is
	// applied directly to this value. Other templates can access the data
	// via the config function.
	Data any `json:"data"`
	// Templates to render the site files.
	Templates SiteSpec `json:"templates"`
	// Each Partial template can be executed by the eval function to
	// compute a result. The map key is used as the name of the partial.
	// The map value is the list of template file paths to parse relative
	// to the template directory.
	Partials map[string]FileList `json:"partials"`
}

// SiteSpec specifies the templates for the site and for each entity type.
// Each template can use the functions defined by package [context].
type SiteSpec struct {
	// Each Base template defines a base that other templates can build on.
	// The map key is used as the name of the base. Base templates can
	// build on other base templates.
	Bases map[string]Spec `json:"bases"`
	// Each Site template is applied to the [Config.Data]. The map key is
	// used as the path for the rendered file relative to the output
	// directory.
	Site map[string]Spec `json:"site"`
	// Templates for each entity type.
	Types map[glx.EntityType]TypeSpec `json:"types"`
}

// TypeSpec specifies templates for an entity type and for entities of that
// type.
type TypeSpec struct {
	// Each Type template is applied to the list of all entities of the
	// type. Each entity list has a distinct type with methods relevant to
	// the type of entity. For example, the list for
	// [glx.EntityTypePersons] is [entity.PersonList], which has a
	// GroupBySurname method. Each list type has an EntityType method that
	// returns the [glx.EntityType] of its elements. Some lists have
	// methods to sort, filter, group, or map the elements in ways useful
	// to templates. The map key is used as the path for the rendered file
	// relative to the output directory.
	Type map[string]Spec `json:"type"`
	// Each Entity template is applied to each entity of the type. Each
	// occurrence of "<id>" in the map key is replaced by the entity's ID,
	// and the result is used as the path for the rendered file relative to
	// the output directory.
	Entity map[string]Spec `json:"entity"`
}

// Spec specifies a template. The template is compiled by cloning the base and
// parsing the files into it. If no base is named, there must be at least one
// file, and the first file must define the body of the template.
//
// If the Spec needs one file and no base, you can represent it in YAML as a
// scalar single string instead of a mapping.
type Spec struct {
	// The name of the base template to build on.
	Base string `json:"base"`
	// The paths to the template files to parse relative to the template
	// directory. If Files needs one file, you can represent it in YAML as
	// a scalar string instead of a sequence of strings.
	Files FileList `json:"files"`
}

// FileList is the list of files to parse into a single template.
type FileList []string

const yamlTagString = "!!str"

// UnmarshalYAML unmarshals the Spec from either a mapping or a single string.
func (s *Spec) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		if n.Tag != yamlTagString {
			return fmt.Errorf("Spec: want mapping or string, got %s: %q", n.Tag, n.Value)
		}
		s.Files = []string{n.Value}
		return nil
	case yaml.MappingNode:
		var local struct {
			Base  string
			Files FileList
		}
		if err := n.Decode(&local); err != nil {
			return err
		}
		s.Base, s.Files = local.Base, local.Files
		return nil
	}

	return fmt.Errorf("Spec: want string or mapping, got %s: %q", n.Tag, n.Value)
}

// UnmarshalYAML unmarshalss the FileList from either a sequence of strings or
// a single string.
func (fl *FileList) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		if n.Tag != yamlTagString {
			return fmt.Errorf("FileList: want string, got %s: %q", n.Tag, n.Value)
		}
		*fl = []string{n.Value}
		return nil
	case yaml.SequenceNode:
		var files []string
		for _, contentNode := range n.Content {
			if contentNode.Tag != yamlTagString {
				return fmt.Errorf("FileList: want string, got %s: %q", contentNode.Tag, contentNode.Value)
			}
			files = append(files, contentNode.Value)
		}
		*fl = files
		return nil
	}
	return fmt.Errorf("FileList: want string or array, got %s: %q", n.Tag, n.Value)
}

// Import for doc reference.
var _ entity.PersonList
var _ context.Context
