package book

import (
	"strings"
	"time"
	"unicode"
)

// stampFormats lists all time formats that are recognized to parse a strings
// representing a time stamp.
var stampFormats = []string{
	time.RFC3339,
	time.RFC850,
	time.ANSIC,
	"2006",
	"2006-01",
	"2006.01",
	"200601",
	"01-2006",
	"01.2006",
	"01/2006",
	"012006",
	"2006.01.02",
	"20060102",
	time.DateOnly,
	"02-01-2006",
	"02.01.2006",
	"02/01/2006",
	"02012006",
	"2006-01-02T15:04:05",
}

// ParseTimestamp parses a time stamp, trying different time format.
func ParseTimestamp(stamp string) (t time.Time, err error) {
	for _, fmt := range stampFormats {
		if t, err = time.Parse(fmt, stamp); err == nil {
			return
		}
	}

	return
}

// NormalizeDate standardizes time stamps format using 2006-01-02 notation.
// If initial date is only a year, or only a year and a month, it does not
// substitute day or month to 01.
func NormalizeDate(stamp string) string {
	t, err := ParseTimestamp(stamp)
	if err != nil {
		return stamp
	}

	switch len(cleanStamp(stamp)) {
	case 4:
		return t.Format("2006")
	case 6:
		return t.Format("2006-01")
	default:
		return t.Format("2006-01-02")
	}
}

// Year get year information from a time stamps.
// Returns empty string if stamp format can be recognized.
func Year(stamp string) string {
	t, err := ParseTimestamp(stamp)
	if err != nil {
		return stamp
	}

	return t.Format("2006")
}

// cleanStamp returns an stamp without any separator or blank.
func cleanStamp(stamp string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, stamp)
}
