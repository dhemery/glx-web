package entity_test

import (
	"bytes"
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
	"github.com/google/go-cmp/cmp"
	"go.yaml.in/yaml/v3"
	"golang.org/x/tools/txtar"
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
			catalog, err := entity.NewCatalog(glxFile, output)
			if err != nil {
				t.Fatalf("compiling catalog: %s", err)
			}

			gotOutputLines := trimmedLines(output.Bytes())
			unmatchedOutputLines, unmatchedPatterns := matchLines(gotOutputLines, tc.wantOutput)
			for _, l := range unmatchedOutputLines {
				t.Errorf("unexpected output: %q", l)
			}
			for _, p := range unmatchedPatterns {
				t.Errorf("unmatched pattern: %q", p)
			}

			// Convert catalog to map[string]any that can be
			// compared to the structure represented by
			// wantCatalog.
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

// Convert in to an object by marshaling and unmarshaling. Use JSON instead of
// YAML because JSON gives more control over marshaling.
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
