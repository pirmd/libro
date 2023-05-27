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
		{"25.Sun Company", map[string]string{"SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"25 - Sun Company", map[string]string{"SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"#25 - Sun Company", map[string]string{"SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces T25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces, T25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"Sun Company (La compagnie des glaces - T25)", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"La compagnie des glaces (Livre 25) - Sun Company", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"La compagnie des glaces (Livre 25) - Sun Company", map[string]string{"Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"[Cycle d'Ambre-01] Les 9 Princes d'Ambre", map[string]string{"Series": "Cycle d'Ambre", "SeriesIndex": "01", "SeriesTitle": "Les 9 Princes d'Ambre"}},
		{"La Trilogie de l'Empire 1 : Fille de l'Empire", map[string]string{"Series": "La Trilogie de l'Empire", "SeriesIndex": "1", "SeriesTitle": "Fille de l'Empire"}},
		{"Unknown", nil},
	}

	for _, tc := range testCases {
		var got map[string]string

		for _, re := range seriesGuessers {
			got = reFindStringSubmatchAsMap(tc.in, re)
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
		{"my/relative/path/GJ Arnaud - [La compagnie des glaces 25] - Sun Company.epub", map[string]string{"Authors": "GJ Arnaud", "Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company"}},
		{"/my/full/path/GJ Arnaud - [La compagnie des glaces 25] - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Series": "La compagnie des glaces", "SeriesIndex": "25", "SeriesTitle": "Sun Company", "Language": "FR"}},
		{"/my/full/path/GJ Arnaud - Sun Company [FR].epub", map[string]string{"Authors": "GJ Arnaud", "Title": "Sun Company", "Language": "FR"}},
		{"/my/full/path/GJ Arnaud - Sun Company.epub", map[string]string{"Authors": "GJ Arnaud", "Title": "Sun Company"}},
		{
			"Feist, Raymond E. - Cycle de la Guerre de la Faille 04 - Trilogie de l_Empire 1 _ Fille de l_Empire, La.epub",
			map[string]string{"Authors": "Feist, Raymond E.", "Series": "Cycle de la Guerre de la Faille", "SeriesIndex": "04", "SeriesTitle": "Trilogie de l_Empire 1 _ Fille de l_Empire, La"},
		},
		{"Unknown", nil},
	}

	for _, tc := range testCases {
		var got map[string]string

		for _, re := range pathGuessers {
			got = reFindStringSubmatchAsMap(tc.in, re)
			if got != nil {
				break
			}
		}
		if fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestContentGuesser(t *testing.T) {
	testCases := []struct {
		in  string
		out map[string]string
	}{
		{in: `ISBN 978-0-596-52068-7`, out: map[string]string{"ISBN": "978-0-596-52068-7"}},
		{in: `ISBN 978 0 596 52068 7`, out: map[string]string{"ISBN": "978 0 596 52068 7"}},
		{in: `EAN: 9780596520687`, out: map[string]string{"ISBN": "9780596520687"}},
		{in: `ISBN-10: 0-596-52068-9`, out: map[string]string{"ISBN": "0-596-52068-9"}},
		{in: `ISBN : 978-2-7470-9059-9`, out: map[string]string{"ISBN": "978-2-7470-9059-9"}},
		{in: `<br>ISBN 978 0 596 52068 7</br>`, out: map[string]string{"ISBN": "978 0 596 52068 7"}},
		{in: "ISBN 99921-58-10-7", out: map[string]string{"ISBN": "99921-58-10-7"}},
		{in: "ISBN 9971-5-0210-0", out: map[string]string{"ISBN": "9971-5-0210-0"}},
		{in: "ISBN 960-425-059-0", out: map[string]string{"ISBN": "960-425-059-0"}},
		{in: "ISBN 80-902734-1-6", out: map[string]string{"ISBN": "80-902734-1-6"}},
		{in: "ISBN 85-359-0277-5", out: map[string]string{"ISBN": "85-359-0277-5"}},
		{in: "ISBN 1-84356-028-3", out: map[string]string{"ISBN": "1-84356-028-3"}},
		{in: "ISBN 0-684-84328-5", out: map[string]string{"ISBN": "0-684-84328-5"}},
		{in: "ISBN 0-8044-2957-X", out: map[string]string{"ISBN": "0-8044-2957-X"}},
		{in: "ISBN 0-85131-041-9", out: map[string]string{"ISBN": "0-85131-041-9"}},
		{in: "ISBN 93-86954-21-4", out: map[string]string{"ISBN": "93-86954-21-4"}},
		{in: "ISBN 0-943396-04-2", out: map[string]string{"ISBN": "0-943396-04-2"}},
		{in: "ISBN 0-9752298-0-X", out: map[string]string{"ISBN": "0-9752298-0-X"}},
		{in: `̀ISBN 12345`, out: nil},
		{in: `11 – X A12616 ISBN 978-2-07-012616-3 13,90`, out: map[string]string{"ISBN": "978-2-07-012616-3"}},
	}

	for _, tc := range testCases {
		got := reFindStringSubmatchAsMap(tc.in, contentGuesser)

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
			got = reFindStringSubmatchAsMap(tc.in, re)
			if got != nil {
				break
			}
		}
		if fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestTitleCleaners(t *testing.T) {
	testCases := []struct {
		in  string
		out map[string]string
	}{
		{"Point cardinal (LITTERATURE)", map[string]string{"Title": "Point cardinal"}},
		{"Point cardinal (Edition française)", map[string]string{"Title": "Point cardinal"}},
		{"Point cardinal (2017)", map[string]string{"Title": "Point cardinal"}},
		{"Point cardinal (POINT)", map[string]string{"Title": "Point cardinal"}},
	}

	for _, tc := range testCases {
		var got map[string]string

		for _, re := range titleCleaners {
			got = reFindStringSubmatchAsMap(tc.in, re)
			if got != nil {
				break
			}
		}

		if fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}
