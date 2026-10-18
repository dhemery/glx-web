package entity

import (
	"cmp"

	"github.com/genealogix/glx/go-glx/glxdate"
)

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
	shiftDay   = widthTieBreaker         // Make room for tie-breaker bits.
	shiftMonth = shiftDay + widthDay     // Make room for day and tie-breaker bits.
	shiftYear  = shiftMonth + widthMonth // Make room for month, day, and tie-breaker bits.

	// The maximum value that will fit in the day field. This value
	// compares higher than any actual day.
	afterAnyActualDay = 1<<widthDay - 1

	// The maximum value that will fit in month field. This compares higher
	// than any actual month.
	afterAnyActualMonth = 1<<widthMonth - 1
)

func sortDate(d glxdate.Date) int {
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

func CompareDates(x, y glxdate.Date) int {
	return cmp.Compare(sortDate(x), sortDate(y))
}
