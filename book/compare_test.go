package book

import (
	"testing"
)

func TestCompareNormalizedDates(t *testing.T) {
	testCases := []struct {
		in1, in2 string
		want     SimilarityLevel
	}{
		{"1976", "1976", AreTheSame},
		{"1976", "1976-01", AreAlmostTheSame},
		{"1976-01-17", "1976-01", AreAlmostTheSame},
		{"1976-01-17", "1976", AreAlmostTheSame},
		{"1976-01-17", "1976-02", AreNotTheSame},
		{"1976-01", "1982", AreNotTheSame},
		{"1976-01-17", "1982", AreNotTheSame},
		{"1976-01-17", "1982-01-17", AreNotTheSame},
		{"1976", "1982", AreNotTheSame},
	}

	for _, tc := range testCases {
		if got := compareNormalizedDates(tc.in1, tc.in2); got != tc.want {
			t.Errorf("Fail to compare date %s vs. %s: wrong equality.\nWant: %v\nGot : %v", tc.in1, tc.in2, tc.want, got)
		}
	}
}
