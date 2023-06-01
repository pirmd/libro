package book

import (
	"io"
	"io/fs"
	"path/filepath"

	"github.com/pirmd/epub"
	"golang.org/x/net/html/atom"

	"github.com/pirmd/libro/book/epubcheck"
	"github.com/pirmd/libro/book/htmlutil"
)

// CheckCompletness assesses whether Book has enough Metadata to be
// identified by the end-user.
func (b *Book) CheckCompleteness() error {
	if b.Title == "" || len(b.Authors) == 0 {
		b.ReportIssue("book has no Title or no Author")
	}

	if len(b.Authors) > 1 {
		b.ReportWarning("book has several Authors. Some might be wrongly considered as book's creator.")
	}

	if b.ISBN == "" || len(b.AlternateISBN) > 0 {
		b.ReportWarning("book ISBN is unknown or has alternate possible values.")
	}

	if b.Publisher == "" || b.PublishedDate == "" {
		b.ReportWarning("book has incomplete publishing information.")
	}

	if (b.Series != "" && b.SeriesIndex == 0) ||
		(b.SeriesIndex != 0 && b.Series == "") ||
		(b.SeriesTitle != "" && (b.SeriesIndex == 0 || b.Series == "")) {
		b.ReportWarning("book seems to belong to a series that is not properly identified.")
	}

	if len(b.Description) < 80 {
		b.ReportWarning("book has no description or a too small description")
	}

	return nil
}

// CheckConformity uses EPUBcheck to verify that the book complies with
// EPUB specification so that it will likely be properly rendered by most reading
// systems.
func (b *Book) CheckConformity() error {
	Debug.Printf("run %s --fatal --error --json - %s", epubcheck.Executable, b.Path)

	report, err := epubcheck.Run(b.Path, "--fatal", "--error", "--warn")
	if err != nil {
		return err
	}

	if l := len(report.Messages); l > 0 {
		for _, m := range report.Messages {
			Verbose.Print(m)
		}

		b.ReportIssue("book's content is likely to have rendering issues (%d findings from EPUBcheck)", l)
	}

	return nil
}

// CheckContentSecurity verifies that Book's content does not contain unsafe
// HTML.
func (b *Book) CheckContentSecurity() error {
	// TODO: change logic for scanning: go through all content
	// (epub.WalkReadingContent), record CSS linked by HTML content in the
	// scanning process then scanCSS

	// TODO: can we be more specific than allowing any attributes, then
	// focusing only on detecting suspicious URL, JS or CSS?
	EPUBScanner := htmlutil.NewPermissiveScanner()
	EPUBScanner.AllowedTags[atom.Img] = []string{"src=__REL_URL", "*"}
	// Add some specific tags and cie encountered in the wild
	EPUBScanner.AllowedTags[atom.Meta] = append([]string{
		"http-equiv=Content-Style-Type",
	}, EPUBScanner.AllowedTags[atom.Meta]...)
	EPUBScanner.AllowedTags[atom.A] = append([]string{
		"tag=**",
	}, EPUBScanner.AllowedTags[atom.A]...)
	EPUBScanner.AllowedTags[atom.Link] = append([]string{
		"tag=**",
	}, EPUBScanner.AllowedTags[atom.Link]...)

	var issues []string
	if err := epub.WalkPublicationResources(b.Path, func(r io.Reader, fi fs.FileInfo) error {
		switch filepath.Ext(fi.Name()) {
		case ".html", ".HTML", ".htm", ".HTM":
			Debug.Printf("scan HTML resources: %s", fi.Name())
			if HTMLissues, err := EPUBScanner.Scan(r); err != nil {
				return err
			} else if len(HTMLissues) > 0 {
				issues = append(issues, HTMLissues...)
			}

		case ".css", ".CSS":
			Debug.Printf("scan CSS resources: %s", fi.Name())
			if CSSissues, err := EPUBScanner.ScanCSS(r); err != nil {
				return err
			} else if len(CSSissues) > 0 {
				issues = append(issues, CSSissues...)
			}

		default:
			return nil
		}

		for _, issue := range issues {
			Verbose.Print(issue)
		}

		return nil
	}); err != nil {
		return err
	}

	if nb := len(issues); nb > 0 {
		b.ReportIssue("book's content contains HTML/CSS with security risks: %d issues detected", nb)
	}

	return nil
}
