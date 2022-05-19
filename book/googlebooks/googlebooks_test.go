package googlebooks

import (
	"encoding/json"
	"testing"

	"github.com/pirmd/verify"
)

const (
	testdata = "./testdata"
)

func TestSearchVolume(t *testing.T) {
	testAPI := API{
		MaxResults: 3,
	}

	testCases := []*VolumeInfo{
		{
			Title: "Leibowitz",
		},

		{
			Title:   "Un cantique pour Leibowitz",
			Authors: []string{"Walter M Miller"},
		},

		{
			Title:     "Un cantique pour Leibowitz",
			Publisher: "Gallimard",
		},

		{
			Identifier: []Identifier{{Type: "ISBN_13", Identifier: "9782072477065"}},
		},
	}

	httpmock := verify.StartMockHTTPResponse(testdata)
	defer httpmock.Stop()

	out := make([][]*VolumeInfo, len(testCases))
	for i, tc := range testCases {
		found, err := testAPI.SearchVolume(tc)
		if err != nil {
			t.Errorf("Fail to search (mocked) googlebooks for %v: %v", tc, err)
		}

		out[i] = found
	}

	got, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatalf("Fail to marshal test output to json: %v", err)
	}

	if failure := verify.MatchGolden(t.Name(), string(got)); failure != nil {
		t.Fatalf("SearchVolume is not as expected:\n%v", failure)
	}
}
