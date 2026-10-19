package entity

import (
	"cmp"

	"github.com/genealogix/glx/go-glx/glxdate"
)

// CompareDates compares two [glxdate.Date] values from the same calendar,
// accounting for range, precision, and qualifiers. The comparison uses the
// year, month, and day returned by Year, Month, and Day respectively. For
// ranged dates, these values represent the start date of the range.
//
// For dates with precision less than [glxdate.PrecisionDay], an omitted month
// or day is filled in with a value that depends on the qualifier. For
// unqualified dates and BEF dates, an omitted month or day is filled in with a
// value less than any actual month or day. For AFT dates, an omitted month
// or day is filled in with a value greater than any actual month or day.
//
// After omitted values are filled in, the year, month, and day of the two
// dates are compared. Two dates with the same year, month, day, and qualifier
// compare equal. For dates that differ only in qualifier, an unqualified date
// compares greater than a BEF date with the same year, month, and day, and
// less than an AFT date with the same year, month, and day.
//
// A quirk of this function gives a non-intuitive result when an AFT date is
// compared to a BEF date just above it. For example, given
//
//	x := glxdate.Parse("AFT 2021-02-02")
//	y := glxdate.Parse("BEF 2021-02-03")
//
// CompareDates(x, y) return -1, even though any event after 2021-02-02 must
// have happened after any event before 2021-02-03. Similar for BEF and AFT
// dates with other precisions, including differing precisions. For example:
//
//	x := glxdate.Parse("AFT 2022")
//	y := glxdate.Parse("BEF 2023-01-01")
//
// Resolving this quirk would require CompareDates to know the details of every
// possible calendar.
//
// The result of comparing dates from different calendars is undefined.
func CompareDates(x, y glxdate.Date) int {
	return cmp.Compare(dateSortValue(x), dateSortValue(y))
}

const (
	// Tie-breakers break ties between otherwise equal dates.
	// tieBreakerBefore < tieBreakerUnqualified < tieBreakerAfter.
	tieBreakerBefore      = 0
	tieBreakerUnqualified = 1
	tieBreakerAfter       = 2

	// Bit width of "fields" for month, day, and tie-breaker.
	widthTieBreaker = 2
	widthDay        = 6 // 6 instead of 5 to allow AFT to use a day value > 31.
	widthMonth      = 4

	// Amount to shift each field to make room for the higher-precision
	// fields.
	shiftDay   = widthTieBreaker         // Make room for tie-breaker field.
	shiftMonth = shiftDay + widthDay     // Make room for day and tie-breaker fields.
	shiftYear  = shiftMonth + widthMonth // Make room for month, day, and tie-breaker fields.

	// The maximum value that will fit in the day field. This value
	// compares higher than any actual day.
	afterAnyActualDay = 1<<widthDay - 1

	// The maximum value that will fit in month field. This compares higher
	// than any actual month.
	afterAnyActualMonth = 1<<widthMonth - 1
)

func dateSortValue(d glxdate.Date) int {
	var tieBreaker int

	// Use BEF and AFT qualifiers to break ties beween otherwise equal
	// dates.
	switch d.Qualifier() {
	case glxdate.QualifierBefore:
		tieBreaker = tieBreakerBefore
	case glxdate.QualifierAfter:
		tieBreaker = tieBreakerAfter
	default:
		tieBreaker = tieBreakerUnqualified
	}

	day, ok := d.Day()
	if !ok {
		// No day is specified. For AFT dates use a "day" that compares
		// higher than any actual day. For other dates use 0, which
		// compares lower.
		day = 0
		if d.Qualifier() == glxdate.QualifierAfter {
			day = afterAnyActualDay
		}
	}
	month, ok := d.Month()
	if !ok {
		// No month is specified. For AFT dates use a "month" that
		// compares higher than any actual month. For other dates use
		// 0, which compares lower.
		month = 0
		if d.Qualifier() == glxdate.QualifierAfter {
			month = afterAnyActualMonth
		}
	}

	return d.Year()<<shiftYear + month<<shiftMonth + day<<shiftDay + tieBreaker
}
