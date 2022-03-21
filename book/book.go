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

    // reList is a regexp that splits a list of values (like Authors or Categories).
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

	// ShortTitle is the book's short title (whithout Series nor SubTitle
	// information).
	ShortTitle string `json:",omitempty"`

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

	// Language is the book's language. It is the two-letter ISO 639-1 code
	// such as 'fr', 'en'.
	Language string `json:",omitempty"`

	// PageCount is total number of pages of this book.
	PageCount int `json:",omitempty"`

	// Categories is the list of subject categories, such as "Fiction",
	// "Suspense".
	Categories []string `json:",omitempty"`
}

// New creates a new Book
func New() *Book {
	return &Book{}
}

// NewFromFile creates a new Book and populates its information according to
// the file's metadata.
func NewFromFile(path string) (*Book, error) {
	b := New()
	b.Path = path

	if err := b.FromFile(); err != nil {
		return nil, err
	}

	return b, nil
}

// FromFile populates Book's information according to the file's metadata.
func (b *Book) FromFile() error {
	Debug.Printf("looking for book's metadata in '%s'", b.Path)

	switch ext := filepath.Ext(b.Path); ext {
	case ".epub":
		if err := b.FromEpub(); err != nil {
			return err
		}

	default:
		return ErrUnknownFormat
	}

	Debug.Printf("found book's metadata: '%#v'", b)
	return nil
}

// FromMap updates a Book's information according to the attributes defined as
// a map where keys are attribute's name (insensitive to case) and value
// is a string representation of the attribute's value.
// For attributes that accept a list of values (like Authors or Categories),
// provided map value should be formatted like "val0 & val1" (individual value
// in as string separated by '&').
// If override flag is on, existing Book's attribute is replaced by the one
// provided in the map ; if override flag is off, Book's attribute is only
// replaced if empty.
func (b *Book) FromMap(m map[string]string, override bool) error {
    Debug.Printf("update book's information from: '%+v'", m)

    for attr, value := range m {
        switch a := strings.Title(attr); a {
        case "Title":
            if b.Title != "" && strings.ToLower(b.Title) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.Title)
            }

            if override || b.Title == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.Title = value
            }

        case "SubTitle":
            if b.SubTitle != "" && strings.ToLower(b.SubTitle) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.SubTitle)
            }

            if override || b.SubTitle == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.SubTitle = value
            }

        case "ShortTitle":
            if b.ShortTitle != "" && strings.ToLower(b.ShortTitle) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.ShortTitle)
            }

            if override || b.ShortTitle == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.ShortTitle = value
            }

        case "Authors":
            v := reList.Split(value, -1)

            if len(b.Authors) != 0 && strings.ToLower(fmt.Sprint(b.Authors)) != strings.ToLower(fmt.Sprint(v)) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%v != %v)", a, v, b.Authors)
            }

            if override || len(b.Authors) == 0 {
                Verbose.Printf("sets new Book's value: %s = %v", a, v)
                b.Authors = v
            }

        case "Publisher":
            if b.Publisher != "" && strings.ToLower(b.Publisher) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.Publisher)
            }

            if override || b.Publisher == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.Publisher = value
            }

        case "PublishedDate":
            if b.PublishedDate != "" && strings.ToLower(b.PublishedDate) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.PublishedDate)
            }

            if override || b.PublishedDate == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.PublishedDate = value
            }

        case "Description":
            if b.Description != "" && strings.ToLower(b.Description) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.Description)
            }

            if override || b.Description == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.Description = value
            }

        case "Series":
            if b.Series != "" && strings.ToLower(b.Series) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.Series)
            }

            if override || b.Series == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.Series = value
            }

        case "SeriesIndex":
            v, err := strconv.ParseFloat(value, 32)
            if err != nil {
                return fmt.Errorf("cannot assign %s to '%s': %v.", value, a, err)
            }

            if b.SeriesIndex != 0 && b.SeriesIndex != v {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%.1f != %.1f)", a, v, b.SeriesIndex)
            }

            if override || b.SeriesIndex == 0 {
                Verbose.Printf("sets new value: %s = %.1f.", a, v)
                b.SeriesIndex = v
            }

        case "ISBN":
            if b.ISBN != "" && strings.ToLower(b.ISBN) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.ISBN)
            }

            if override || b.ISBN == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.ISBN = value
            }

        case "Language":
            if b.Language != "" && strings.ToLower(b.Language) != strings.ToLower(value) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%s != %s)", a, value, b.Language)
            }

            if override || b.Language == "" {
                Verbose.Printf("sets new Book's value: %s = %s", a, value)
                b.Language = value
            }

        case "PageCount":
            v, err := strconv.Atoi(value)
            if err != nil {
                return fmt.Errorf("cannot assign %s to '%s': %v.", value, a, err)
            }

            if b.PageCount != 0 && b.PageCount != v {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%d != %d)", a, v, b.PageCount)
            }

            if override || b.PageCount == 0 {
                Verbose.Printf("sets new value: %s = %d.", a, v)
                b.PageCount = v
            }

        case "Categories":
            v := reList.Split(value, -1)

            if len(b.Categories) != 0 && strings.ToLower(fmt.Sprint(b.Categories)) != strings.ToLower(fmt.Sprint(v)) {
                Debug.Printf("new Book's value for '%s' is different from the existing one (%v != %v)", a, v, b.Categories)
            }

            if override || len(b.Categories) == 0 {
                Verbose.Printf("sets new Book's value: %s = %v", a, v)
                b.Categories = v
            }

        default:
            return fmt.Errorf("cannot set unknown attribute '%s'", a)
        }
    }

    return nil
}

// String proposes an easy-to-read raw representation of a Book.
func (b Book) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "Path         :\t%v", b.Path)
	fmt.Fprintf(&s, "\nTitle        :\t%v", b.Title)
	fmt.Fprintf(&s, "\nAuthors      :\t%v", b.Authors)

	if b.ISBN != "" {
		fmt.Fprintf(&s, "\nISBN         :\t%v", b.ISBN)
	}

	if b.SubTitle != "" {
		fmt.Fprintf(&s, "\nSubTitle     :\t%v", b.SubTitle)
	}

	if b.ShortTitle != "" {
		fmt.Fprintf(&s, "\nShortTitle   :\t%v", b.ShortTitle)
	}

	if b.Series != "" || b.SeriesIndex != 0 {
		fmt.Fprintf(&s, "\nSeries       :\t%v", b.Series)
		fmt.Fprintf(&s, "\nSeriesIndex  :\t%.1f", b.SeriesIndex)
	}

	if b.Description != "" {
		fmt.Fprintf(&s, "\nDescription  :\t%v", b.Description)
	}

	if b.Publisher != "" {
		fmt.Fprintf(&s, "\nPublisher    :\t%v", b.Publisher)
	}

	if b.PublishedDate != "" {
		fmt.Fprintf(&s, "\nPublishedDate:\t%v", b.PublishedDate)
	}

	if b.Language != "" {
		fmt.Fprintf(&s, "\nLanguage     :\t%v", b.Language)
	}

	if b.PageCount > 0 {
		fmt.Fprintf(&s, "\nPageCount    :\t%d", b.PageCount)
	}

	if len(b.Categories) > 0 {
		fmt.Fprintf(&s, "\nCategories   :\t%v", b.Categories)
	}

	return s.String()
}
