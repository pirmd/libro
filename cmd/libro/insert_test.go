package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pirmd/verify"
)

func TestCreate(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	t.Run("Default", func(t *testing.T) {
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

			if err := testApp.RunInsertSubcmd([]string{bInJSON}); err != nil {
				t.Errorf("Fail to create book for %v: %v", b, err)
			}

		}

		got, err := testApp.List()
		if err != nil {
			t.Fatalf("Fail to read library's status: %v", err)
		}
		fmt.Fprintf(testApp.out, "\nFinal list of books in library:\n%s\n", strings.Join(got, "\n"))

		if failure := verify.MatchGolden(t.Name(), testApp.out.String()); failure != nil {
			t.Fatalf("Output is not as expected.\n%v", failure)
		}
	})
}
