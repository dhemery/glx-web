package entity

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNewPersonName(t *testing.T) {
	cases := []struct {
		args    []string
		want    *PersonName
		wantErr error
	}{
		{
			args: []string{"prefix", "prefix"},
			want: &PersonName{Prefix: "prefix"},
		},
		{
			args: []string{"given", "given"},
			want: &PersonName{Given: "given"},
		},
		{
			args: []string{"nickname", "nickname"},
			want: &PersonName{Nickname: "nickname"},
		},
		{
			args: []string{"surname_prefix", "surname prefix"},
			want: &PersonName{SurnamePrefix: "surname prefix"},
		},
		{
			args: []string{"surname", "surname"},
			want: &PersonName{Surname: "surname"},
		},
		{
			args: []string{"suffix", "suffix"},
			want: &PersonName{Suffix: "suffix"},
		},
		{
			args: []string{
				"suffix", "x",
				"surname", "s",
				"surname_prefix", "sp",
				"nickname", "n",
				"given", "g",
				"prefix", "p",
			},
			want: &PersonName{
				Prefix:        "p",
				Given:         "g",
				Nickname:      "n",
				SurnamePrefix: "sp",
				Surname:       "s",
				Suffix:        "x",
			},
		},
		{
			args:    []string{"given", "given", "surname"},
			wantErr: ErrNotPairs([]string{"given", "given", "surname"}),
		},
		{
			args:    []string{"given", "given", "monkey", "monkey value"},
			wantErr: ErrUnknownKey("monkey"),
		},
	}

	for _, tc := range cases {
		got, err := NewPersonName(tc.args)
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s\n   got %#v,\n  want %#v", tc.args, err, tc.wantErr)

		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s\n   got %#v,\n  want %#v", tc.args, got, tc.want)
		}
	}
}

func TestPersonNameMerge(t *testing.T) {
	cases := map[string]struct {
		original     *PersonName
		replacements *PersonName
		want         *PersonName
	}{
		"original with no blanks fields": {
			// All fields non-blank
			original: &PersonName{
				Prefix:        "original prefix",
				Given:         "original given",
				Nickname:      "original nickname",
				SurnamePrefix: "original surname prefix",
				Surname:       "original surname",
				Suffix:        "original suffix",
			},
			replacements: &PersonName{
				Prefix:        "replacement prefix",
				Given:         "replacement given",
				Nickname:      "replacement nickname",
				SurnamePrefix: "replacement surname prefix",
				Surname:       "replacement surname",
				Suffix:        "replacement suffix",
			},
			// Same as original. No fields changed.
			want: &PersonName{
				Prefix:        "original prefix",
				Given:         "original given",
				Nickname:      "original nickname",
				SurnamePrefix: "original surname prefix",
				Surname:       "original surname",
				Suffix:        "original suffix",
			},
		},
		"original with all blank fields": {
			original: &PersonName{},
			replacements: &PersonName{
				Prefix:        "replacement prefix",
				Given:         "replacement given",
				Nickname:      "replacement nickname",
				SurnamePrefix: "replacement surname prefix",
				Surname:       "replacement surname",
				Suffix:        "replacement suffix",
			},
			// All values replaced.
			want: &PersonName{
				Prefix:        "replacement prefix",
				Given:         "replacement given",
				Nickname:      "replacement nickname",
				SurnamePrefix: "replacement surname prefix",
				Surname:       "replacement surname",
				Suffix:        "replacement suffix",
			},
		},
		"nil original": {
			original: nil,
			replacements: &PersonName{
				Nickname: "replacement nickname",
			},
			// Same as replacements. No other fields filled in.
			want: &PersonName{
				Nickname: "replacement nickname",
			},
		},
		"nil replacement": {
			original:     &PersonName{Given: "given"},
			replacements: nil,
			// Same as original. No other fields filled in.
			want: &PersonName{Given: "given"},
		},
		"nil original and replacement": {
			original:     nil,
			replacements: nil,
			want:         nil,
		},
	}

	for name, tc := range cases {
		got := tc.original.Merge(tc.replacements)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s\n   got %#v,\n  want %#v", name, got, tc.want)
		}
	}
}

func TestPersonNameString(t *testing.T) {
	var personName = PersonName{
		Prefix:        "prefix",
		Given:         "given",
		Nickname:      "nickname",
		SurnamePrefix: "surnameprefix",
		Surname:       "surname",
		Suffix:        "suffix",
	}

	want := `prefix given "nickname" surnameprefix surname suffix`
	got := personName.String()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPersonNameGoString(t *testing.T) {
	var personName = PersonName{
		Prefix:        "prefix",
		Given:         "given",
		Nickname:      "nickname",
		SurnamePrefix: "surnameprefix",
		Surname:       "surname",
		Suffix:        "suffix",
	}

	want := "&PersonName{" +
		`Prefix: "prefix", ` +
		`Given: "given", ` +
		`Nickname: "nickname", ` +
		`SurnamePrefix: "surnameprefix", ` +
		`Surname: "surname", ` +
		`Suffix: "suffix"` +
		"}"

	got := personName.GoString()
	if got != want {
		t.Errorf("\n   got %q,\n  want %q", got, want)
	}
}

func TestPersonNameFormat(t *testing.T) {
	var personName = PersonName{
		Prefix:        "prefix",
		Given:         "given",
		Nickname:      "nickname",
		SurnamePrefix: "surnameprefix",
		Surname:       "surname",
		Suffix:        "suffix",
	}

	cases := map[string]string{
		// Verbs defined by package fmt.
		"%s": personName.String(),
		"%q": fmt.Sprintf("%q", personName.String()),
		"%v": personName.GoString(),
		"%x": fmt.Sprintf("%x", personName.String()),
		"%X": fmt.Sprintf("%X", personName.String()),

		// Custom verbs specific to PersonName.
		"%<": "prefix",
		"%g": "given",
		"%n": "nickname",
		"%{": "surnameprefix",
		"%S": "surname",
		"%>": "suffix",
		// Custom verbs specific to PersonName.

		"%20.10<": "prefix",
	}

	for format, want := range cases {
		got := personName.Formatted(format)
		if got != want {
			t.Errorf("format %q:\n   got %q,\n  want %q", format, got, want)
		}
	}
}
