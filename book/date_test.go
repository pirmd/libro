package book

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	testCases := []struct {
		in   string
		want time.Time
	}{
		{"1976", time.Date(1976, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"1976-01", time.Date(1976, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"01.1976", time.Date(1976, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"011976", time.Date(1976, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"1976-01-17", time.Date(1976, 1, 17, 0, 0, 0, 0, time.UTC)},
		{"1976.01.17", time.Date(1976, 1, 17, 0, 0, 0, 0, time.UTC)},
		{"17/01/1976", time.Date(1976, 1, 17, 0, 0, 0, 0, time.UTC)},
		{"19760117", time.Date(1976, 1, 17, 0, 0, 0, 0, time.UTC)},
		{"17011976", time.Date(1976, 1, 17, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		got, err := ParseTimestamp(tc.in)
		if err != nil {
			t.Errorf("Fail to parse time for %s: %v", tc.in, err)
		}

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("Fail to parse time for %s.\nWant: %v\nGot : %v", tc.in, tc.want, got)
		}
	}
}

func TestNormalizeDate(t *testing.T) {
	testCases := []struct {
		in   string
		want string
	}{
		{"1976", "1976"},
		{"1976-01", "1976-01"},
		{"01.1976", "1976-01"},
		{"011976", "1976-01"},
		{"1976-01-17", "1976-01-17"},
		{"1976.01.17", "1976-01-17"},
		{"17/01/1976", "1976-01-17"},
		{"19760117", "1976-01-17"},
		{"17011976", "1976-01-17"},
	}

	for _, tc := range testCases {
		got := NormalizeDate(tc.in)

		if got != tc.want {
			t.Errorf("Fail to normalize date for %s.\nWant: %v\nGot : %v", tc.in, tc.want, got)
		}
	}
}

func TestComareNormalizedDate(t *testing.T) {
	testCases := []struct {
		in1, in2  string
		want      string
		wantEqual bool
	}{
		{"1976", "1976", "1976", true},
		{"1976", "1976-01", "1976-01", true},
		{"1976-01-17", "1976-01", "1976-01-17", true},
		{"1976-01-17", "1976", "1976-01-17", true},
		{"1976-01-17", "1976-02", "", false},
		{"1976-01", "1982", "", false},
		{"1976-01-17", "1982", "", false},
		{"1976-01-17", "1982-01-17", "", false},
		{"1976", "1982", "", false},
	}

	for _, tc := range testCases {
		got, gotEqual := CompareNormalizedDate(tc.in1, tc.in2)

		if gotEqual != tc.wantEqual {
			t.Errorf("Fail to compare date %s vs. %s: wrong equality.\nWant: %v\nGot : %v", tc.in1, tc.in2, tc.wantEqual, gotEqual)
		}
		if got != tc.want {
			t.Errorf("Fail to compare date %s vs. %s: most precise date is not as expected.\nWant: %v\nGot : %v", tc.in1, tc.in2, tc.want, got)
		}
	}
}
