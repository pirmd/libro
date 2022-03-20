package book_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/pirmd/libro/book"

	"github.com/pirmd/verify"
)

const (
	testdataBooks = "../testdata/books" //Use test data of the main package
)

func TestFromFile(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	book.Verbose, book.Debug = verify.NewLogger(t), verify.NewLogger(t)

	out := make([]*book.Book, len(testCases))
	for i, tc := range testCases {
		b, err := book.NewFromFile(tc)
		if err != nil {
			t.Errorf("Fail to get metadata for %s: %v", tc, err)
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
