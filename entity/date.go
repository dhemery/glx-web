package entity

import (
	"cmp"

	"github.com/genealogix/glx/go-glx"
	"github.com/genealogix/glx/go-glx/glxdate"
)

// CompareDateStrings compares x and y by computing a datestamp for each and
// comparing the datestamps. It returns
//
//	-1 if DateStringStamp(x) is less than DateStringStamp(y),
//	 0 if DateStringStamp(x) equals DateStringStamp(y),
//	+1 if DateStringStamp(x) is greater than DateStringStamp(y).
//
// See [DateStamp] for details of how the datestamp is computed.
func CompareDateStrings(x, y glx.DateString) int {
	return cmp.Compare(DateStringStamp(x), DateStringStamp(y))
}

// DateStringStamp returns an integer that represents the sort order of the
// date. See [DateStamp] for details of how the datestamp is computed.
func DateStringStamp(date glx.DateString) int {
	d, _ := date.Parse()
	return DateStamp(d)
}

// CompareDates compares x and y by comparing their datestamps.
// It returns
//
//	-1 if DateStamp(x) is less than DateStamp(y),
//	 0 if DateStamp(x) equals DateStamp(y),
//	+1 if DateStamp(x) is greater than DateStamp(y).
//
// See [DateStamp] for details of how the datestamp is computed.
func CompareDates(x, y glxdate.Date) int {
	return cmp.Compare(DateStamp(x), DateStamp(y))
}

const (
	// Bit width of "fields" for month, day, and offset.
	widthOffset = 4
	widthDay    = 5
	widthMonth  = 4

	// Amount to shift each field to make room for the higher-precision
	// fields.
	shiftDay   = widthOffset             // Make room for offset field.
	shiftMonth = shiftDay + widthDay     // Make room for day and offset fields.
	shiftYear  = shiftMonth + widthMonth // Make room for month, day, and offset fields.
)

// DateStamp returns an integer that represents the sort order of the date.
// The datestamp accounts for the year, month, day, and qualifier returned by
// Year, Month, Day, and Qualifier respectively. If date is a range, these
// values represent the start date of the range.
//
// For dates with precision less than [glxdate.PrecisionDay], an omitted month
// or day is filled in with a value that depends on the qualifier. For
// unqualified dates and BEF dates, an omitted month or day is filled in with
// 0. For AFT dates, an omitted month or day is filled in with a value greater
// than any actual month or day.
//
// After omitted values are filled in, the datestamp is computed from resulting
// year, month, day, and qualifier. Two dates with the same year, month, day,
// and qualifier have the same datestamp. For dates that differ only in
// qualifier, the timestamp of an unqualified date is greater than that of a
// BEF date, and less than that of an AFT date.
//
// A quirk of this computation gives a non-intuitive result when the timestamp
// of a date qualified by AFT is compared to that of a date just above it
// qualified by BEF. For example, given
//
//	x, _ := glxdate.Parse("AFT 2021-02-02")
//	y, _ := glxdate.Parse("BEF 2021-02-03")
//	_ = DateStamp(x) < DateStamp(y) // true
//
// DateStamp(x) < DateStamp(y) even though any event after 2021-02-02 must have
// happened after any event before 2021-02-03. The quirk applies similarly for
// adjacent BEF and AFT dates with other precisions, including differing
// precisions. For example:
//
//	x := glxdate.Parse("AFT 2022")
//	y := glxdate.Parse("BEF 2023-01-01")
//	_ = DateStamp(x) < DateStamp(y) // true
//
// Resolving this quirk would require DateStamp to know the last day of every
// month in every year in every possible calendar.
//
// The result of comparing the datestamps of dates from different calendars is
// undefined.
func DateStamp(date glxdate.Date) int {
	var year, month, day, offset int

	switch date.Qualifier() {
	case glxdate.QualifierBefore:
		year, month, day = dateParts(date)
		offset = int(date.Precision())
	case glxdate.QualifierAfter:
		year, month, day = nextDateParts(date)
		offset = 4 + int(date.Precision())
	case glxdate.QualifierNone:
		year, month, day = dateParts(date)
		offset = 8 + int(date.Precision())
	}

	return year<<shiftYear + month<<shiftMonth + day<<shiftDay + offset
}

func nextDateParts(date glxdate.Date) (year, month, day int) {
	year = date.Year()
	if year <= 0 {
		year++
	}
	month, _ = date.Month()
	day, _ = date.Day()
	switch date.Precision() {
	case glxdate.PrecisionDay:
		if day < 31 {
			return year, month, day + 1
		}
		fallthrough
	case glxdate.PrecisionMonth:
		if month < 12 {
			return year, month + 1, 1
		}
		fallthrough
	case glxdate.PrecisionYear:
		return year + 1, 1, 1
	}
	return year, 1, 1
}

func dateParts(date glxdate.Date) (int, int, int) {
	year := date.Year()
	if year <= 0 {
		year++
	}
	month, ok := date.Month()
	if !ok {
		month = 1
	}

	day, ok := date.Day()
	if !ok {
		day = 1
	}
	return year, month, day
}
