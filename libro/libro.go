package libro

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pirmd/libro/book"
)

const (
	// defaultLocation is the default location scheme when creating/updating a
	// new element in Libro's library.
	// By default, Libro's sorts books as:
	//  Author - [Series SeriesIndex] - SeriesTitle [Language].Ext
	defaultLocation = `
    {{- define "author" }}{{if .Authors}}{{index .Authors 0}}{{else}}unknown{{end}}{{end -}}
    {{- define "series" }}{{if .Series}} - [{{.Series}} {{.SeriesIndex}}]{{end}}{{end -}}
    {{- define "title" }}{{if .SeriesTitle}} - {{.SeriesTitle}}{{else}} - {{.Title}}{{end}}{{end -}}
    {{- define "lang" }}{{if .Language}} [{{.Language}}]{{end}}{{end -}}

    {{- print (tmpl "author" .) (tmpl "series" . ) (tmpl "title" .) (tmpl "lang" .) (ext .Path) }}`
)

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

	// UseGuesser, if set,  tries to complete book's metadata by guessing
	// missing information using media's filename and title.
	// Default to false (do not try guessing missing metadata)
	UseGuesser bool

	// LocationTmpl is a text.Template that determines the standardized media
	// files location in the collection based on their metadata.
	// Default to nil (keep item location as-is)
	LocationTmpl *template.Template
}

// New creates a new Libro.
func New() *Libro {
	tmpl := template.New("location").Option("missingkey=error")
	tmpl = tmpl.Funcs(FilepathFuncMap).Funcs(TmplFuncMap(tmpl))

	return &Libro{
		Root:         ".",
		Verbose:      log.New(io.Discard, "", 0),
		Debug:        log.New(io.Discard, "debug:", 0),
		LocationTmpl: template.Must(tmpl.Parse(defaultLocation)),
	}
}

// Read extracts all possible information about a book.
func (lib *Libro) Read(path string) (*book.Book, error) {
	book.Verbose, book.Debug = lib.Verbose, lib.Debug

	lib.Verbose.Printf("Read information from book's file")
	b, err := book.NewFromFile(path)
	if err != nil {
		return nil, err
	}

	// TODO: split b.Guess to separate guessing from filename from guessing Series and Cie
	//       guess series at the end (after googlebooks) as googlebooks does
	//       not provide series information.
	if lib.UseGuesser {
		lib.Verbose.Print("Guess information from book's Filename and Title")
		if err := b.Guess(); err != nil {
			return nil, err
		}
	}

	if lib.UseGooglebooks {
		lib.Verbose.Print("Get book's information from Googlebooks")
		if err := b.FromGooglebooks(); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// Create inserts a new book in Libro's collection.
//
// It determines the location in the collection by executing Libro.LocationTmpl
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

	if lib.LocationTmpl == nil {
		lib.Verbose.Printf("Done (no template to relocate book)")
		return nil
	}

	buff := new(bytes.Buffer)
	if err := lib.LocationTmpl.Execute(buff, b); err != nil {
		return err
	}
	path := filepath.Clean(os.ExpandEnv(buff.String()))

	dst := lib.fullpath(path)
	lib.Debug.Printf("new location of book is '%s'", dst)

	dontNeedCopy, err := samePath(dst, b.Path)
	if err != nil {
		return err
	}
	if dontNeedCopy {
		lib.Verbose.Printf("Done (destination is the same as current one)")
		return nil
	}

	lib.Verbose.Printf("copy book to '%s'", dst)
	if err := copyFile(dst, b.Path); err != nil {
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
