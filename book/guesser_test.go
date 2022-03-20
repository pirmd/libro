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
		{"Sun Company (La compagnie des glaces 25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "ShortTitle": "Sun Company"}},
		{"Sun Company - La compagnie des glaces 25", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "ShortTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces #25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "ShortTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces n°25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "ShortTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces Series 25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "ShortTitle": "Sun Company"}},
		{"Book 25 of La compagnie des glaces", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25"}},
		{"Unknown", nil},
	}

	for _, tc := range testCases {
		if got := seriesGuesser.GuessFrom(tc.in); fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestPathGuesser(t *testing.T) {
	testCases := []struct {
		in  string
		out map[string]string
	}{
		{"GJ Arnaud - [La compagnie des glaces 25] - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Series": "La compagnie des glaces", "SeriesIndex": "25", "Title": "Sun Company", "Language": "FR"}},
		{"my/relative/path/GJ Arnaud - [La compagnie des glaces 25] - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Series": "La compagnie des glaces", "SeriesIndex": "25", "Title": "Sun Company", "Language": "FR"}},
		{"/my/full/path/GJ Arnaud - [La compagnie des glaces 25] - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Series": "La compagnie des glaces", "SeriesIndex": "25", "Title": "Sun Company", "Language": "FR"}},
		{"/my/full/path/GJ Arnaud - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Title": "Sun Company", "Language": "FR"}},
		{"Unknown", nil},
	}

	for _, tc := range testCases {
		if got := pathGuesser.GuessFrom(tc.in); fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestString2Categories(t *testing.T) {
	testCases := []struct {
		in  string
		out []string
	}{
		{"A & B & C", []string{"A", "B", "C"}},
		{"A& B & C", []string{"A", "B", "C"}},
		{"A, B & C", []string{"A", "B", "C"}},
		{"A", []string{"A"}},
	}

	for _, tc := range testCases {
		if got := reCategories.Split(tc.in, -1); fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}
