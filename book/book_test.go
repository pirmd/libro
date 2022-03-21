package book

import (
	"encoding/json"
    "fmt"
	"path/filepath"
	"testing"

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

	Verbose, Debug = verify.NewLogger(t), verify.NewLogger(t)

	out := make([]*Book, len(testCases))
	for i, tc := range testCases {
		b, err := NewFromFile(tc)
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

func TestString2List(t *testing.T) {
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
        if got := reList.Split(tc.in, -1); fmt.Sprint(got) != fmt.Sprint(tc.out) {
            t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
        }
    }
}

func TestNewFromMapWithoutOverride(t *testing.T) {
    testCases := []struct {
        in   *Book
        inM  map[string]string
        out  *Book
    }{
        {
            &Book{},
            map[string]string{
                "Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
            },
            &Book{
                Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Language: "FR",
            },
        },

        {
            &Book{ Title: "Mon père fouettard", Categories: []string{"Biographie"} },
            map[string]string{
                "Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
            },
            &Book{
                Title: "Mon père fouettard", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Categories: []string{"Biographie"}, Language: "FR",
            },
        },

        {
            &Book{Authors: []string{"Mini Moi"}, PublishedDate: "2002"},
            map[string]string{
                "Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
            },
            &Book{
                Title: "Mon père, ce héros", Authors: []string{"Mini Moi"}, PublishedDate: "2002", Language: "FR",
            },
        },
    }

    Verbose, Debug = verify.NewLogger(t), verify.NewLogger(t)

    for _, tc := range testCases {
        if err := tc.in.FromMap(tc.inM, false); err != nil {
            t.Fatalf("fail to update Book: %v", err)
        }

        if failure := verify.Equal(tc.out, tc.in); failure != nil {
            t.Errorf("Update Book from map %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.inM, tc.out, tc.in)
        }
    }
}

func TestNewFromMapWithOverride(t *testing.T) {
    testCases := []struct {
        in   *Book
        inM  map[string]string
        out  *Book
    }{
        {
            &Book{},
            map[string]string{
                "Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
            },
            &Book{
                Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Language: "FR",
            },
        },

        {
            &Book{ Title: "Mon père fouettard", Categories: []string{"Biographie"} },
            map[string]string{
                "Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
            },
            &Book{
                Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, Categories: []string{"Biographie"}, PublishedDate: "1980", Language: "FR",
            },
        },

        {
            &Book{Authors: []string{"Mini Moi"}, PublishedDate: "2002"},
            map[string]string{
                "Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
            },
            &Book{
                Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Language: "FR",
            },
        },
    }

    Verbose, Debug = verify.NewLogger(t), verify.NewLogger(t)

    for _, tc := range testCases {
        if err := tc.in.FromMap(tc.inM, true); err != nil {
            t.Fatalf("fail to update Book: %v", err)
        }

        if failure := verify.Equal(tc.out, tc.in); failure != nil {
            t.Errorf("Update Book from map %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.inM, tc.out, tc.in)
        }
    }
}
