package book

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
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
