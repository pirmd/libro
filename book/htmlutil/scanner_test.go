package htmlutil

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/html/atom"
)

const (
	testData = "./testdata"
)

func TestReSpaceSeparatedNames(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{"toto", true},
		{"toto123", true},
		{"toto titi", true},
		{"javascript:alert(1)", false},
	}

	for _, tc := range testCases {
		got := reSpaceSeparatedNames.MatchString(tc.in)
		if got != tc.out {
			t.Errorf("Fail detect unsafe attributes '%s': got %v but should get %v", tc.in, got, tc.out)
		}
	}
}

func TestIsInList(t *testing.T) {
	testCases := []struct {
		in  []string
		out bool
	}{
		{in: []string{"toto"}, out: true},
		{in: []string{"titi", "toto"}, out: true},
		{in: []string{"titi"}, out: false},
	}

	for _, tc := range testCases {
		got := isInList("toto", tc.in)
		if got != tc.out {
			t.Errorf("Fail to identify string 'toto' in list '%v': got %v but should get %v", tc.in, got, tc.out)
		}
	}
}

func TestMinimalScanner(t *testing.T) {
	t.Run("XSS", func(t *testing.T) {
		testScannerWithRule(t, NewMinimalScanner())
	})
}

func TestScannerWithStyle(t *testing.T) {
	testScanner := NewScannerWithStyle()

	t.Run("XSS", func(t *testing.T) {
		testScannerWithRule(t, testScanner)
	})

	t.Run("ShouldPass", func(t *testing.T) {
		testCases := []string{
			`<a href="http://www.google.com"></a>`,
			`<img src="giraffe.gif" />`,
			`<img src="http://www.myspace.com/img.gif"/>`,
			`<span class="foo">Hello World</span>`,
			`<span class="foo bar654">Hello World</span>`,
			`<style type="text/css">body {background:yellow;}</style>`,
			`<div style="color:green">`,
		}

		for i, tc := range testCases {
			issues, err := testScanner.Scan(strings.NewReader(tc))
			if err != nil {
				t.Errorf("[line %d] Scan of '%s' failed: %v", i, tc, err)
			}

			if issues != nil {
				for _, msg := range issues {
					t.Log(msg)
				}

				t.Errorf("[line %d] Scan did not pass for '%s'.", i, tc)
			}
		}
	})
}

func TestPermissiveScanner(t *testing.T) {
	testScanner := NewPermissiveScanner()

	t.Run("XSS", func(t *testing.T) {
		testScannerWithRule(t, testScanner)
	})
}

func TestScanner(t *testing.T) {
	// TestScanner uses a slightly less permissive set of rules compared to
	// ScannerWithStyle for testing purpose.
	testScanner := NewScannerWithStyle()
	testScanner.AllowedTags[atom.A] = []string{"href=__URL_?"}
	testScanner.AllowedCSSFunctions = []string{"not", "local"}
	testScanner.AllowedCSSAtKeywords = []string{"@font-face"}

	t.Run("XSS", func(t *testing.T) {
		testScannerWithRule(t, testScanner)
	})

	t.Run("ShouldPass", func(t *testing.T) {
		testCases := []string{
			`<a href="http://www.google.com"></a>`,
			`<img src="giraffe.gif" />`,
			`<img src="http://www.myspace.com/img.gif"/>`,
			`<span class="foo">Hello World</span>`,
			`<span class="foo bar654">Hello World</span>`,
			`<style type="text/css">body {background:yellow;}</style>`,
			`<style> div.stub *:not([title^="si on"]) { color : lime }</style>`,
			`<style>@font-face { font-family: MyHelvetica; src: local("Helvetica Neue Bold"), url(MgOpenModernaBold.ttf); }</style>`,
			`<div style="color:green">`,

			`<a href="?q=1&r=2"></a>`,
			`<a href="?q=%7B%22value%22%3A%22a%22%7D"></a>`,
			`<a href="?q=1&r=2&s=:foo@"></a>`,
		}

		for i, tc := range testCases {
			issues, err := testScanner.Scan(strings.NewReader(tc))
			if err != nil {
				t.Errorf("[line %d] Scan of '%s' failed: %v", i, tc, err)
			}

			if issues != nil {
				for _, msg := range issues {
					t.Log(msg)
				}

				t.Errorf("[line %d] Scan did not pass for '%s'.", i, tc)
			}
		}
	})
}

func testScannerWithRule(t *testing.T, scanner *Scanner) {
	if err := filepath.Walk(testData, func(path string, fi os.FileInfo, werr error) error {
		if werr != nil {
			return werr
		}

		if fi.IsDir() {
			return nil
		}

		r, err := os.Open(path)
		if err != nil {
			t.Fatalf("cannot read test data in %s: %v", path, err)
		}
		defer r.Close()

		t.Run(fi.Name(), func(t *testing.T) {
			line := 0
			s := bufio.NewScanner(r)
			for s.Scan() {
				line++
				tc := s.Text()
				if tc == "" || strings.HasPrefix(tc, "//") {
					continue // empty line or comment
				}

				issues, err := scanner.Scan(strings.NewReader(tc))
				if err != nil {
					t.Errorf("[line %d] Scan of '%s' failed: %v", line, tc, err)
				}

				if strings.HasPrefix(fi.Name(), "xss") {
					if issues == nil {
						t.Errorf("[line %d] Scan pass for '%s'.", line, tc)
					}
				} else {
					if issues != nil {
						for _, msg := range issues {
							t.Log(msg)
						}

						t.Errorf("[line %d] Scan did not pass for '%s'.", line, tc)
					}
				}
			}
		})
		return nil
	}); err != nil {
		t.Fatalf("cannot read test data in %s: %v", testData, err)
	}
}
