package main

import (
	"embed"

	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pirmd/libro/book"
	"github.com/pirmd/libro/util"
)

//go:embed templates/name/*
var nameTmplDir embed.FS

// Libro represents a collection of media and its associated management
// facilities.
type Libro struct {
	// Root is the location of Libro's library.
	Root string

	// Verbose is the logger for providing low interest messages to the user.
	Verbose *log.Logger

	// Debug is the logger for providing messages supposed to help the
	// developer understand his/her mistakes.
	Debug *log.Logger

	// UseGooglebooks, if set, will complete book's missing metadata by
	// searching Googlebooks.
	// Default to false (do not try fetching missing metadata)
	UseGooglebooks bool

	// MaxSearchResults defines the maximum number of results to consider when
	// looking for a book.
	// Default to 3
	MaxSearchResults int

	// UseGuesser, if set,  tries to complete book's metadata by guessing
	// missing information using media's filename and title.
	// Default to false (do not try guessing missing metadata)
	UseGuesser bool

	// PathTmpl is a text.Template that determines the standardized media
	// files location in the collection based on their metadata.
	// Default to nil (keep item location as-is)
	PathTmpl *template.Template
}

// NewLibro creates a new Libro.
func NewLibro() *Libro {
	tmpl := template.New("location").Option("missingkey=error")
	tmpl = tmpl.Funcs(util.StringsFuncMap).Funcs(util.FilepathFuncMap).Funcs(util.TmplFuncMap(tmpl))
	tmpl = template.Must(tmpl.ParseFS(nameTmplDir, "templates/name/*"))

	return &Libro{
		Root:             ".",
		Verbose:          log.New(io.Discard, "", 0),
		Debug:            log.New(io.Discard, "debug:", 0),
		PathTmpl:         template.Must(tmpl.Parse(`{{template "fullname.gotmpl" .}}`)),
		MaxSearchResults: 3,
	}
}

// Read extracts all possible information about a book.
func (lib *Libro) Read(path string) (*book.Book, error) {
	lib.Verbose.Printf("Read information from book's file")
	b, err := book.NewFromFile(path)
	if err != nil {
		return nil, err
	}

	if lib.UseGuesser {
		lib.Verbose.Print("Clean book's metadata")
		if err := b.CleanMetadata(); err != nil {
			return nil, err
		}

		lib.Verbose.Print("Guess information from book's Filename")
		if err := lib.guessFromFilename(b); err != nil {
			return nil, err
		}

		lib.Verbose.Print("Guess information from book's Content")
		if err := lib.guessFromContent(b); err != nil {
			return nil, err
		}
	}

	if lib.UseGooglebooks {
		lib.Verbose.Print("Get book from Googlebooks")
		if err := lib.searchOnGooglebooks(b); err != nil {
			return nil, err
		}
	}

	if lib.UseGuesser {
		lib.Verbose.Print("Guess information from book's Metadata")
		if err := b.GuessFromMetadata(); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// Create inserts a new book in Libro's collection.
//
// It determines the location in the collection by executing Libro.PathTmpl
// against book's metadata. Location can be relative or absolute, relative
// location are relative to the Libro's root folder.
// Location can contain reference to environment variables that are expanded to
// determine the target location to store the book's file.
//
// Create operation will fail if the target location already exists.
func (lib *Libro) Create(b *book.Book) error {
	lib.Verbose.Printf("Insert book into library in '%s'", lib.Root)

	if b.Path == "" {
		lib.Verbose.Printf("Done (no file attached to book)")
		return nil
	}

	if lib.PathTmpl == nil {
		lib.Verbose.Printf("Done (no template to relocate book)")
		return nil
	}

	buff := new(bytes.Buffer)
	if err := lib.PathTmpl.Execute(buff, b); err != nil {
		return err
	}
	path := filepath.Clean(os.ExpandEnv(buff.String()))

	dst := lib.fullpath(path)
	lib.Debug.Printf("new location of book is '%s'", dst)

	dontNeedCopy, err := util.SamePath(dst, b.Path)
	if err != nil {
		return err
	}
	if dontNeedCopy {
		lib.Verbose.Printf("Done (destination is the same as current one)")
		return nil
	}

	lib.Verbose.Printf("copy book to '%s'", dst)
	if err := util.CopyFile(dst, b.Path); err != nil {
		return err
	}

	b.Path = path

	return nil
}

// fullpath returns the full path to interact with Libro's collection. If
// path is relative, fullpath returns its full location inside Libro's
// root folder.  If path is absolute, fullpath returns its "clean"
// representation (filepath.Clean).
//
// Relative path are "secured" to some point by ignoring any indication
// pointing outside of Libro's root.
//
// fullPath does not check whether the fullpath exists or makes sense.
func (lib *Libro) fullpath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(lib.Root, filepath.Clean("/"+path))
}

func (lib *Libro) guessFromFilename(b *book.Book) error {
	guessedBook, err := book.NewFromFilename(b.Path)
	if err != nil {
		return err
	}

	if guessedBook == nil {
		lib.Debug.Print("Filename contains no relevant information")
		return nil
	}

	switch lvl, rational := b.CompareWith(guessedBook); lvl {
	case book.AreTheSame, book.AreAlmostTheSame:
		lib.Debug.Print("Filename is consistent with Metadata")
		b.CompleteFrom(guessedBook)

	case book.AreNotComparable:
		// Either Book's Metadata or Book's Filename where not complete
		// enough to determine whether they are consistent or not. We try
		// to get most use of available information.
		lib.Debug.Print("cannot check whether Filename is consistent with Metadata")
		b.CompleteFrom(guessedBook)

	default:
		b.ReportIssue("filename is not consistent with Metadata: %s", rational)
	}

	return nil
}

func (lib *Libro) guessFromContent(b *book.Book) error {
	guessedBook, err := book.NewFromContent(b.Path)
	if err != nil {
		return err
	}

	if guessedBook == nil {
		lib.Debug.Print("Content contains no relevant information")
		return nil
	}

	switch lvl, rational := b.CompareWith(guessedBook); lvl {
	case book.AreTheSame, book.AreAlmostTheSame:
		lib.Debug.Print("Content is consistent with Metadata")
		b.CompleteFrom(guessedBook)

	case book.AreNotComparable:
		// Either Book's Metadata or Book's Filename where not complete enough
		// to determine whether they are consistent or not. We try to get most
		// use of available information.
		lib.Debug.Print("cannot check whether Content is consistent with Metadata")
		b.CompleteFrom(guessedBook)

	default:
		b.ReportIssue("content is not consistent with Metadata: %s", rational)
	}

	return nil
}

func (lib *Libro) searchOnGooglebooks(b *book.Book) error {
	matches, err := b.SearchOnGooglebooks(lib.MaxSearchResults)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		b.ReportWarning("no match found on Googlebooks")
		return nil
	}

	bestMatch := matches[0]
	for _, m := range matches {
		if m.ISBN == b.ISBN {
			bestMatch = m
			break
		}
	}

	switch lvl, rational := b.CompareWith(bestMatch); lvl {
	case book.AreTheSame:
		lib.Debug.Printf("Googlebooks returns the same book: %s", rational)
		b.ReplaceFrom(bestMatch)

	case book.AreAlmostTheSame:
		lib.Debug.Printf("Googlebooks returns a similar book: %s", rational)
		b.CompleteFrom(bestMatch)

	case book.AreNotTheSame:
		b.ReportIssue("Googlebooks returns a different book: %s", rational)

	default:
		// difficult to take an informed decision, leave it to the end-user.
		for _, match := range matches {
			b.ReportSimilarBook(match)
		}
	}

	return nil
}
