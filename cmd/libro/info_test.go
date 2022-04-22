package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pirmd/verify"
)

func TestRunInfoSubCmd(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	testRunInfoSubcmd := func(args ...string) func(*testing.T) {
		return func(t *testing.T) {
			testApp := newTestApp(t)

			for _, tc := range testCases {
				if err := testApp.RunInfoSubcmd(append(args, tc)); err != nil {
					t.Errorf("Fail to read information for %s: %v", tc, err)
				}

				fmt.Fprintln(testApp.Stdout)
			}

			if failure := verify.MatchGolden(t.Name(), testApp.Out()); failure != nil {
				t.Fatalf("Output is not as expected.\n%v", failure)
			}
		}
	}

	t.Run("Default", func(t *testing.T) {
		testRunInfoSubcmd()(t)
	})

	t.Run("WithPlaintextTemplate", func(t *testing.T) {
		testTmpl := `{{template "plaintext" .}}`
		testRunInfoSubcmd("--format=" + testTmpl)(t)
	})

	t.Run("WithGuesser", func(t *testing.T) {
		testRunInfoSubcmd("-use-guesser")(t)
	})

	t.Run("WithGooglebooks", func(t *testing.T) {
		httpmock := verify.StartMockHTTPResponse(testdata)
		defer httpmock.Stop()

		testRunInfoSubcmd("-use-googlebooks")(t)
	})

	t.Run("WithGooglebooksAndGuesser", func(t *testing.T) {
		httpmock := verify.StartMockHTTPResponse(testdata)
		defer httpmock.Stop()

		testRunInfoSubcmd("-use-guesser", "-use-googlebooks")(t)
	})
}
