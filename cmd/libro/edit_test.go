package main

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/pirmd/verify"
)

func TestEdit(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	t.Run("DoNothing", func(t *testing.T) {
		testApp := newTestApp(t)

		for _, tc := range testCases {
			b, err := testApp.Library.Read(tc)
			if err != nil {
				t.Errorf("Fail to read information for %s: %v", tc, err)
			}

			j, err := json.Marshal(b)
			if err != nil {
				t.Errorf("Fail to convert book %s to JSON: %v", tc, err)
			}
			bInJSON := string(j)

			if err := testApp.RunEditSubcmd([]string{"-dont-edit", bInJSON}); err != nil {
				t.Errorf("Fail to edit book %v: %v", b, err)
			}

		}

		if failure := verify.MatchGolden(t.Name(), testApp.out.String()); failure != nil {
			t.Fatalf("Output is not as expected.\n%v", failure)
		}
	})

	t.Run("SetDefault", func(t *testing.T) {
		testApp := newTestApp(t)

		for _, tc := range testCases {
			b, err := testApp.Library.Read(tc)
			if err != nil {
				t.Errorf("Fail to read information for %s: %v", tc, err)
			}

			j, err := json.Marshal(b)
			if err != nil {
				t.Errorf("Fail to convert book %s to JSON: %v", tc, err)
			}
			bInJSON := string(j)

			if err := testApp.RunEditSubcmd([]string{"-dont-edit", "-default", "Categories=libro&testing", bInJSON}); err != nil {
				t.Errorf("Fail to create book for %v: %v", b, err)
			}

		}

		if failure := verify.MatchGolden(t.Name(), testApp.out.String()); failure != nil {
			t.Fatalf("Output is not as expected.\n%v", failure)
		}
	})

	t.Run("SetNew", func(t *testing.T) {
		testApp := newTestApp(t)

		for _, tc := range testCases {
			b, err := testApp.Library.Read(tc)
			if err != nil {
				t.Errorf("Fail to read information for %s: %v", tc, err)
			}

			j, err := json.Marshal(b)
			if err != nil {
				t.Errorf("Fail to convert book %s to JSON: %v", tc, err)
			}
			bInJSON := string(j)

			if err := testApp.RunEditSubcmd([]string{"-dont-edit", "-set", "Categories=libro&testing", bInJSON}); err != nil {
				t.Errorf("Fail to create book for %v: %v", b, err)
			}

		}

		if failure := verify.MatchGolden(t.Name(), testApp.out.String()); failure != nil {
			t.Fatalf("Output is not as expected.\n%v", failure)
		}
	})
}
