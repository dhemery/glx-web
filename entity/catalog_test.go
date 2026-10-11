package entity

import (
	"bytes"
	"encoding/json/v2"
	"io/fs"
	"maps"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/dhemery/glx-web/internal/marshal"
	"github.com/genealogix/glx/go-glx"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v3"
)

// JSON marshal options to represent the catalog with only the details we want
// to compare. Certain details of GLX vocabularies (e.g. Description) are
// verbose and are not useful for checking the results of NewCatalog.
var catalogMarshalOptions = json.JoinOptions(
	json.FormatNilMapAsNull(true),
	json.FormatNilSliceAsNull(true),
	json.OmitZeroStructFields(true),
	json.WithMarshalers(
		json.JoinMarshalers(
			json.MarshalToFunc(marshal.Date),
			json.MarshalToFunc(marshal.PropertyDefinition),
			json.MarshalToFunc(marshal.FieldDefinition)),
	),
)

// TestNewCatalogWarnings tests that [NewCatalog] emits warnings about
// troublesome conditions not prevented by glx validation. Each test case is a
// [txtar.Archive] read from testdata, contianing a comment, one or more GLX
// files, and a want.yaml file. The GLX files are looaded into a [glx.GLXFile]
// to pass to NewCatalog. The comment is a regexp pattern that the output
// written by NewCatalog must match. The want.yaml file is a YAML
// representation of the desired catalog.
func TestNewCatalogWarnings(t *testing.T) {
	for name, tc := range textArchives(t, "testdata/catalog/warning-tests/*") {
		t.Run(filepath.Base(name), func(t *testing.T) {
			fsys, err := txtar.FS(tc)
			if err != nil {
				t.Fatal(err)
			}

			glxfile := glxFile(t, fsys)

			output := &bytes.Buffer{}
			cat := NewCatalog(glxfile, output)

			gotOutput := output.Bytes()

			wantOutput := regexp.MustCompile(string(bytes.TrimSpace(tc.Comment)))
			if !wantOutput.Match(gotOutput) {
				t.Errorf("\n error output: %s\nwant match of: %s", gotOutput, tc.Comment)
			}

			wantCatObject := unmarshalYAML(t, fsys, "want-catalog.yaml")
			gotCatObject := catAsObject(t, cat)

			if diff := cmp.Diff(wantCatObject, gotCatObject); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func catAsObject(t *testing.T, cat *Catalog) any {
	t.Helper()
	j, err := json.Marshal(cat, catalogMarshalOptions)
	if err != nil {
		t.Fatal("marshaling catalog:", err)
	}

	var out any
	if err := json.Unmarshal(j, &out); err != nil {
		t.Fatal("unmarshaling catalog:", err)
	}
	return out
}

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
