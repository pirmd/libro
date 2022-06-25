package book

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestReFindStringSubmatchAsMap(t *testing.T) {
	testRe := regexp.MustCompile(`A:\s(?P<A>[0-9]+)\sB:\s(?P<B>[0-9]+)`)
	testCases := []struct {
		in  string
		out map[string]string
	}{
		{in: "A: 314\nB: 42", out: map[string]string{"A": "314", "B": "42"}},
		{in: "A: 314\nC: 42", out: nil},
	}

	for _, tc := range testCases {
		got := reFindStringSubmatchAsMap(tc.in, testRe)

		if !reflect.DeepEqual(tc.out, got) {
			t.Errorf("Find submatches in '%s' failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}

func TestReFindReaderSubmatchAsMap(t *testing.T) {
	testRe := regexp.MustCompile(`A:\s(?P<A>[0-9]+)\sB:\s(?P<B>[0-9]+)`)
	testCases := []struct {
		in  string
		out []map[string]string
	}{
		{in: "A: 314\nB: 42", out: []map[string]string{{"A": "314", "B": "42"}}},
		{in: "A: 314\nB: 42\nA: 17\nB: 01", out: []map[string]string{{"A": "314", "B": "42"}, {"A": "17", "B": "01"}}},
		{in: "A: 314\nC: 42", out: nil},
	}

	for _, tc := range testCases {
		got := reFindReaderSubmatchAsMap(strings.NewReader(tc.in), testRe)

		if !reflect.DeepEqual(tc.out, got) {
			t.Errorf("Find submatches in '%s' failed:\nWant: %#v\nGot : %#v\n\n", tc.in, tc.out, got)
		}
	}
}
