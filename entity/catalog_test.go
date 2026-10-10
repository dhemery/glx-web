package entity

import (
	"bytes"
	"maps"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/genealogix/glx/go-glx"
	"golang.org/x/tools/txtar"
)

// TestNewCatalogWarnings tests that [NewCatalog] emits warnings about
// troublesome conditions not prevented by glx validation. Each test case is a
// [txtar.Archive] read from testdata. Its Files are GLX files to load into a
// [glx.GLXFile] to pass to NewCatalog. Its Comment is a regexp pattern that
// the output written by NewCatalog must match.
func TestNewCatalogWarnings(t *testing.T) {
	for name, tc := range textArchives(t, "testdata/catalog/warning-tests/*") {
		t.Run(filepath.Base(name), func(t *testing.T) {
			glxfile := glxFile(t, tc.Files)

			output := &bytes.Buffer{}
			_ = NewCatalog(glxfile, output)

			got := output.Bytes()

			want := regexp.MustCompile(string(bytes.TrimSpace(tc.Comment)))
			if !want.Match(got) {
				t.Errorf("\n error output: %s\nwant match of: %s", got, tc.Comment)
			}
		})
	}
}

func glxFile(t *testing.T, ff []txtar.File) *glx.GLXFile {
	t.Helper()
	files := map[string][]byte{}

	// Add the standard vocabularies.
	maps.Copy(files, glx.StandardVocabularies())

	// Add files from the test case.
	for _, f := range ff {
		files[f.Name] = f.Data
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
		t.Fatalf("globbing %s: %s", pattern, err)
	}

	if len(fnames) == 0 {
		t.Fatal("no files matching", pattern)
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
