package entity_test

import (
	"testing"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx/glxdate"
)

type dateCompareTest struct {
	x, y string
	want int
}

var dateCompareTests = []dateCompareTest{
	{"2021-12-31", "2021-12-31", 0},
	{"2021-12-31", "2022-01-01", -1},
	{"2022-01-01", "2022-01-02", -1},
	{"2022-01", "2022-01-01", -1}, // Missing day < day 01
	{"2022", "2022-01", -1},       // Missing month < month 01
	{"2022", "2022-01-01", -1},

	{"BEF 2022-01-01", "2022-01-01", -1},
	{"BEF 2022-01", "2022-01", -1},
	{"BEF 2022", "2022", -1},
	{"BEF 2022", "2022-01", -1},       // BEF without month < first month
	{"BEF 2022", "2022-01-01", -1},    // BEF without month < day of year
	{"BEF 2022-01", "2022-01-01", -1}, // BEF without day < first day of month

	{"AFT 2022-01-01", "2022-01-01", 1},
	{"AFT 2022-01", "2022-01", 1},
	{"AFT 2022", "2022", 1},
	{"AFT 2022", "2022-12", 1},       // AFT without month > last month of year
	{"AFT 2022", "2022-12-31", 1},    // AFT without month > last day of year
	{"AFT 2022-01", "2022-01-31", 1}, // AFT without day > last day of month
}

func TestCompareDates(t *testing.T) {
	for _, tc := range dateCompareTests {
		x, _ := glxdate.Parse(tc.x)
		y, _ := glxdate.Parse(tc.y)

		got := entity.CompareDates(x, y)
		if got != tc.want {
			t.Errorf("CompareDates(x=%q, y=%q) = %d, want %d", tc.x, tc.y, got, tc.want)
		}

		got = entity.CompareDates(y, x)
		want := -tc.want
		if got != want {
			t.Errorf("CompareDates(y=%q, x=%q) = %d, want %d", tc.y, tc.x, got, want)
		}
	}
}
