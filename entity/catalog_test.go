package entity_test

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"testing"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v3"
)

// TestCatalogCompile tests how [entity.NewCatalog] compiles GLX archives. Each
// test case is a [txtar.Archive] read from testdata. The Archive contains a
// set of files and a comment descibing the responsibilities being tested. The
// files are:
//
//   - One or more .glx files. These are loaded into a [glx.GLXFile], which is
//     passed to NewCatalog.
//   - want-catalog.yaml: A YAML representation of the desired catalog
//     structure. This YAML representation is simplified to omit certain
//     details that are not needed to varify catalog structure.
//   - want-errors: A list of regexp patterns that the output written by
//     NewCatalog must match. Each non-blank output line must match one
//     pattern. Each pattern must match one line of output.
func TestCatalogCompile(t *testing.T) {
	testFS := os.DirFS(".")
	testFileNames, err := fileNamesMatching(testFS, "testdata/catalog/compile/*")
	if err != nil {
		t.Fatalf("test setup: finding test case files: %s", err)
	}

	for _, testFileName := range testFileNames {
		t.Run(filepath.Base(testFileName), func(t *testing.T) {
			tc, err := readCatalogTestCase(testFS, testFileName)
			if err != nil {
				t.Fatalf("test setup: test case: %s", err)
			}

			glxFile, err := loadGLXFile(tc.glxFiles)
			if err != nil {
				if se, ok := errors.AsType[*glx.StructuredValidationError](err); ok {
					// Log the specific errors tucked inside.
					for _, e := range se.Errors {
						t.Log(e)
					}
				}
				t.Fatalf("test setup: loading GLX archive: %s", err)
			}

			output := &bytes.Buffer{}
			catalog := entity.NewCatalog(glxFile, output)

			gotOutputLines := trimmedLines(output.Bytes())
			unmatchedOutputLines, unmatchedPatterns := matchLines(gotOutputLines, tc.wantOutput)
			for _, l := range unmatchedOutputLines {
				t.Errorf("unexpected output: %q", l)
			}
			for _, p := range unmatchedPatterns {
				t.Errorf("unmatched pattern: %q", p)
			}

			// Marshal and unmarshal the catalog to convert it to a
			// map[string]any that can be compared to the structure
			// represented by wantCatalog.
			gotCatalog, err := asObject(catalog, catalogMarshalOptions)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.wantCatalog, gotCatalog); diff != "" {
				t.Error(diff)
			}
		})
	}
}

type catalogTestCase struct {
	glxFiles    map[string][]byte
	wantOutput  []*regexp.Regexp
	wantCatalog any
}

func readCatalogTestCase(fsys fs.FS, fname string) (*catalogTestCase, error) {
	b, err := fs.ReadFile(fsys, fname)
	if err != nil {
		return nil, err
	}

	tarFS, err := txtar.FS(txtar.Parse(b))
	if err != nil {
		return nil, err
	}

	glxFileNames, err := fileNamesMatching(tarFS, "*.glx")
	if err != nil {
		return nil, err
	}

	glxFiles := make(map[string][]byte, len(glxFileNames))
	for _, fname := range glxFileNames {
		b, err := fs.ReadFile(tarFS, fname)
		if err != nil {
			return nil, err
		}
		glxFiles[fname] = b
	}

	outputPatterns, err := fs.ReadFile(tarFS, "want-errors")
	if err != nil {
		return nil, err
	}

	outputMatchers, err := compileRegexps(trimmedLines(outputPatterns))
	if err != nil {
		return nil, err
	}

	catalogYAML, err := fs.ReadFile(tarFS, "want-catalog.yaml")
	if err != nil {
		return nil, err
	}

	var wantCatalog any
	err = yaml.Unmarshal(catalogYAML, &wantCatalog)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling want-catalog.yaml: %w", err)
	}

	return &catalogTestCase{
		glxFiles:    glxFiles,
		wantOutput:  outputMatchers,
		wantCatalog: wantCatalog,
	}, err
}

func fileNamesMatching(fsys fs.FS, pattern string) ([]string, error) {
	names, err := fs.Glob(fsys, pattern)
	if err != nil {
		return nil, fmt.Errorf("globbing %q: %w", pattern, err)
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no files matching %q", pattern)
	}
	return names, nil
}

func matchLines(lines []string, patterns []*regexp.Regexp) ([]string, []*regexp.Regexp) {
	unmatchedPatterns := slices.Clone(patterns)
	unmatchedLines := slices.Clone(lines)

	for p, pattern := range slices.Backward(unmatchedPatterns) {
		for l, line := range slices.Backward(unmatchedLines) {
			if pattern.MatchString(line) {
				unmatchedLines = slices.Delete(unmatchedLines, l, l+1)
				unmatchedPatterns = slices.Delete(unmatchedPatterns, p, p+1)
				break
			}
		}
	}
	return unmatchedLines, unmatchedPatterns
}

func compileRegexps(patterns []string) ([]*regexp.Regexp, error) {
	var rr []*regexp.Regexp

	for _, pattern := range patterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		rr = append(rr, r)
	}

	return rr, nil
}

func trimmedLines(bb []byte) []string {
	var lines []string
	for line := range bytes.Lines(bytes.TrimSpace(bb)) {
		trimmed := string(bytes.TrimSpace(line))
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

// Convert in to an object by marshaling and unmarshaling. We're using instead
// of YAML because JSOM gives more control over marshaling.
func asObject(in any, opts json.Options) (any, error) {
	j, err := json.Marshal(in, opts)
	if err != nil {
		return nil, fmt.Errorf("marshaling: %w", err)
	}

	var out any
	if err := json.Unmarshal(j, &out); err != nil {
		return nil, fmt.Errorf("unmarshaling: %w", err)
	}
	return out, nil
}

func loadGLXFile(files map[string][]byte) (*glx.GLXFile, error) {
	ff := maps.Clone(files)
	// Add the standard vocabularies.
	maps.Copy(ff, glx.StandardVocabularies())

	serializer := glx.NewSerializer(nil)
	glxfile, dups, err := serializer.DeserializeMultiFileFromMap(ff)
	if err != nil {
		return nil, err
	}

	if len(dups) > 0 {
		return nil, fmt.Errorf("duplicates in GLX archive: %v", dups)
	}

	return glxfile, nil
}

// JSON marshal options to reduce the JSON representation of the catalog to
// only the details we want to compare. The resulting HSON omits each entity's
// back-references and inner GLX entity. It also omits details of GLX
// vocabularies (e.g. Description) that are useful for validating catalog
// structure.
var catalogMarshalOptions = json.JoinOptions(
	json.OmitZeroStructFields(true),
	json.WithMarshalers(
		json.JoinMarshalers(
			// Entities
			json.MarshalToFunc(marshalAssertion),
			json.MarshalToFunc(marshalCitation),
			json.MarshalToFunc(marshalEvent),
			json.MarshalToFunc(marshalPerson),
			json.MarshalToFunc(marshalPlace),

			json.MarshalToFunc(marshalParticipant),
			json.MarshalToFunc(marshalPropertyEntry),

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

func marshalAssertion(out *jsontext.Encoder, in *entity.Assertion) error {
	v := struct {
		Date        glxdate.Date        `json:"date"`
		Participant *entity.Participant `json:"participant"`
		Property    string              `json:"property"`
		Subject     string              `json:"subject"`
	}{
		Date:        in.Date,
		Participant: in.Participant,
		Property:    in.Property,
		Subject:     in.Subject.ID(),
	}
	return json.MarshalEncode(out, &v)
}

func marshalBoolValue(out *jsontext.Encoder, in entity.BoolValue) error {
	return out.WriteToken(jsontext.String(fmt.Sprintf("BoolValue{%t}", in)))
}

func marshalCitation(out *jsontext.Encoder, in *entity.Citation) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.PagePath()))
	}
	v := struct {
		Properties map[string]entity.Property `json:"properties"`
	}{
		Properties: in.Properties,
	}
	return json.MarshalEncode(out, &v)
}

func marshalDate(out *jsontext.Encoder, in glxdate.Date) error {
	return out.WriteToken(jsontext.String(in.Raw()))
}

func marshalEvent(out *jsontext.Encoder, in *entity.Event) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.PagePath()))
	}
	v := struct {
		Date         glxdate.Date               `json:"date"`
		Participants []*entity.Participant      `json:"participants"`
		Place        string                     `json:"place"`
		Properties   map[string]entity.Property `json:"properties"`
		Type         *entity.VocabularyValue    `json:"type"`
	}{
		Date:         in.Date,
		Participants: in.Participants,
		Properties:   in.Properties,
		Type:         in.Type,
	}
	if in.Place != nil {
		v.Place = in.Place.ID()
	}
	return json.MarshalEncode(out, &v)
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
	return out.WriteToken(jsontext.String(fmt.Sprintf("IntValue{%d}", in)))
}

func marshalParticipant(out *jsontext.Encoder, in *entity.Participant) error {
	v := struct {
		Person     string                     `json:"person"`
		Role       *entity.VocabularyValue    `json:"role"`
		Properties map[string]entity.Property `json:"properties"`
	}{
		Role:       in.Role,
		Properties: in.Properties,
	}
	if in.Person != nil {
		v.Person = in.Person.ID()
	}
	return json.MarshalEncode(out, &v)
}

func marshalPerson(out *jsontext.Encoder, in *entity.Person) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.PagePath()))
	}
	v := struct {
		Properties map[string]entity.Property `json:"properties"`
	}{
		Properties: in.Properties,
	}
	return json.MarshalEncode(out, &v)
}

func marshalPlace(out *jsontext.Encoder, in *entity.Place) error {
	if out.StackDepth() > 2 { // Marshaling a reference
		return out.WriteToken(jsontext.String(in.PagePath()))
	}
	v := struct {
		Name       string                     `json:"name"`
		Parent     string                     `json:"parent"`
		Properties map[string]entity.Property `json:"properties"`
		Type       *entity.VocabularyValue    `json:"type"`
	}{
		Name:       in.Name,
		Properties: in.Properties,
		Type:       in.Type,
	}
	if in.Parent != nil {
		v.Parent = in.Parent.ID()
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

func marshalPropertyEntry(out *jsontext.Encoder, in entity.PropertyEntry) error {
	v := struct {
		Date   glxdate.Date                    `json:"date"`
		Fields map[string]entity.PropertyField `json:"fields,omitempty"`
		Value  entity.Stringer                 `json:"value"`
	}{
		Date:   in.Date,
		Fields: in.Fields,
		Value:  in.Value,
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
