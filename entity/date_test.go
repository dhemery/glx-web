package entity_test

import (
	"testing"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx"
)

func TestCompareDates(t *testing.T) {
	tt := [][]glx.DateString{
		// Each slice is ordered as desired.
		{
			// Adjacent dates with each possible qualifier.
			"BEF 2022-05-15",
			"    2022-05-15",
			"AFT 2022-05-15",
			"BEF 2022-05-16",
			"    2022-05-16",
			"AFT 2022-05-16",
		},
		{
			// Across a month boundary. Every possible qualifier
			// and precision from the last date of one month to the
			// first date of the next.
			"BEF 2022-05-31",
			"    2022-05-31",
			"AFT 2022-05-31",
			"AFT 2022-05",
			"BEF 2022-06",
			"    2022-06",
			"BEF 2022-06-01",
			"    2022-06-01",
			"AFT 2022-06-01",
		},
		{
			// Across a year boundary. Every possible qualifier and
			// precision from the last date of one year to the
			// first date of the next.
			"BEF 2021-12-31",
			"    2021-12-31",
			"AFT 2021-12-31",
			"AFT 2021-12",
			"AFT 2021",
			"BEF 2022",
			"    2022",
			"BEF 2022-01",
			"    2022-01",
			"BEF 2022-01-01",
			"    2022-01-01",
			"AFT 2022-01-01",
		},
		{
			// Across an epoch boundary. Every possible qualifier
			// and precision from the final date BCE to the initial
			// date AD.
			"BEF 0001-12-31 BCE",
			"    0001-12-31 BCE",
			"AFT 0001-12-31 BCE",
			"AFT 0001-12 BCE",
			"AFT 0001 BCE",
			"BEF 0001",
			"    0001",
			"BEF 0001-01",
			"    0001-01",
			"BEF 0001-01-01",
			"    0001-01-01",
			"AFT 0001-01-01",
		},
	}

	for _, dateStrings := range tt {
		// Compare each date to itself.
		for _, d := range dateStrings {
			want := 0
			got := entity.CompareDateStrings(d, d)
			if got != want {
				t.Errorf("CompareDateStrings(%q, %q) = %d, want %d", d, d, got, want)
			}

		}

		// Compare each date to each date later in the slice.
		for i, x := range dateStrings[:len(dateStrings)-1] {
			for _, y := range dateStrings[i+1:] {
				got := entity.CompareDateStrings(x, y)
				// The slice is ordered as desired, so we want
				// the result to indicate x < y.
				want := -1
				if got != want {
					t.Errorf("CompareDateStrings(x=%q, y=%q) = %d, want %d", x, y, got, want)
				}

				// Swap the order of the arguments, and so swap
				// the sign of the desired result.
				got = entity.CompareDateStrings(y, x)
				want = -want
				if got != 1 {
					t.Errorf("CompareDateStrings(y=%q, x=%q) = %d, want %d", y, x, got, want)
				}
			}
		}
	}
}
