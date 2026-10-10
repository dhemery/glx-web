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
	"unknown-property": {
		content: `
persons:
  person-unknown-property:
    properties:
      monkey: "this property is not defined"`,
		want: `persons\[person-unknown-property\].Properties.monkey: unknown property`,
	},
	"unknown-field": {
		content: `
persons:
  person-unknown-field:
    properties:
      name:
        value: Person Unknown Field
        fields:
          monkey: "this field is not defined"`,
		want: `persons\[person-unknown-field\].Properties.name.Fields.monkey: unknown field`,
	},
	"unparseable-property-fields": {
		content: `
persons:
  person-unparseable-fields-type:
    properties:
      name:
        value: Jeff
        fields: "string instead of map of fields"`,

		want: `persons\[person-unparseable-fields-type\].Properties.name.Fields: cannot parse type.*discarding all fields`,
	},
	"primitive-value-wrong-type": {
		content: `
persons:
  person-unparseable-primitive-type:
    properties:
      occupation:
        value: [array, instead, of, primitive, type]`,
		want: `persons\[person-unparseable-primitive-type\].Properties.occupation.Value: cannot parse primitive value type`,
	},
}

func TestErrorHandling(t *testing.T) {
	vocabs := glx.StandardVocabularies()
	serializer := glx.NewSerializer(nil)

	for name, tc := range errorTests {
		t.Run(name, func(t *testing.T) {
			files := map[string][]byte{}

			// Add the standard vocabularies.
			maps.Copy(files, vocabs)

			// Add the test content.
			files[name+".glx"] = []byte(strings.TrimSpace(tc.content))

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
				t.Errorf("\n error output: %s\nwant match of: %s", got, tc.want)
			}
		})
	}
}
