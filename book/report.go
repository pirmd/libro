package book

import (
	"fmt"
)

// ReportIssue reports a (possible) issue encountered during Book's processing
// that deserves end-user attention.
func (b *Book) ReportIssue(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	b.ToReview = append(b.ToReview, msg)
	Verbose.Print("warn: " + msg)
}

// ReportSuggestion reports a (possible) better choice for Book's
// attribute(s) that can benefit to end-user.
func (b *Book) ReportSuggestion(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	b.ToReview = append(b.ToReview, msg)
	Verbose.Print(msg)
}
