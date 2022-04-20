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
	"2006-01-02",
	"2006.01.02",
	"20060102",
	"02-01-2006",
	"02.01.2006",
	"02/01/2006",
	"02012006",
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

// CompareNormalizedDate compares two 'normalized' date and returns the most
// 'precise' one in case of equality.
func CompareNormalizedDate(date1, date2 string) (string, bool) {
	d1, d2 := date1, date2
	if len(date1) < len(date2) {
		d1, d2 = date2, date1
	}

	if strings.HasPrefix(d1, d2) {
		return d1, true
	}
	return "", false
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
