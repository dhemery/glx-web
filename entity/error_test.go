package entity

import (
	"bytes"
	"maps"
	"regexp"
	"strings"
	"testing"

	"github.com/genealogix/glx/go-glx"
)

type errorTest struct {
	// File content that passes GLX validation (perhaps with warnings) but
	// includes an anomaly that glx-web must warn about. The test trims
	// leading and trailing space.
	content string
	// Regexp pattern for the desired output.
	want string
}

var errorTests = map[string]errorTest{
	"warning-unknown-property.glx": {
		content: `
persons:
  person-id:
    properties:
      monkey: "monkey value"`,
		want: `persons\[person-id\].Properties.monkey: unknown property`,
	},
}

func TestErrorHandling(t *testing.T) {
	vocabs := glx.StandardVocabularies()
	serializer := glx.NewSerializer(nil)

	for name, tc := range errorTests {
		t.Run(name, func(t *testing.T) {

			files := map[string][]byte{}
			// Prepare the standard vocabularies.
			maps.Copy(files, vocabs)
			files[name] = []byte(strings.TrimSpace(tc.content))

			glxfile, dups, err := serializer.DeserializeMultiFileFromMap(files)
			if err != nil {
				t.Fatal(err)
			}
			if len(dups) > 0 {
				t.Fatal("dups", dups)
			}

			errout := &bytes.Buffer{}
			_ = NewCatalog(glxfile, errout)

			got := errout.String()

			want := regexp.MustCompile(tc.want)
			if !want.MatchString(got) {
				t.Errorf("\n          got: %s\nwant match of: %s", got, tc.want)
			}
		})
	}
}
