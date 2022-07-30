package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pirmd/verify"
)

const (
	testdata      = "./testdata" //Use test data of the main package
	testdataBooks = testdata + "/books"
)

type testApp struct {
	*App
	*verify.TestFolder
}

func newTestApp(tb testing.TB) *testApp {
	app := NewApp()

	app.Stdout = new(bytes.Buffer)

	testLog := verify.NewLogger(tb)
	app.Verbose, app.Debug = testLog, testLog

	return &testApp{
		App:        app,
		TestFolder: verify.MustNewTestFolder(tb),
	}
}

func TestRunInfoSubcmd(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	testRunInfoSubcmd := func(args ...string) func(*testing.T) {
		args = append([]string{"-format={{toPrettyJSON .}}", "info"}, args...)

		return func(t *testing.T) {
			testApp := newTestApp(t)

			for _, tc := range testCases {
				if err := testApp.Run(append(args, tc)); err != nil {
					t.Errorf("Fail to read information for %s: %v", tc, err)
				}

				fmt.Fprintln(testApp.Stdout)
			}

			got := testApp.Stdout.(*bytes.Buffer).String()
			if failure := verify.MatchGolden(t.Name(), got); failure != nil {
				t.Fatalf("Output is not as expected.\n%v", failure)
			}
		}
	}

	t.Run("Default", func(t *testing.T) {
		testRunInfoSubcmd()(t)
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

func TestRunInsertSubcmd(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	testRunInsertSubcmd := func(args ...string) func(*testing.T) {
		args = append([]string{"-format={{toPrettyJSON .}}", "insert"}, args...)

		return func(t *testing.T) {
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

				if err := testApp.Run(append(args, "-root="+testApp.TestFolder.Root, bInJSON)); err != nil {
					t.Errorf("Fail to create book for %v: %v", b, err)
				}

				fmt.Fprintln(testApp.Stdout)
			}

			ls, err := testApp.List()
			if err != nil {
				t.Fatalf("Fail to read library's status: %v", err)
			}
			fmt.Fprintf(testApp.Stdout, "\nFinal list of books in library:\n%s\n", strings.Join(ls, "\n"))

			got := testApp.Stdout.(*bytes.Buffer).String()
			if failure := verify.MatchGolden(t.Name(), got); failure != nil {
				t.Fatalf("Output is not as expected.\n%v", failure)
			}
		}
	}

	t.Run("Default", func(t *testing.T) {
		testRunInsertSubcmd()(t)
	})

	t.Run("WithRenameByAuthorTemmplate", func(t *testing.T) {
		testRunInsertSubcmd(`-rename={{template "fullname_byauthor.gotmpl" .}}`)(t)
	})

	t.Run("WithRenameShortByAuthorTemmplate", func(t *testing.T) {
		testRunInsertSubcmd(`-rename={{template "shortname_byauthor.gotmpl" .}}`)(t)
	})
}

func TestRunCheckSubcmd(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	testRunCheckSubcmd := func(args ...string) func(*testing.T) {
		args = append([]string{"-format={{toPrettyJSON .}}", "check"}, args...)

		return func(t *testing.T) {
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

				if err := testApp.Run(append(args, bInJSON)); err != nil {
					fmt.Fprintf(testApp.Stdout, "\nERROR during check: %v\n", err)
				}

				fmt.Fprintln(testApp.Stdout)
			}

			got := testApp.Stdout.(*bytes.Buffer).String()
			if failure := verify.MatchGolden(t.Name(), got); failure != nil {
				t.Fatalf("Output is not as expected.\n%v", failure)
			}
		}
	}

	t.Run("Default", func(t *testing.T) {
		testRunCheckSubcmd()(t)
	})

	t.Run("WithCompletenessCheck", func(t *testing.T) {
		testRunCheckSubcmd("-completeness")(t)
	})

	t.Run("WithConformityCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test of conformity check in short mode.")
		}

		testRunCheckSubcmd("-conformity")(t)
	})

	t.Run("WithExitIfIssue", func(t *testing.T) {
		testRunCheckSubcmd("-fail-on-issue")(t)
	})
}

func TestRunEditSubcmd(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	testRunEditSubcmd := func(args ...string) func(*testing.T) {
		args = append([]string{"-format={{toPrettyJSON .}}", "edit"}, args...)

		return func(t *testing.T) {
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

				if err := testApp.Run(append(args, bInJSON)); err != nil {
					t.Errorf("Fail to edit book for %v: %v", b, err)
				}

				fmt.Fprintln(testApp.Stdout)
			}

			got := testApp.Stdout.(*bytes.Buffer).String()
			if failure := verify.MatchGolden(t.Name(), got); failure != nil {
				t.Fatalf("Output is not as expected.\n%v", failure)
			}
		}
	}

	t.Run("DoNothing", func(t *testing.T) {
		testRunEditSubcmd("-dont-edit")(t)
	})

	t.Run("SetDefault", func(t *testing.T) {
		testRunEditSubcmd("-dont-edit", "-default", "Subject=libro&testing")(t)
	})

	t.Run("SetNew", func(t *testing.T) {
		testRunEditSubcmd("-dont-edit", "-set", "Subject=libro&testing")(t)
	})
}

func TestBookTemplates(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataBooks, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s: %v", testdataBooks, err)
	}

	t.Run("WithPlaintextTemplate", func(t *testing.T) {
		args := []string{`-format={{template "book.txt.gotmpl" .}}`, "info"}

		testApp := newTestApp(t)

		for _, tc := range testCases {
			if err := testApp.Run(append(args, tc)); err != nil {
				t.Errorf("Fail to read information for %s: %v", tc, err)
			}

			fmt.Fprintln(testApp.Stdout)
		}

		got := testApp.Stdout.(*bytes.Buffer).String()
		if failure := verify.MatchGolden(t.Name(), got); failure != nil {
			t.Fatalf("Output is not as expected.\n%v", failure)
		}
	})

}
