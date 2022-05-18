package book

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/pirmd/verify"
)

func TestFromGooglebooks(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	httpmock := verify.StartMockHTTPResponse(testdata)
	defer httpmock.Stop()

	Verbose, Debug = verify.NewLogger(t), verify.NewLogger(t)

	out := make([]*Book, len(testCases))
	for i, tc := range testCases {
		b, err := NewFromFile(tc)
		if err != nil {
			t.Errorf("Fail to get metadata for %s: %v", tc, err)
		}

		if err := b.FromGooglebooks(3); err != nil {
			t.Errorf("Fail to search (mocked) googlebooks for %s: %v", tc, err)
		}

		out[i] = b
	}

	got, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatalf("Fail to marshal test output to json: %v", err)
	}

	if failure := verify.MatchGolden(t.Name(), string(got)); failure != nil {
		t.Fatalf("Metadata is not as expected:\n%v", failure)
	}
}
