package book

import (
	"fmt"
	"testing"
)

func TestSeriesGuesser(t *testing.T) {
	testCases := []struct {
		in  string
		out map[string]string
	}{
		{"Sun Company (La compagnie des glaces 25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"Sun Company - La compagnie des glaces 25", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces #25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces n°25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces Series 25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"Book 25 of La compagnie des glaces", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25"}},
		{"Unknown", nil},
	}

	for _, tc := range testCases {
		var got map[string]string

		for _, re := range seriesGuessers {
			got = submatchAsMap(tc.in, re)
			if got != nil {
				break
			}
		}

		if fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestPathGuesser(t *testing.T) {
	testCases := []struct {
		in  string
		out map[string]string
	}{
		{"GJ Arnaud - [La compagnie des glaces 25] - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company", "Language": "FR"}},
		{"my/relative/path/GJ Arnaud - [La compagnie des glaces 25] - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company", "Language": "FR"}},
		{"/my/full/path/GJ Arnaud - [La compagnie des glaces 25] - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company", "Language": "FR"}},
		{"/my/full/path/GJ Arnaud - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Title": "Sun Company", "Language": "FR"}},
		{"Unknown", nil},
	}

	for _, tc := range testCases {
		var got map[string]string

		for _, re := range pathGuessers {
			got = submatchAsMap(tc.in, re)
			if got != nil {
				break
			}
		}
		if fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestTitleCleaner(t *testing.T) {
	testCases := []struct {
		in  string
		out map[string]string
	}{
		{
			"Sun Company / La compagnie des glaces 25",
			map[string]string{"SubTitle": "La compagnie des glaces 25", "Title": "Sun Company"},
		},
		{
			"Sun Company (French Edition)",
			map[string]string{"Title": "Sun Company"},
		},
		{"Unknown", nil},
	}

	for _, tc := range testCases {
		var got map[string]string

		for _, re := range titleCleaners {
			got = submatchAsMap(tc.in, re)
			if got != nil {
				break
			}
		}
		if fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}
