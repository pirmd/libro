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
	reList = regexp.MustCompile(`\s?[&,]\s?`)
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
	ISBN string `json:",omitempty"`

	// SubTitle is the book's sub-title.
	SubTitle string `json:",omitempty"`

	// Publisher is the publisher of this book.
	Publisher string `json:",omitempty"`

	// PublishedDate is the date of publication of this book.
	PublishedDate string `json:",omitempty"`

	// Description is the synopsis of the book. The text of the description
	// is formatted in HTML and includes simple formatting elements.
	Description string `json:",omitempty"`

	// Series is the series to which this book belongs to.
	Series string `json:",omitempty"`

	// SeriesIndex is the position in the series to which the book belongs to.
	SeriesIndex float64 `json:",omitempty"`

	// SeriesTitle is the book's title in th eseries (whithout Series nor
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
}

// New creates a new Book
func New() *Book {
	return &Book{}
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
			b.Authors = reList.Split(value, -1)

		case "Publisher":
			b.Publisher = value

		case "PublishedDate":
			b.PublishedDate = value

		case "Description":
			b.Description = value

		case "Series":
			b.Series = value

		case "SeriesIndex":
			var err error
			if b.SeriesIndex, err = strconv.ParseFloat(value, 32); err != nil {
				return nil, fmt.Errorf("cannot assign %s to '%s': %v", value, a, err)
			}

		case "ISBN":
			b.ISBN = value

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
		if b.Title != "" && !strings.EqualFold(b.Title, b1.Title) {
			Debug.Printf("new Title (%s) is different from the existing one (%s)", b1.Title, b.Title)
		}
		if override || b.Title == "" {
			Verbose.Printf("sets Title to %s", b1.Title)
			b.Title = b1.Title
		}
	}

	if len(b1.Authors) > 0 {
		if len(b.Authors) != 0 && !strings.EqualFold(fmt.Sprint(b.Authors), fmt.Sprint(b1.Authors)) {
			Debug.Printf("new Authors (%v) is different from the existing one (%v)", b1.Authors, b.Authors)
		}
		if override || len(b.Authors) == 0 {
			Verbose.Printf("sets Authors to %v", b1.Authors)
			b.Authors = append([]string{}, b1.Authors...)
		}
	}

	if b1.ISBN != "" {
		// TODO: implements better ISBN  comparison and replacement logic (if ISBN_10 vs. same in ISBN_13)
		if b.ISBN != "" && b.ISBN != b1.ISBN {
			Debug.Printf("new ISBN (%s) is different from the existing one (%s)", b1.ISBN, b.ISBN)
		}
		if override || b.ISBN == "" {
			Verbose.Printf("sets ISBN to %s", b1.ISBN)
			b.ISBN = b1.ISBN
		}
	}

	if b1.SubTitle != "" {
		if b.SubTitle != "" && !strings.EqualFold(b.SubTitle, b1.SubTitle) {
			Debug.Printf("new SubTitle (%s) is different from the existing one (%s)", b1.SubTitle, b.SubTitle)
		}
		if override || b.SubTitle == "" {
			Verbose.Printf("sets SubTitle to %s", b1.SubTitle)
			b.SubTitle = b1.SubTitle
		}
	}

	if b1.Publisher != "" {
		if b.Publisher != "" && !strings.EqualFold(b.Publisher, b1.Publisher) {
			Debug.Printf("new Publisher (%s) is different from the existing one (%s)", b1.Publisher, b.Publisher)
		}
		if override || b.Publisher == "" {
			Verbose.Printf("sets Publisher to %s", b1.Publisher)
			b.Publisher = b1.Publisher
		}
	}

	if b1.PublishedDate != "" {
		// TODO: implements better date comparison and replacement logic (if similar, use the most specific date)
		if b.PublishedDate != "" && b.PublishedDate != b1.PublishedDate {
			Debug.Printf("new PublishedDate (%s) is different from the existing one (%s)", b1.PublishedDate, b.PublishedDate)
		}
		if override || b.PublishedDate == "" {
			Verbose.Printf("sets PublishedDate to %s", b1.PublishedDate)
			b.PublishedDate = b1.PublishedDate
		}
	}

	if b1.Description != "" {
		if b.Description != "" && !strings.EqualFold(b.Description, b1.Description) {
			Debug.Printf("new Description (%.12v) is different from the existing one (%.12v)", b1.Description, b.Description)
		}
		if override || b.Description == "" {
			Verbose.Printf("sets Description to %s", b1.Description)
			b.Description = b1.Description
		}
	}

	if b1.Series != "" {
		if b.Series != "" && !strings.EqualFold(b.Series, b1.Series) {
			Debug.Printf("new Series (%s) is different from the existing one (%s)", b1.Series, b.Series)
		}
		if override || b.Series == "" {
			Verbose.Printf("sets Series to %s", b1.Series)
			b.Series = b1.Series
		}
	}

	if b1.SeriesIndex != 0 {
		if b.SeriesIndex != 0 && b.SeriesIndex != b1.SeriesIndex {
			Debug.Printf("new SeriesIndex (%.1f) is different from the existing one (%.1f)", b1.SeriesIndex, b.SeriesIndex)
		}
		if override || b.SeriesIndex == 0 {
			Verbose.Printf("sets SeriesIndex to %.1f", b1.SeriesIndex)
			b.SeriesIndex = b1.SeriesIndex
		}
	}

	if b1.SeriesTitle != "" {
		if b.SeriesTitle != "" && !strings.EqualFold(b.SeriesTitle, b1.SeriesTitle) {
			Debug.Printf("new SeriesTitle (%s) is different from the existing one (%s)", b1.SeriesTitle, b.SeriesTitle)
		}
		if override || b.SeriesTitle == "" {
			Verbose.Printf("sets SeriesTitle to %s", b1.SeriesTitle)
			b.SeriesTitle = b1.SeriesTitle
		}
	}

	if b1.Language != "" {
		if b.Language != "" && !strings.EqualFold(b.Language, b1.Language) {
			Debug.Printf("new Language (%s) is different from the existing one (%s)", b1.Language, b.Language)
		}
		if override || b.Language == "" {
			Verbose.Printf("sets Language to %s", b1.Language)
			b.Language = b1.Language
		}
	}

	if b1.PageCount != 0 {
		if b.PageCount != 0 && b.PageCount != b1.PageCount {
			Debug.Printf("new PageCount (%d) is different from the existing one (%d)", b1.PageCount, b.PageCount)
		}
		if override || b.PageCount == 0 {
			Verbose.Printf("sets PageCount to %d", b1.PageCount)
			b.PageCount = b1.PageCount
		}
	}

	if len(b1.Subject) > 0 {
		if len(b.Subject) != 0 && !strings.EqualFold(fmt.Sprint(b.Subject), fmt.Sprint(b1.Subject)) {
			Debug.Printf("new Subject (%v) is different from the existing one (%v)", b1.Subject, b.Subject)
		}
		if override || len(b.Subject) == 0 {
			Verbose.Printf("sets new Subject to %v", b1.Subject)
			b.Subject = append([]string{}, b1.Subject...)
		}
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

// String proposes an easy-to-read raw representation of a Book.
func (b Book) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "Path         :\t%s", b.Path)
	fmt.Fprintf(&s, "\nTitle        :\t%s", b.Title)
	fmt.Fprintf(&s, "\nAuthors      :\t%s", strings.Join(b.Authors, " & "))

	if b.ISBN != "" {
		fmt.Fprintf(&s, "\nISBN         :\t%s", b.ISBN)
	}

	if b.SubTitle != "" {
		fmt.Fprintf(&s, "\nSubTitle     :\t%s", b.SubTitle)
	}

	if b.SeriesTitle != "" || b.Series != "" || b.SeriesIndex != 0 {
		fmt.Fprintf(&s, "\nSeriesTitle  :\t%s", b.SeriesTitle)
		fmt.Fprintf(&s, "\nSeries       :\t%s", b.Series)
		fmt.Fprintf(&s, "\nSeriesIndex  :\t%.1f", b.SeriesIndex)
	}

	if b.Description != "" {
		fmt.Fprintf(&s, "\nDescription  :\t%s", b.Description)
	}

	if b.Publisher != "" {
		fmt.Fprintf(&s, "\nPublisher    :\t%s", b.Publisher)
	}

	if b.PublishedDate != "" {
		fmt.Fprintf(&s, "\nPublishedDate:\t%s", b.PublishedDate)
	}

	if b.Language != "" {
		fmt.Fprintf(&s, "\nLanguage     :\t%s", b.Language)
	}

	if b.PageCount > 0 {
		fmt.Fprintf(&s, "\nPageCount    :\t%d", b.PageCount)
	}

	if len(b.Subject) > 0 {
		fmt.Fprintf(&s, "\nSubject      :\t%s", strings.Join(b.Subject, " & "))
	}

	return s.String()
}
