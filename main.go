package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pirmd/libro/book"
	"github.com/pirmd/libro/util"
)

var (
	myname    = filepath.Base(os.Args[0])
	myversion = "v?.?.?-?" //should be set using: go build -ldflags "-X main.myversion=X.X.X"

	//go:embed templates/book/*
	bookTmplDir embed.FS
)

// App is a wrapper around a Libro object that implements command line
// facilities to interact with the user.
// App is also supposed to offer a not that complicated API for testing the
// app.
type App struct {
	// Verbose is the logger for providing low interest messages to the user.
	Verbose *log.Logger

	// Debug is the logger for providing messages supposed to help the
	// developer understand his/her mistakes.
	Debug *log.Logger

	// Library points to the underlying Libro object.
	Library *Libro

	// Formatter is the go template used to pretty print an object.
	Formatter *template.Template

	// Stdout is the standard output where to print app's result.
	// Except for test, it is usually os.Sdtout.
	Stdout io.Writer
}

// NewApp creates a new App
func NewApp() *App {
	tmpl := template.New("formatter").Option("missingkey=error")
	tmpl = tmpl.Funcs(util.SerializationFuncMap).Funcs(util.StringsFuncMap)
	tmpl = template.Must(tmpl.ParseFS(bookTmplDir, "templates/book/*.gotmpl"))

	app := &App{
		Verbose:   log.New(io.Discard, "", 0),
		Debug:     log.New(io.Discard, "debug:", 0),
		Stdout:    os.Stdout,
		Library:   NewLibro(),
		Formatter: template.Must(tmpl.Parse(`{{toJSON .}}`)),
	}

	// if we are printing to a TTY, we use a format that is easier to read for a human.
	if fi, _ := os.Stdout.Stat(); (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		template.Must(app.Formatter.Parse(`{{template "book.txt.gotmpl" .}}`))
	}

	app.Library.Verbose, app.Library.Debug = app.Verbose, app.Debug
	book.Verbose, book.Debug = app.Verbose, app.Debug

	return app
}

// Run runs the App.
func (app *App) Run(args []string) error {
	fs := flag.NewFlagSet(myname, flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [option...] <commands> [arguments]\n", myname)
		fmt.Fprintf(fs.Output(), "Commands:\n")
		fmt.Fprintf(fs.Output(), "    info       retrieve information from an epub\n")
		fmt.Fprintf(fs.Output(), "    insert     insert an epub to the library\n")
		fmt.Fprintf(fs.Output(), "    edit       edit information about an epub\n")
		fmt.Fprintf(fs.Output(), "    version    print %s version\n", myname)
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	fs.Var(util.NewLogSwitcher(app.Verbose), "verbose", "print messages of low interest")
	fs.Var(util.NewLogSwitcher(app.Debug, app.Verbose), "debug", "print cryptic messages supposed to help the developer understand his/her mistakes")
	fs.Var(util.NewGoTemplate(app.Formatter), "format", "set output format using golang text/template")
	fs.Var(util.NewGoTemplateFS(app.Formatter), "format-tmpl", "loads user-defined format template(s) from golang text/template definition files")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("%v\nRun %s -help", err, fs.Name())
	}

	switch cmd := fs.Arg(0); cmd {
	case "version":
		fmt.Fprintf(fs.Output(), "%s version %s\n", myname, myversion)

	case "info":
		return app.RunInfoSubcmd(fs.Args()[1:])

	case "insert", "add":
		return app.RunInsertSubcmd(fs.Args()[1:])

	case "edit":
		return app.RunEditSubcmd(fs.Args()[1:])

	default:
		return fmt.Errorf("'%[1]s %s' unknown command\nRun %[1]s -help", fs.Name(), cmd)
	}

	return nil
}

// RunInfoSubcmd executes the "info" sub-command.
func (app *App) RunInfoSubcmd(args []string) error {
	fs := flag.NewFlagSet(myname+" info", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [option...] FILENAME\n", fs.Name())
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	fs.BoolVar(&app.Library.UseGuesser, "use-guesser", false, "completes book's metadata by guessing lacking information from book's filename and title")
	fs.BoolVar(&app.Library.UseGooglebooks, "use-googlebooks", false, "completes book's metadata by searching lacking information from Googlebooks")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("%v\nRun %s -help", err, fs.Name())
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("wrong number of arguments\nRun %s -help", fs.Name())
	}
	path := fs.Arg(0)

	b, err := app.Library.Read(path)
	if err != nil {
		return fmt.Errorf("cannot retrieve information about '%s': %v", path, err)
	}

	if err := app.Formatter.Execute(app.Stdout, b); err != nil {
		return fmt.Errorf("fail to display book information: %v", err)
	}
	fmt.Fprint(app.Stdout)

	return nil
}

// RunInsertSubcmd executes the "insert" sub-command.
func (app *App) RunInsertSubcmd(args []string) error {
	fs := flag.NewFlagSet(myname+" insert", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [option...] BOOKinJSON\n", fs.Name())
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	fs.StringVar(&app.Library.Root, "root", app.Library.Root, "root folder where the books library is to be found")
	fs.Var(util.NewGoTemplate(app.Library.PathTmpl), "rename", "sets filename format using golang text/template")
	fs.Var(util.NewGoTemplateFS(app.Library.PathTmpl), "rename-tmpl", "loads user-defined filename template(s) from golang text/template definition files")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("%v\nRun %s -help", err, fs.Name())
	}

	var bookJSON io.Reader
	switch fs.NArg() {
	case 0:
		if fi, _ := os.Stdin.Stat(); (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return fmt.Errorf("invalid number of argument(s)\nRun %s -help", fs.Name())
		}
		bookJSON = os.Stdin
	case 1:
		bookJSON = strings.NewReader(fs.Arg(0))
	default:
		return fmt.Errorf("invalid number of argument(s)\nRun %s -help", fs.Name())
	}

	b := book.New()
	if err := json.NewDecoder(bookJSON).Decode(&b); err != nil {
		return fmt.Errorf("fail to decode book's JSON: %v", err)
	}

	if err := app.Library.Create(b); err != nil {
		return fmt.Errorf("fail to add new book: %v", err)
	}

	if err := app.Formatter.Execute(app.Stdout, b); err != nil {
		return fmt.Errorf("fail to display book information: %v", err)
	}
	fmt.Fprint(app.Stdout)

	return nil
}

// RunEditSubcmd executes the "edit" sub-command.
func (app *App) RunEditSubcmd(args []string) error {
	fs := flag.NewFlagSet(myname+" edit", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [option...] BOOKinJSON\n", fs.Name())
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	var auto bool
	fs.BoolVar(&auto, "auto", false, "do not trigger an editor if libro get the impression that Book's information is good enough")

	var dontedit bool
	fs.BoolVar(&dontedit, "dont-edit", false, "do not trigger any editor at all. Supersedes 'auto' flag")

	var editor string
	fs.StringVar(&editor, "editor", os.Getenv("EDITOR"), "sets editor's name to use for editing Book's information")

	setAttr := make(map[string]string)
	fs.Var(util.NewKV(setAttr), "set", "set a new value for a book's attribute (format attribute=value)")

	defaultAttr := make(map[string]string)
	fs.Var(util.NewKV(defaultAttr), "default", "set a new value for a book's attribute if the attribute is not yet set (format attribute=value)")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("%v\nRun %s -help", err, fs.Name())
	}

	var bookJSON io.Reader
	switch fs.NArg() {
	case 0:
		if fi, _ := os.Stdin.Stat(); (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return fmt.Errorf("invalid number of argument(s)\nRun %s -help", fs.Name())
		}
		bookJSON = os.Stdin
	case 1:
		bookJSON = strings.NewReader(fs.Arg(0))
	default:
		return fmt.Errorf("invalid number of argument(s)\nRun %s -help", fs.Name())
	}

	b := book.New()
	if err := json.NewDecoder(bookJSON).Decode(&b); err != nil {
		return fmt.Errorf("fail to decode book's JSON: %v", err)
	}

	if len(defaultAttr) != 0 {
		app.Verbose.Printf("Set default value for book's information")
		if err := b.CompleteFromMap(defaultAttr); err != nil {
			return fmt.Errorf("fail to set default value: %v", err)
		}
	}

	if len(setAttr) != 0 {
		app.Verbose.Printf("Set new value for book's information")
		if err := b.ReplaceFromMap(setAttr); err != nil {
			return fmt.Errorf("fail to set new value: %v", err)
		}
	}

	app.Verbose.Printf("Edit book's information")
	switch {
	case dontedit:
		app.Verbose.Printf("manual edition of book's information has been prevented by '-dont-edit' flag")
	case editor == "":
		app.Verbose.Printf("no editor has been defined. Set $EDITOR global var or use -editor command line flag")
	case auto && (b.IsComplete() && len(b.ToReview) == 0):
		app.Verbose.Printf("no need to edit book's information that seems good enough to me")
	default:
		var err error
		if b, err = editBook(editor, b); err != nil {
			return fmt.Errorf("fail to edit book: %v", err)
		}
	}

	if err := app.Formatter.Execute(app.Stdout, b); err != nil {
		return fmt.Errorf("fail to display book information: %v", err)
	}
	fmt.Fprint(app.Stdout)

	return nil
}

func main() {
	app := NewApp()

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}
}
