package book

import (
	"time"
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
	"01-2006",
	"01.2006",
	"01/2006",
	"2006-01-02",
	"2006.01.02",
	"02-01-2006",
	"02.01.2006",
	"02/01/2006",
}

// ParseTime parses a time stamp, trying different time format.
func ParseTime(stamp string) (t time.Time, err error) {
	for _, fmt := range stampFormats {
		if t, err = time.Parse(fmt, stamp); err == nil {
			return
		}
	}

	return
}

// NormalizeDate standardize time stamp format using 2006-01-02 notation.
// If initial date is only a year, or only a year and a month, it does not
// substitute day or month to 01.
func NormalizeDate(stamp string) string {
	t, err := ParseTime(stamp)
	if err != nil {
		return stamp
	}
	switch precision := len(stamp); {
	case precision == 4:
		return t.Format("2006")
	case precision < 8:
		return t.Format("2006-01")
	default:
		return t.Format("2006-01-02")
	}
}
