package entity_test

import (
	"testing"

	"github.com/dhemery/glx-web/entity"
	"github.com/genealogix/glx/go-glx/glxdate"
)

func TestCompareDates(t *testing.T) {
	tt := [][]string{
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
	}

	for _, dateStrings := range tt {
		var dates []glxdate.Date
		for _, dateString := range dateStrings {
			date, _ := glxdate.Parse(dateString)
			dates = append(dates, date)
		}

		// Compare each date to itself.
		for _, d := range dates {
			want := 0
			got := entity.CompareDates(d, d)
			if got != want {
				t.Errorf("CompareDates(%q, %q) = %d, want %d", d, d, got, want)
			}

		}

		// Compare each date to each date later in the slice.
		for i, x := range dates[:len(dates)-1] {
			for _, y := range dates[i+1:] {
				got := entity.CompareDates(x, y)
				// The slice is ordered as desired, so we want
				// the result to indicate x < y.
				want := -1
				if got != want {
					t.Errorf("CompareDates(x=%q, y=%q) = %d, want %d", x, y, got, want)
				}

				// Swap the order of the arguments, and so swap
				// the sign of the desired result.
				got = entity.CompareDates(y, x)
				want = -want
				if got != 1 {
					t.Errorf("CompareDates(y=%q, x=%q) = %d, want %d", y, x, got, want)
				}
			}
		}
	}
}
