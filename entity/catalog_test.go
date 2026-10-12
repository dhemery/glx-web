package entity_test

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io/fs"
	"maps"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v3"
)

// TestCatalogCompile tests how [entity.NewCatalog] compiles GLX archives. Each
// test case is a [txtar.Archive] read from testdata, contianing a comment, one
// or more GLX files, and a want-catalog.yaml file. The GLX files are looaded
// into a [glx.GLXFile] to pass to NewCatalog. The comment is a regexp pattern
// that the output written by NewCatalog must match. The want-catalog.yaml file
// is a YAML representation of the desired catalog structure. This YAML
// representation is simplified to omit certain details that are not needed to
// varify catalog structure.
func TestCatalogCompile(t *testing.T) {
	for name, tc := range textArchives(t, "testdata/catalog/compile/*") {
		t.Run(filepath.Base(name), func(t *testing.T) {
			fsys, err := txtar.FS(tc)
			if err != nil {
				t.Fatal(err)
			}

			glxfile := glxFile(t, fsys)

			output := &bytes.Buffer{}
			cat := entity.NewCatalog(glxfile, output)

			gotOutput := output.Bytes()

			wantOutput := regexp.MustCompile(string(bytes.TrimSpace(tc.Comment)))
			if !wantOutput.Match(gotOutput) {
				t.Errorf("\n error output: %s\nwant match of: %s", gotOutput, tc.Comment)
			}

			// Marshal and unmarshal the catalog to convert it to a
			// map[string]any that can be compared to the structure
			// represented by want-catalog.yaml.
			gotCatalog := asObject(t, cat, catalogMarshalOptions)
			wantCatalog := unmarshalYAML(t, fsys, "want-catalog.yaml")

			if diff := cmp.Diff(wantCatalog, gotCatalog); diff != "" {
				t.Error(diff)
			}
		})
	}
}

// Convert in to an object by marshaling and unmarshaling. We're using instead
// of YAML because JSOM gives more control over marshaling.
func asObject(t *testing.T, in any, opts json.Options) any {
	t.Helper()
	j, err := json.Marshal(in, opts)
	if err != nil {
		t.Fatal("marshaling:", err)
	}

	var out any
	if err := json.Unmarshal(j, &out); err != nil {
		t.Fatal("unmarshaling:", err)
	}
	return out
}

// Unmarshal the named YAML file.
func unmarshalYAML(t *testing.T, fsys fs.FS, name string) any {
	t.Helper()
	j, err := fs.ReadFile(fsys, name)
	if err != nil {
		t.Fatalf("reading %s: %s", name, err)
	}

	var out any
	if err := yaml.Unmarshal(j, &out); err != nil {
		t.Fatalf("unmarshaling %s: %s", name, err)
	}
	return out
}

func glxFile(t *testing.T, fsys fs.FS) *glx.GLXFile {
	const pattern = "*.glx"
	t.Helper()
	fnames, err := fs.Glob(fsys, pattern)
	if err != nil {
		t.Fatalf("globbing %q: %s", pattern, err)
	}

	if len(fnames) == 0 {
		t.Fatalf("no files matching %q", pattern)
	}

	files := map[string][]byte{}

	// Add the standard vocabularies.
	maps.Copy(files, glx.StandardVocabularies())

	// Add files from the test case.
	for _, fname := range fnames {
		content, err := fs.ReadFile(fsys, fname)
		if err != nil {
			t.Fatal(err)
		}
		files[fname] = content
	}

	serializer := glx.NewSerializer(nil)
	glxfile, dups, err := serializer.DeserializeMultiFileFromMap(files)
	if err != nil {
		t.Fatal("deserializing GLX archive:", err)
	}

	if len(dups) > 0 {
		t.Fatal("duplicates in GLX archive:", dups)
	}

	return glxfile
}

func textArchives(t *testing.T, pattern string) map[string]*txtar.Archive {
	t.Helper()
	fnames, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("globbing %q: %s", pattern, err)
	}

	if len(fnames) == 0 {
		t.Fatalf("no files matching %q", pattern)
	}

	tests := make(map[string]*txtar.Archive, len(fnames))

	for _, fname := range fnames {
		ar, err := txtar.ParseFile(fname)
		if err != nil {
			t.Fatal("reading text archive:", err)
		}

		tests[fname] = ar
	}

	return tests
}

// JSON marshal options to reduce the JSON representation of the catalog to
// only the details we want to compare. The resulting HSON omits each entity's
// back-references and inner GLX entity. It also omits details of GLX
// vocabularies (e.g. Description) that are useful for validating catalog
// structure.
var catalogMarshalOptions = json.JoinOptions(
	json.FormatNilMapAsNull(true),
	json.FormatNilSliceAsNull(true),
	json.OmitZeroStructFields(true),
	json.WithMarshalers(
		json.JoinMarshalers(
			// Entities
			json.MarshalToFunc(marshalPerson),
			json.MarshalToFunc(marshalPlace),

			// Primitive Values
			json.MarshalToFunc(marshalBoolValue),
			json.MarshalToFunc(marshalDate),
			json.MarshalToFunc(marshalIntValue),
			json.MarshalToFunc(marshalStringValue),

			// GLX types
			json.MarshalToFunc(marshalFieldDefinition),
			json.MarshalToFunc(marshalPropertyDefinition),
			json.MarshalToFunc(marshalVocabularyEntry),
		),
	),
)

func marshalBoolValue(out *jsontext.Encoder, in entity.BoolValue) error {
	return out.WriteToken(jsontext.String(fmt.Sprintf("BoolValue{%t}", in)))
}

func marshalDate(out *jsontext.Encoder, in glxdate.Date) error {
	return out.WriteToken(jsontext.String(in.Raw()))
}

func marshalFieldDefinition(out *jsontext.Encoder, in *glx.FieldDefinition) error {
	v := struct {
		Label     string `json:"label"`
		ValueType string `json:"value_type"`
	}{
		Label:     in.Label,
		ValueType: in.ValueType,
	}
	return json.MarshalEncode(out, &v)
}

// Marshal IntValue as a formatted string instead of a number. Using string
// prevents json.Unmarshal from unmarshaling the value as float64. The special
// format allows Diff to match the IntValue type as well as the value.
func marshalIntValue(out *jsontext.Encoder, in entity.IntValue) error {
	return out.WriteToken(jsontext.String(fmt.Sprintf("IntValue(%d)", in)))
}

func marshalPerson(out *jsontext.Encoder, in *entity.Person) error {
	v := struct {
		Properties map[string]entity.Property `json:"properties"`
	}{
		Properties: in.Properties,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPlace(out *jsontext.Encoder, in *entity.Place) error {
	v := struct {
		Parent     *entity.Place              `json:"parent"`
		Properties map[string]entity.Property `json:"properties"`
		Type       entity.VocabularyValue     `json:"type"`
	}{
		Properties: in.Properties,
		Type:       in.Type,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPropertyDefinition(out *jsontext.Encoder, in *glx.PropertyDefinition) error {
	v := struct {
		Label          string `json:"label"`
		ValueType      string `json:"value_type"`
		ReferenceType  string `json:"reference_type"`
		VocabularyType string `json:"vocabulary_type"`
	}{
		Label:          in.Label,
		ValueType:      in.ValueType,
		ReferenceType:  in.ReferenceType,
		VocabularyType: in.VocabularyType,
	}
	return json.MarshalEncode(out, &v)
}

func marshalStringValue(out *jsontext.Encoder, in entity.StringValue) error {
	return out.WriteToken(jsontext.String(fmt.Sprintf("StringValue{%q}", in)))
}

func marshalVocabularyEntry(out *jsontext.Encoder, in *glx.VocabularyEntry) error {
	v := struct {
		Label string `json:"label"`
	}{
		Label: in.Label,
	}

	return json.MarshalEncode(out, &v)
}
