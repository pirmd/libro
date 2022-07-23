package book

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Verbose is the logger of book package that provides feedback of operation
	// done on books.
	Verbose = log.New(io.Discard, log.Prefix(), log.Flags())

	// Debug is the logger of book package that provides information for
	// debugging purpose.
	Debug = log.New(io.Discard, log.Prefix(), log.Flags())

	// ErrUnknownFormat is raised if supplied file format is unknown.
	ErrUnknownFormat = errors.New("unknown file format")

	// reList is a regexp that splits a list of values (like Authors or Subject).
	reList = regexp.MustCompile(`\s?&\s?`)

	// reAuthName is a regexp that splits an Author's name into its surname and
	// forename.
	reAuthName = regexp.MustCompile(`\s?,\s?`)
)

// Book represents a book.
type Book struct {
	// Path is the location of the book's file in the file-system.
	Path string

	// Title is the book's title.
	Title string

	// Authors is the list names of the authors and/or editors for this book.
	Authors []string

	// ISBN is the unique industry standard identifier for this book.
	// libro tends to prefer ISBN_13 format when available or when it can be
	// derived from an ISBN_10.  ISBN10 and ISBN13 methods can be invoked to
	// convert from one format to the other.
	// Most Book's functions dealing with ISBN will better work if ISBN is
	// 'normalized' using Book.SetISBN.
	ISBN string `json:",omitempty"`

	// SubTitle is the book's sub-title.
	SubTitle string `json:",omitempty"`

	// Publisher is the publisher of this book.
	Publisher string `json:",omitempty"`

	// PublishedDate is the date of publication of this book.
	// `libro` tries to normalize dates using '2006-01-02' format. When
	// 'precision' of date is not enough to capture known month or days, date
	// is cut to '2006-01' or simply to '2006'.
	// Most Book's functions dealing with Date will better work if Date is
	// 'normalized' using Book.SetPublishedDate.
	PublishedDate string `json:",omitempty"`

	// Description is the synopsis of the book. The text of the description
	// is formatted in HTML and includes simple formatting elements.
	Description string `json:",omitempty"`

	// Series is the series to which this book belongs to.
	Series string `json:",omitempty"`

	// SeriesIndex is the position in the series to which the book belongs to.
	SeriesIndex float64 `json:",omitempty"`

	// SeriesTitle is the book's title in the series (without Series nor
	// SubTitle information).
	SeriesTitle string `json:",omitempty"`

	// Language is the book's language. It is the two-letter ISO 639-1 code
	// such as 'fr', 'en'.
	Language string `json:",omitempty"`

	// PageCount is total number of pages of this book.
	PageCount int64 `json:",omitempty"`

	// Subject is the list of subject categories, such as "Fiction",
	// "Suspense".
	Subject []string `json:",omitempty"`

	// ToReview collects messages that report events encountered during Book's
	// processing that deserve end-user attention.
	*Report
}

// New creates a new Book
func New() *Book {
	return &Book{
		Report: NewReport(),
	}
}

// NewFromFile creates a new Book and populates its information according to
// the file's metadata.
func NewFromFile(path string) (*Book, error) {
	Debug.Printf("create a Book from '%s'", path)

	switch ext := filepath.Ext(path); ext {
	case ".epub":
		return NewFromEpub(path)

	default:
		return nil, ErrUnknownFormat
	}
}

// SetISBN sets Book's ISBN and tries to normalize it to ISBN_13 format.
// SetISBN reports non-recognized ISBN but do not fail.
func (b *Book) SetISBN(isbn string) {
	normISBN, err := NormalizeISBN(isbn)
	if err != nil {
		b.ReportIssue("non-supported ISBN (%s): %v", isbn, err)
		return
	}

	b.ISBN = normISBN
}

// SetPublishedDate sets Book's PublishedDate and tries to normalize its
// format.
func (b *Book) SetPublishedDate(date string) {
	b.PublishedDate = NormalizeDate(date)
}

// SetAuthors sets Book's Authors and tries to keep Authors' names and
// surnames in a pre-defined order.
func (b *Book) SetAuthors(authors []string) {
	b.Authors = make([]string, len(authors))

	for i, auth := range authors {
		if name := reAuthName.Split(auth, 2); len(name) == 2 {
			b.Authors[i] = strings.TrimSpace(name[1] + " " + name[0])
		} else {
			b.Authors[i] = auth
		}
	}
}

// SetDescription sets Book's Description and tries to clean it from un-helping
// HTML formatting directives.
func (b *Book) SetDescription(desc string) {
	// TODO: use a better 'html to text' converter that preserves formatting
	unhtml, err := getRawTextFromHTML(strings.NewReader(desc))
	if err != nil {
		Debug.Printf("fail to clean Description from HTML tags: %v", err)
		b.Description = desc
	}

	cleanDesc, err := io.ReadAll(unhtml)
	if err != nil {
		Debug.Printf("fail to clean Description from HTML tags: %v", err)
		b.Description = desc
	}

	b.Description = string(cleanDesc)
}

// NewFromMap creates a Book's from to the attributes defined as a map
// where keys are attribute's name (insensitive to case) and value is a string
// representation of the attribute's value.
// For attributes that accept a list of values (like Authors or Subject),
// provided map value should be formatted like "val0 & val1" (individual value
// in as string separated by '&').
func NewFromMap(m map[string]string) (*Book, error) {
	b := New()

	for attr, value := range m {
		switch a := strings.Title(attr); a {
		case "Title":
			b.Title = value

		case "SubTitle":
			b.SubTitle = value

		case "SeriesTitle":
			b.SeriesTitle = value

		case "Authors":
			b.SetAuthors(reList.Split(value, -1))

		case "Publisher":
			b.Publisher = value

		case "PublishedDate":
			b.SetPublishedDate(value)

		case "Description":
			b.SetDescription(value)

		case "Series":
			b.Series = value

		case "SeriesIndex":
			var err error
			if b.SeriesIndex, err = strconv.ParseFloat(value, 32); err != nil {
				return nil, fmt.Errorf("cannot assign %s to '%s': %v", value, a, err)
			}

		case "ISBN":
			b.SetISBN(value)

		case "Language":
			b.Language = value

		case "PageCount":
			var err error
			if b.PageCount, err = strconv.ParseInt(value, 10, 0); err != nil {
				return nil, fmt.Errorf("cannot assign %s to '%s': %v", value, a, err)
			}

		case "Subject":
			b.Subject = reList.Split(value, -1)

		default:
			return nil, fmt.Errorf("cannot set unknown attribute '%s'", a)
		}
	}

	return b, nil
}

// MergeWith merges Book with 'b1' Book.
// If override is set, Book's attributes are replaced by the none-empty
// corresponding attribute of 'b1' Book.
func (b *Book) MergeWith(b1 *Book, override bool) {
	if b1.Title != "" {
		if b.Title == "" {
			Verbose.Printf("set empty Title to %v", b1.Title)
			b.Title = b1.Title
		} else if override {
			if b.compareTitleWith(b1) < AreAlmostTheSame {
				b.ReportIssue("changed Title from %v to %v", b.Title, b1.Title)
			} else {
				Verbose.Printf("changed Title from %v to %v", b.Title, b1.Title)
			}

			b.Title = b1.Title
		}
	}

	if len(b1.Authors) > 0 {
		if len(b.Authors) == 0 {
			Verbose.Printf("set empty Authors to %v", b1.Authors)
			b.Authors = append([]string{}, b1.Authors...)
		} else if override {
			if b.compareAuthorsWith(b1) < AreAlmostTheSame {
				b.ReportIssue("changed Authors from %v to %v", b.Authors, b1.Authors)
			} else {
				Verbose.Printf("changed Authors from %v to %v", b.Authors, b1.Authors)
			}

			b.Authors = append([]string{}, b1.Authors...)
		}
	}

	if b1.ISBN != "" {
		if b.ISBN == "" {
			b.ReportIssue("set empty ISBN to %v", b1.ISBN)
			b.ISBN = b1.ISBN
		} else if override && b.compareIdentifierWith(b1) != AreTheSame {
			b.ReportIssue("changed ISBN from %v to %v", b.ISBN, b1.ISBN)
			b.ISBN = b1.ISBN
		} else if b.compareIdentifierWith(b1) != AreTheSame {
			b.ReportIssue("found a different ISBN: %v (vs. %s)", b1.ISBN, b.ISBN)
		}
	}

	if b1.SubTitle != "" {
		if b.SubTitle == "" {
			Verbose.Printf("set empty SubTitle to %s", b1.SubTitle)
			b.SubTitle = b1.SubTitle
		} else if override {
			if b.compareSubTitleWith(b1) < AreAlmostTheSame {
				b.ReportIssue("changed SubTitle from %v to %v", b.SubTitle, b1.SubTitle)
			} else {
				Verbose.Printf("changed SubTitle from %v to %v", b.SubTitle, b1.SubTitle)
			}
			b.SubTitle = b1.SubTitle
		}
	}

	if b1.Publisher != "" {
		if b.Publisher == "" {
			Verbose.Printf("set empty Publisher to %s", b1.Publisher)
			b.Publisher = b1.Publisher
		} else if override {
			if b.comparePublisherWith(b1) < AreAlmostTheSame {
				b.ReportIssue("changed Publisher from %v to %v", b.Publisher, b1.Publisher)
			} else {
				Verbose.Printf("changed Publisher from %v to %v", b.Publisher, b1.Publisher)
			}
			b.Publisher = b1.Publisher
		}
	}

	if b1.PublishedDate != "" {
		if b.PublishedDate == "" {
			Verbose.Printf("set empty PublishedDate to %s", b1.PublishedDate)
			b.PublishedDate = b1.PublishedDate
		} else if override {
			if b.comparePublishedDateWith(b1) < AreAlmostTheSame {
				b.ReportIssue("changed PublishedDate from %v to %v", b.PublishedDate, b1.PublishedDate)
			} else {
				Verbose.Printf("changed PublishedDate from %v to (more precise) %v", b.PublishedDate, b1.PublishedDate)
			}
			b.PublishedDate = b1.PublishedDate
		} else if b.comparePublishedDateWith(b1) == AreAlmostTheSame && len(b1.PublishedDate) > len(b.PublishedDate) {
			Verbose.Printf("changed PublishedDate from %v to (more precise) %v", b.PublishedDate, b1.PublishedDate)
			b.PublishedDate = b1.PublishedDate
		}
	}

	if b1.Description != "" {
		if b.Description == "" {
			Verbose.Printf("set empty Description to %.12v", b1.Description)
			b.Description = b1.Description
		} else if override && !strings.EqualFold(b.Description, b1.Description) {
			Verbose.Printf("changed Description from %.12v to %.12v", b.Description, b1.Description)
			b.Description = b1.Description
		}
	}

	if b1.Series != "" {
		if b.Series == "" {
			Verbose.Printf("set empty Series to %v", b1.Series)
			b.Series = b1.Series
		} else if override && !strings.EqualFold(b.Series, b1.Series) {
			Verbose.Printf("changed Series from %v to %v", b.Series, b1.Series)
			b.Series = b1.Series
		}
	}

	if b1.SeriesIndex != 0 {
		if b.SeriesIndex == 0 {
			Verbose.Printf("set empty SeriesIndex to %v", b1.SeriesIndex)
			b.SeriesIndex = b1.SeriesIndex
		} else if override && (b.SeriesIndex != b1.SeriesIndex) {
			b.ReportIssue("changed SeriesIndex from %v to %v", b.SeriesIndex, b1.SeriesIndex)
			b.SeriesIndex = b1.SeriesIndex
		}
	}

	if b1.SeriesTitle != "" {
		if b.SeriesTitle == "" {
			Verbose.Printf("set empty SeriesTitle to %v", b1.SeriesTitle)
			b.SeriesTitle = b1.SeriesTitle
		} else if override && !strings.EqualFold(b.SeriesTitle, b1.SeriesTitle) {
			b.ReportIssue("changed SeriesTitle from %v to %v", b.SeriesTitle, b1.SeriesTitle)
			b.SeriesTitle = b1.SeriesTitle
		}
	}

	if b1.Language != "" {
		if b.Language == "" {
			Verbose.Printf("set empty Language to %v", b1.Language)
			b.Language = b1.Language
		} else if override && !strings.EqualFold(b.Language, b1.Language) {
			Verbose.Printf("changed Language from %v to %v", b.Language, b1.Language)
			b.Language = b1.Language
		}
	}

	if b1.PageCount != 0 {
		if b.PageCount == 0 {
			Verbose.Printf("set empty PageCount to %v", b1.PageCount)
			b.PageCount = b1.PageCount
		} else if override && (b.PageCount != b1.PageCount) {
			Verbose.Printf("changed PageCount from %v to %v", b.PageCount, b1.PageCount)
			b.PageCount = b1.PageCount
		}
	}

	if len(b1.Subject) > 0 {
		if len(b.Subject) == 0 {
			Verbose.Printf("set empty Subject to %v", b1.Subject)
			b.Subject = append([]string{}, b1.Subject...)
		} else if override && b.compareSubjectWith(b1) != AreTheSame {
			Verbose.Printf("changed Subject from %v to %v", b.Subject, b1.Subject)
			b.Subject = append([]string{}, b1.Subject...)
		}
	}

	if len(b1.Report.Issues) > 0 {
		//TODO: I'm relatively defensive here by reporting any issues even the
		//one encountered on intermediate book attributes consolidation that
		//even might not end in the final result.
		//Does it make sense as b1 issues might not be relevant for the final
		//consolidation of attributes?
		b.Report.Issues = append(b.Report.Issues, b1.Report.Issues...)
	}
}

// CompleteFrom completes Book's attributes by setting empty values to the
// corresponding value of the provided 'b1' Book.
// CompleteFrom is a shortcut to call MergeWith with override set to false.
func (b *Book) CompleteFrom(b1 *Book) {
	b.MergeWith(b1, false)
}

// ReplaceFrom completes and replaces Book's attributes using the non-empty
// corresponding value of the provided 'b1' Book.
// ReplaceFrom is a shortcut to call MergeWith with override set to true.
func (b *Book) ReplaceFrom(b1 *Book) {
	b.MergeWith(b1, true)
}

// CompleteFromMap completes Book's attributes by setting empty values to the
// corresponding value of the provided map.
// map format is similar to NewFromMap.
func (b *Book) CompleteFromMap(m map[string]string) error {
	b1, err := NewFromMap(m)
	if err != nil {
		return err
	}

	b.CompleteFrom(b1)
	return nil
}

// ReplaceFromMap completes and replaces Book's attributes using the non-empty
// corresponding value of the provided map.
// map format is similar to NewFromMap.
func (b *Book) ReplaceFromMap(m map[string]string) error {
	b1, err := NewFromMap(m)
	if err != nil {
		return err
	}

	b.ReplaceFrom(b1)
	return nil
}
