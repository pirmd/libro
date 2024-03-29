package book

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pirmd/verify"
)

const (
	testdata      = "../testdata" //Use test data of the main package
	testdataBooks = testdata + "/books"
)

func TestString2List(t *testing.T) {
	testCases := []struct {
		in  string
		out []string
	}{
		{"A & B & C", []string{"A", "B", "C"}},
		{"A& B & C", []string{"A", "B", "C"}},
		{"A et B", []string{"A", "B"}},
		{"A", []string{"A"}},
	}

	for _, tc := range testCases {
		if got := reList.Split(tc.in, -1); fmt.Sprint(got) != fmt.Sprint(tc.out) {
			t.Errorf("Guessing %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestCleanAuthorName(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"Victor Hugo", "Victor Hugo"},
		{"Hugo,Victor", "Victor Hugo"},
		{"Victor HUGO", "Victor Hugo"},
		{"HUGO,Victor", "Victor Hugo"},
	}

	for _, tc := range testCases {
		if got := cleanAuthorName(tc.in); got != tc.out {
			t.Errorf("Clean Author's name %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestNewFromFile(t *testing.T) {
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

func TestNewFromMap(t *testing.T) {
	testCases := []struct {
		in  map[string]string
		out *Book
	}{
		{
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Language: "FR",
				Report: NewReport(),
			},
		},

		{
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Skywalker,Luke", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Language: "FR",
				Report: NewReport(),
			},
		},

		{
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Skywalker,Luke et Mini MOI", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker", "Mini Moi"}, PublishedDate: "1980", Language: "FR",
				Report: NewReport(),
			},
		},
	}

	Verbose, Debug = verify.NewLogger(t), verify.NewLogger(t)

	for _, tc := range testCases {
		got, err := NewFromMap(tc.in)
		if err != nil {
			t.Fatalf("fail to create Book: %v", err)
		}

		if failure := verify.Equal(tc.out, got); failure != nil {
			t.Errorf("Create Book from map %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestCompleteFromMap(t *testing.T) {
	testCases := []struct {
		in  *Book
		inM map[string]string
		out *Book
	}{
		{
			New(),
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Language: "FR",
				Report: NewReport(),
			},
		},

		{
			&Book{Title: "Mon père fouettard", Subject: []string{"Biographie"}, Report: NewReport()},
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père fouettard", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Subject: []string{"Biographie"}, Language: "FR",
				Report: NewReport(),
			},
		},

		{
			&Book{Authors: []string{"Mini Moi"}, PublishedDate: "2002", Report: NewReport()},
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père, ce héros", Authors: []string{"Mini Moi"}, PublishedDate: "2002", Language: "FR",
				Report: NewReport(),
			},
		},
	}

	Verbose, Debug = verify.NewLogger(t), verify.NewLogger(t)

	for _, tc := range testCases {
		if err := tc.in.CompleteFromMap(tc.inM); err != nil {
			t.Fatalf("fail to update Book: %v", err)
		}

		if failure := verify.Equal(tc.out, tc.in); failure != nil {
			t.Errorf("Update Book from map %#v failed:\nWant: %#v\nGot : %#v\n\n", tc.inM, tc.out, tc.in)
		}
	}
}

func TestReplaceFromMap(t *testing.T) {
	testCases := []struct {
		in  *Book
		inM map[string]string
		out *Book
	}{
		{
			New(),
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Language: "FR",
				Report: NewReport(),
			},
		},

		{
			&Book{Title: "La gloire de mon père", Subject: []string{"Biographie"}, Report: NewReport()},
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, Subject: []string{"Biographie"}, PublishedDate: "1980", Language: "FR",
				Report: &Report{
					Issues:       []string{"changed Title from La gloire de mon père to Mon père, ce héros"},
					Warnings:     []string{},
					SimilarBooks: []*Book{},
				},
			},
		},

		{
			&Book{Authors: []string{"Mini Moi"}, PublishedDate: "2002", Report: NewReport()},
			map[string]string{
				"Title": "Mon père, ce héros", "Authors": "Luke Skywalker", "PublishedDate": "1980", "Language": "FR",
			},
			&Book{
				Title: "Mon père, ce héros", Authors: []string{"Luke Skywalker"}, PublishedDate: "1980", Language: "FR",
				Report: &Report{
					Issues:       []string{"changed Authors from [Mini Moi] to [Luke Skywalker]"},
					Warnings:     []string{"changed PublishedDate from 2002 to 1980"},
					SimilarBooks: []*Book{},
				},
			},
		},
	}

	Verbose, Debug = verify.NewLogger(t), verify.NewLogger(t)

	for i, tc := range testCases {
		if err := tc.in.ReplaceFromMap(tc.inM); err != nil {
			t.Fatalf("[tc #%d] fail to update Book: %v", i, err)
		}

		if failure := verify.Equal(tc.out, tc.in); failure != nil {
			t.Errorf("[tc #%d] Update Book from map %#v failed:\nWant: %#v\nGot : %#v\n\n", i, tc.inM, tc.out, tc.in)
		}
	}
}
