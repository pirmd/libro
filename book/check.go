package book

import (
	"github.com/pirmd/libro/book/epubcheck"
)

// QualityLevel indicate the similarity level between two elements.
type QualityLevel int

const (
	// Bad indicates that quality level is not acceptable
	Bad QualityLevel = iota
	// Poor indicates tat quality level is poor.
	Poor
	// Average indicates that quality level is acceptable
	Average
	// Good indicates that quality level is good.
	Good
	// Perfect indicates that quality level is perfect.
	Perfect
)

// String outputs a human understandable description of a QualityLevel.
func (lvl QualityLevel) String() string {
	return [...]string{"bad", "poor", "average", "good", "as perfect as it can be"}[lvl]
}

// Identifiability assesses whether Book has enough Metadata to be identified
// by the end-user.
func (b *Book) Identifiability() QualityLevel {
	if b.Title == "" || len(b.Authors) == 0 {
		Verbose.Printf("book has no Title or no known Author(s)")
		return Bad
	}

	if b.Publisher == "" || b.PublishedDate == "" {
		if b.ISBN == "" {
			Verbose.Printf("book has a Title and Author(s) but no complete Publication information")
			return Average
		}
	}

	if b.Series != "" && b.SeriesIndex != 0 {
		return Perfect
	}

	return Good
}

// InformationCompleteness assesses whether Book's information is complete.
func (b *Book) InformationCompleteness() QualityLevel {
	if b.Title == "" || len(b.Authors) == 0 {
		Verbose.Printf("book has no Title or Author(s)")
		return Bad
	}

	if b.Description == "" {
		Verbose.Printf("book has no Description")
		return Average
	}

	if b.Publisher == "" || b.PublishedDate == "" {
		Verbose.Printf("book has no complete Publication information")
		return Average
	}

	if b.ISBN == "" {
		Verbose.Printf("book has no ISBN")
		return Average
	}

	if (b.Series != "" && b.SeriesIndex == 0) ||
		(b.SeriesIndex != 0 && b.Series == "") ||
		(b.SeriesTitle != "" && (b.SeriesIndex == 0 || b.Series == "")) {
		Verbose.Printf("book has incomplete Series information")

		return Average
	}

	return Good
}

// CanBeRendered uses EPUBcheck to verify that the book will likely be
// properly rendered by most reading systems.
func (b *Book) CanBeRendered() error {
	Verbose.Print("Verify that Book can be rendered by most reading systems")
	Debug.Printf("run %s --fatal --error --json - %s", epubcheck.Executable, b.Path)

	report, err := epubcheck.Run(b.Path, "--fatal", "--error", "--warn")
	if err != nil {
		return err
	}

	if l := len(report.Messages); l > 0 {
		for _, m := range report.Messages {
			Verbose.Print(m)
		}

		b.ReportIssue("Book's content is likely to have rendering issues (%d findings from epubcheck)", l)
	}

	return nil
}
