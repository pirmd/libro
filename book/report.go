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

	// Warnings collects messages that report events encountered during Book's
	// processing that might deserve end-user attention.
	Warnings []string `json:",omitempty"`

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
		Warnings:     []string{},
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

// ReportWarning reports a (possible) issue encountered during Book's processing
// that light deserves end-user attention.
func (r *Report) ReportWarning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	r.Warnings = append(r.Warnings, msg)
	Verbose.Print("warn: " + msg)
}

// ReportSimilarBook reports possible similar Book.
func (r *Report) ReportSimilarBook(book *Book) {
	r.SimilarBooks = append(r.SimilarBooks, book)
	Verbose.Printf("found similar book: %#v", book)
}

// HasIssue returns whether Report contains at least one Issue.
func (r Report) HasIssue() bool {
	return len(r.Issues) > 0
}

// HasWarning returns whether Report contains at least one Warning.
func (r Report) HasWarning() bool {
	return len(r.Warnings) > 0
}

// HasSimilarBook returns whether Report contains at least one similar Book.
func (r Report) HasSimilarBook() bool {
	return len(r.SimilarBooks) > 0
}

// NeedReview returns whether Report contains messages that benefit from a
// end-user review.
func (r Report) NeedReview() bool {
	return r.HasIssue() || r.HasWarning() || r.HasSimilarBook()
}
