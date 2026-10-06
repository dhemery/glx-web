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
			t.Errorf("%s:\n  want %#v,\n   got %#v", tc.args, tc.wantErr, err)

		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s\n  want %#v,\n   got %#v", tc.args, tc.want, got)
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
		t.Errorf("want %q, got %q", want, got)
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
		t.Errorf("\n  want %q,\n   got %q", want, got)
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

		// Verbs specific to PersonName.
		"%<": "prefix",
		"%g": "given",
		"%n": "nickname",
		"%{": "surnameprefix",
		"%S": "surname",
		"%>": "suffix",
	}
	for format, want := range cases {
		got := personName.Formatted(format)
		if got != want {
			t.Errorf("format %q: want %q, got %q", format, want, got)
		}
	}
}
