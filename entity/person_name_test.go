package entity

import (
	"fmt"
	"testing"
)

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
