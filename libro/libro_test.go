package libro_test

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pirmd/libro/book"
	"github.com/pirmd/libro/libro"

	"github.com/pirmd/verify"
)

const (
	testdataBooks = "../testdata/books" //Use test data of the main package
)

type testLibro struct {
	*libro.Libro
	*verify.TestFolder
}

func newTestLibro(tb testing.TB) *testLibro {
	testLog := verify.NewLogger(tb)
	testDir := verify.MustNewTestFolder(tb)

	testLib := libro.New()
	testLib.Root = testDir.Root
	testLib.Verbose, testLib.Debug = testLog, testLog

	book.Verbose, book.Debug = testLog, testLog

	return &testLibro{
		Libro:      testLib,
		TestFolder: testDir,
	}
}

func TestRead(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	t.Run("Default", func(t *testing.T) {
		library := newTestLibro(t)

		out := make([]*book.Book, len(testCases))
		for i, tc := range testCases {
			b, err := library.Read(tc)
			if err != nil {
				t.Errorf("Fail to read information for %s: %v", tc, err)
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
	})

	t.Run("WithGuesser", func(t *testing.T) {
		library := newTestLibro(t)
		library.UseGuesser = true

		out := make([]*book.Book, len(testCases))
		for i, tc := range testCases {
			b, err := library.Read(tc)
			if err != nil {
				t.Errorf("Fail to read information for %s: %v", tc, err)
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
	})
}

func TestCreate(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	t.Run("Default", func(t *testing.T) {
		library := newTestLibro(t)

		for _, tc := range testCases {
			b, err := library.Read(tc)
			if err != nil {
				t.Errorf("Fail to read information for %s: %v", tc, err)
			}

			if err := library.Create(b); err != nil {
				t.Errorf("Fail to create book for %#v: %v", b, err)
			}

		}

		got, err := library.List()
		if err != nil {
			t.Fatalf("Fail to read library's status: %v", err)
		}

		if failure := verify.MatchGolden(t.Name(), strings.Join(got, "\n")); failure != nil {
			t.Fatalf("Library' final state is not as expected:\n%v", failure)
		}
	})

	t.Run("WithGuesser", func(t *testing.T) {
		library := newTestLibro(t)
		library.UseGuesser = true

		for _, tc := range testCases {
			b, err := library.Read(tc)
			if err != nil {
				t.Errorf("Fail to read information for %s: %v", tc, err)
			}

			if err := library.Create(b); err != nil {
				t.Errorf("Fail to create book for %#v: %v", b, err)
			}

		}

		got, err := library.List()
		if err != nil {
			t.Fatalf("Fail to read library's status: %v", err)
		}

		if failure := verify.MatchGolden(t.Name(), strings.Join(got, "\n")); failure != nil {
			t.Fatalf("Library' final state is not as expected:\n%v", failure)
		}
	})
}
