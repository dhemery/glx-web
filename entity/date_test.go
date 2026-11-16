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
			"BEF 2020-05-15",
			"    2020-05-15",
			"BEF 2020-05-16",
			"AFT 2020-05-15",
			"    2020-05-16",
			"AFT 2020-05-16",
		},
		{
			// Across a month boundary. Every possible qualifier
			// and precision from the last date of one month to the
			// first date of the next.
			"BEF 2021-05-31",
			"    2021-05-31",
			"BEF 2021-06",
			"BEF 2021-06-01",
			"AFT 2021-05",
			"AFT 2021-05-31",
			"    2021-06",
			"    2021-06-01",
			"AFT 2021-06-01",
		},
		{
			// Across a year boundary. Every possible qualifier and
			// precision from the last date of one year to the
			// first date of the next.
			"BEF 2022-12-31",
			"    2022-12-31",
			"BEF 2023",
			"BEF 2023-01",
			"BEF 2023-01-01",
			"AFT 2022",
			"AFT 2022-12",
			"AFT 2022-12-31",
			"    2023",
			"    2023-01",
			"    2023-01-01",
			"AFT 2023-01-01",
		},
		{
			// Across the end of a short month. Every possible
			// qualifier and precision from the last date of one
			// month to the first date of the next.
			//
			// This demonstrates the quirk in the algorithm: If a
			// month has fewer than 31 days, a date AFT the last
			// day of the month sorts too early.
			"BEF 2024-02-29",
			"    2024-02-29",
			"AFT 2024-02-29", // Sorts as if "at least 2024-02-30"
			"BEF 2024-03",
			"BEF 2024-03-01",
			"AFT 2024-02",
			"    2024-03",
			"    2024-03-01",
			"AFT 2024-03-01",
		},
		{
			// Across an epoch boundary. Every possible qualifier
			// and precision from the final date BCE to the initial
			// date AD.
			"BEF 0001-12-31 BCE",
			"    0001-12-31 BCE",
			"BEF 0001",
			"BEF 0001-01",
			"BEF 0001-01-01",
			"AFT 0001 BCE",
			"AFT 0001-12 BCE",
			"AFT 0001-12-31 BCE",
			"    0001",
			"    0001-01",
			"    0001-01-01",
			"AFT 0001-01-01",
		},
		{
			"", // Zero date comes before earliest valid date
			"BEF 9999 BCE",
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
				// 3he slice is ordered as desired, so we want
				// the result to indicate x < y.
				want := -1
				got := entity.CompareDateStrings(x, y)
				if got != want {
					t.Errorf("CompareDateStrings(x=%q, y=%q) = %d, want %d", x, y, got, want)
				}

				// Swap the order of the arguments, and so swap
				// the sign of the desired result.
				want = -want
				got = entity.CompareDateStrings(y, x)
				if got != 1 {
					t.Errorf("CompareDateStrings(y=%q, x=%q) = %d, want %d", y, x, got, want)
				}
			}
		}
	}
}
