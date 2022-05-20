package book

import (
	"fmt"
)

// Report represents a set of indications about libro's automatic or
// semi-automatic activities to the end-user that deserve attention or
// arbitration.
type Report struct {
	// Issues collects messages that report events encountered during Book's
	// processing that deserve end-user attention.
	Issues []string `json:",omitempty"`

	// SimilarBooks collects alternative Book's metadata that are possibly
	// better or more complete than actual metada set. `libro`is for some reasons
	// usually not sure enough whether they are corresponding to exactly the
	// same book.
	SimilarBooks []*Book `json:",omitempty"`
}

// NewReport creates a new empty Report.
func NewReport() *Report {
	return &Report{
		Issues:       []string{},
		SimilarBooks: []*Book{},
	}
}

// ReportIssue reports a (possible) issue encountered during Book's processing
// that deserves end-user attention.
func (r *Report) ReportIssue(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	r.Issues = append(r.Issues, msg)
	Verbose.Print("warn: " + msg)
}

// ReportSimilarBook reports possible similar Book.
func (r *Report) ReportSimilarBook(book *Book) {
	r.SimilarBooks = append(r.SimilarBooks, book)
	Verbose.Printf("found similar book: %#v", book)
}

// NeedReview returns whether Report contains significant elements for the
// end-user to consider.
func (r Report) NeedReview() bool {
	return len(r.Issues) > 0 || len(r.SimilarBooks) > 0
}
